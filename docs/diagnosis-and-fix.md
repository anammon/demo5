# AssistApp 部署问题诊断与解决文档

> 文件目的：记录外网无法访问 / 登录失败问题的诊断过程、根因、所做修改、验证步骤与后续建议，便于归档与复现。

## 一、问题概述

- 表现：通过 cloudflared 暴露的域名可以访问前端静态页面，但登录/注册等 API 请求失败，浏览器报错：
  "Cross-Origin Request Blocked... http://localhost:3000/user/login (Reason: CORS request did not succeed). Status code: (null)"，同时 axios 抛出 `Network Error`。
- 附带问题：页面中“漂流瓶”图标有缺失或重复显示；点击漂流瓶后未进入交互式页面（没有扔/捡按钮），而是显示放大的图标或通用详情页。

## 二、核心根因

1. 前端构建产物中硬编码了 `http://localhost:3000`（`VITE_API_BASE` 被注入或回退为 localhost），导致外网客户端访问 cloud domain 时，JS bundle 尝试向客户端本地的 `localhost:3000` 发请求，造成请求未到后端而报 `Network Error`。

2. 后端未映射 `public/` 下的某些静态路径（如 `/icons`），导致图片 404，造成页面缺图或样式异常。

3. 前端在 Home 页面同时渲染了一个固定的“漂流瓶”入口，并且后端 App 列表中又包含名为“漂流瓶”的动态应用，未去重导致重复显示；点击动态卡片会跳到 `/apps/:id`（通用详情），而不是跳到交互式 `/apps/bottle` 页面。

4. 虽然 CORS 是需要关注的点，但真正阻断请求的是前端发向错误主机（localhost）而不是 CORS 中间件未注册——需确保 CORS 中间件在后端最先注册以免误导诊断。

## 三、已实施的修复

以下修改已经在代码库中完成并提交：

- frontend/.env
  - 原：`VITE_API_BASE=http://localhost:3000`
  - 现在：`VITE_API_BASE=`（清空生产注入值）
  - 目的：避免把开发时的 localhost 注入生产 bundle。

- frontend/src/services/api.ts
  - 改为：`const API_BASE = import.meta.env.VITE_API_BASE ?? ""`（默认同源），并保留统一的 axios 实例与 token interceptor。
  - 目的：构建后默认同源请求，生产环境不会向客户端本地发起请求。

- frontend/src/pages/Login.tsx 与 frontend/src/pages/Register.tsx
  - 将直接 `axios.post("http://localhost:3000/...")` 替换为 `api.post("/user/...")`，统一使用 `src/services/api.ts`。
  - 目的：消除硬编码 URL，统一管理 baseURL 与认证 header。

- frontend/src/pages/Home.tsx
  - 添加 `hasBottle` 判断：仅当动态应用列表中不存在“漂流瓶”时才渲染固定入口卡片。
  - 修改点击逻辑：若 app 名称包含“漂流瓶”或英文 `bottle`，点击跳转到 `/apps/bottle`（专用交互页面），否则跳转到 `/apps/:id`（通用详情）。
  - 目的：避免重复入口与错误跳转。

- backend/main.go
  - 增加静态映射：`r.Static("/icons", "D:/gin/demo5/frontend/dist/icons")`；保留 `r.Static("/assets", ...)`。
  - 确认 CORS 中间件 `r.Use(middlewares.CorsMiddleware())` 在最前。
  - 修复了监听地址拼接逻辑（避免 `0.0.0.0::3000`）。
  - 目的：让 `dist/icons/*` 等资源能被正确提供；保证 CORS 头在所有响应中返回。

## 四、验证步骤

1. 构建前端（在 Windows cmd）

```cmd
cd /d d:\gin\demo5\frontend
# 若使用 pnpm
pnpm build
# 或 npm
npm run build
```

2. 检查构建产物中是否存在 localhost 注入（PowerShell）：

