import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { fetchAppById, updateApp } from "../../services/app";

export default function UpdateApp() {
  const navigate = useNavigate();
  const [idInput, setIdInput] = useState("");
  const [form, setForm] = useState({
    name: "",
    description: "",
    icon: "",
    type: "",
    tags: "",
    status: "published",
  });
  const [loading, setLoading] = useState(false);
  const [loadingApp, setLoadingApp] = useState(false);

  useEffect(() => {
    // try read id from query string
    try {
      const params = new URLSearchParams(window.location.search);
      const id = params.get("id") || params.get("app_id");
      if (id) {
        setIdInput(id);
        loadApp(id);
      }
    } catch (e) {}
  }, []);

  const onChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) =>
    setForm({ ...form, [e.target.name]: e.target.value });

  const loadApp = async (id: string) => {
    setLoadingApp(true);
    try {
      const data = await fetchAppById(id);
      setForm({
        name: data.Name || data.name || "",
        description: data.Description || data.description || "",
        icon: data.Icon || data.icon || "",
        type: data.Type || data.type || "",
        tags: data.Tags || data.tags || "",
        status: data.Status || data.status || "published",
      });
    } catch (e) {
      alert("加载应用信息失败");
    } finally {
      setLoadingApp(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!idInput) {
      alert("请输入要更新的应用 ID");
      return;
    }
    setLoading(true);
    try {
      const id = Number(idInput);
      const payload = {
        id,
        name: form.name.trim(),
        description: form.description.trim(),
        icon: form.icon.trim(),
        type: form.type.trim(),
        tags: form.tags.trim(),
        status: form.status,
      };
      await updateApp(id, payload);
      navigate("/home");
    } catch (err: any) {
      alert("更新失败: " + (err.response?.data?.error || err.message || "未知错误"));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="relative w-screen h-screen flex items-center justify-center overflow-hidden">
      <div className="absolute inset-0 bg-gradient-to-br from-blue-200 via-blue-100 to-white"></div>
      <div className="relative z-10 bg-white shadow-2xl rounded-2xl p-8 w-full max-w-2xl">
        <h2 className="text-2xl font-bold mb-4">更新应用</h2>

        <div className="mb-4">
          <label className="block text-sm font-medium mb-1">应用 ID</label>
          <div className="flex gap-2">
            <input value={idInput} onChange={(e) => setIdInput(e.target.value)} className="flex-1 p-2 border rounded" />
            <button onClick={() => loadApp(idInput)} disabled={!idInput || loadingApp} className="px-4 py-2 bg-gray-600 text-white rounded">
              {loadingApp ? "加载中..." : "加载"}
            </button>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">应用名称</label>
            <input name="name" value={form.name} onChange={onChange} required className="w-full p-2 border rounded" />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">应用描述</label>
            <textarea name="description" value={form.description} onChange={onChange} required className="w-full p-2 border rounded" rows={4} />
          </div>
          <div className="grid grid-cols-2 gap-3">
            <input name="type" value={form.type} onChange={onChange} placeholder="类型" className="p-2 border rounded" />
            <input name="status" value={form.status} onChange={onChange} placeholder="状态" className="p-2 border rounded" />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">标签（逗号分隔）</label>
            <input name="tags" value={form.tags} onChange={onChange} className="w-full p-2 border rounded" />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">图标 URL（可选）</label>
            <input name="icon" value={form.icon} onChange={onChange} className="w-full p-2 border rounded" />
          </div>
          <button type="submit" disabled={loading} className="w-full py-2 bg-blue-600 text-white rounded">
            {loading ? "更新中..." : "更新应用"}
          </button>
        </form>
      </div>
    </div>
  );
}