```cmd
powershell -Command "Select-String -Path 'dist\\**\\*' -Pattern 'localhost:3000' -SimpleMatch"
```
期望：无匹配。如仍有匹配，则回退至源码查找未替换的硬编码。

3. 启动后端并提供最新 `dist`

```cmd
cd /d d:\gin\demo5
go run main.go
```

4. 在外网通过 cloudflared 域访问页面并在 DevTools → Network 检查：
- 静态资源 `/assets/*`、`/icons/*` 返回 200。
- POST `/user/login` 的 Request URL 指向 cloud domain（或为相对路径 `/user/login`），而不是 `localhost:3000`。
- 点击漂流瓶后，URL 跳转到 `/apps/bottle` 并显示扔瓶子/捡瓶子/匿名选项等交互控件。

## 五、预防与改进建议

1. 构建与环境管理
   - 生产环境不能依赖开发 `.env` 的默认值，建议 CI/CD 注入生产 env，或在代码中默认使用同源空字符串。
   - `.env` 不应提交带开发 localhost 的配置到主分支。

2. API 调用规范
   - 前端所有网络请求必须通过单一 `api` 实例（`frontend/src/services/api.ts`），便于管理 baseURL、token 注入与错误拦截。

3. 静态文件托管
   - 若后端托管前端 dist，建议映射整个 `dist` 根目录给静态文件，并保证 API 路由优先注册或通过 NoRoute 兜底。但要小心根路径映射可能与 API 路由冲突（需按先后注册或 API 前缀匹配）。
   - 更好的做法是使用反向代理（如 nginx）把 `/api` 代理到后端，由 nginx 提供静态文件（更稳定与高性能）。

4. 路由与应用元数据
   - 为每个插件/应用增加 `type` 或 `slug` 字段（后端 App 模型），前端根据该字段决定跳转到哪个专用页面，而不是根据名称字符串匹配。

5. 发布验证脚本
   - 在构建流程中加入简单的 smoke tests：
     - 检查 dist 中不包含 `localhost`。
     - 检查 `icons/` `assets/` 是否存在关键文件。
     - 对关键 API 做一次快速请求，验证返回状态与 CORS 头（若需要）。

## 六、常见排查命令（备份）

- 检查是否有进程监听 3000（Windows）：

```cmd
netstat -ano | findstr 3000
```

- 使用 curl 验证 cloud domain（在可访问外网的终端上）：

```cmd
curl -i -X OPTIONS https://your-cloud-domain/user/login
curl -i -X POST https://your-cloud-domain/user/login -H "Content-Type: application/json" -d "{\"identifier\":\"test\",\"password\":\"pass\"}"
```

## 七、变更清单（供审计）

- `frontend/.env` — 清空 `VITE_API_BASE`
- `frontend/src/services/api.ts` — 统一 baseURL 为 `import.meta.env.VITE_API_BASE ?? ""`
- `frontend/src/pages/Login.tsx` — 使用 `api.post("/user/login")`
- `frontend/src/pages/Register.tsx` — 使用 `api.post("/user/register")`
- `frontend/src/pages/Home.tsx` — 去重漂流瓶入口，点击流中对漂流瓶跳转到 `/apps/bottle`；避免重复渲染
- `backend/main.go` — 增加 `r.Static("/icons", ".../dist/icons")`，CORS middleware 确认在最前，修复端口拼接

## 八、下一步（可选）

- 我可以为你：
  - 在本地执行 `pnpm build` 或 `npm run build` 并把构建结果与检查日志返回给你。
  - 把 App 模型扩展为 `type/slug` 字段，并在前端基于该字段路由到专用页面（如 `bottle`）。
  - 编写一个简单的 GitHub Actions workflow，在每次 push/合并到 main 时自动构建并运行 smoke tests。

如需我生成 PDF 或其他格式文档（例如 Word），请告诉我目标格式，我会把 Markdown 转为指定格式并放在仓库中。

---

文档由修复过程整理，若需要补充日志片段（错误前端 Network / 后端日志），请把对应截图或日志文本贴上，我可以把它们附到该文档中以便审计。