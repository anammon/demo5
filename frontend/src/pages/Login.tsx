import { useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";

export default function Login() {
  const [accountOrEmail, setAccountOrEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await axios.post(
        "http://localhost:3000/user/login",
        {
          identifier: accountOrEmail,
          password,
        },
        { headers: { "Content-Type": "application/json" } }
      );

      localStorage.setItem("token", res.data.token);
      alert("登录成功！欢迎 " + res.data.user.name);
      navigate("/home");
    } catch (err: any) {
      alert("登录失败：" + (err.response?.data?.error || "未知错误"));
    }
  };

  return (
    <div className="relative w-screen h-screen flex items-center justify-center overflow-hidden">
      {/* 背景渐变动画 */}
      <div className="absolute inset-0 bg-gradient-to-br from-blue-200 via-blue-100 to-white animate-gradient-x"></div>
      {/* 光晕效果 */}
      <div className="absolute -top-40 -left-40 w-[500px] h-[500px] bg-blue-300 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-pulse"></div>
      <div className="absolute -bottom-40 -right-40 w-[500px] h-[500px] bg-purple-300 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-pulse delay-700"></div>

      {/* 登录卡片 */}
      <div className="relative z-10 bg-white shadow-2xl rounded-2xl p-12 w-full max-w-md">
        <h2 className="text-3xl font-bold text-center text-blue-600 mb-8">
          欢迎使用 AssistApp
        </h2>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              账户或邮箱
            </label>
            <input
              type="text"
              value={accountOrEmail}
              onChange={(e) => setAccountOrEmail(e.target.value)}
              required
              className="w-full px-4 py-3 rounded-xl border border-gray-300 focus:ring-2 focus:ring-blue-400 focus:outline-none"
              placeholder="输入账户名或邮箱"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              密码
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className="w-full px-4 py-3 rounded-xl border border-gray-300 focus:ring-2 focus:ring-blue-400 focus:outline-none"
              placeholder="输入密码"
            />
          </div>
          <button
            type="submit"
            className="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-medium transition"
          >
            登录
          </button>
        </form>
        <p className="text-center text-sm text-gray-500 mt-6">
          没有账号？{" "}
          <a href="/register" className="text-blue-600 hover:underline">
            注册
          </a>
        </p>
      </div>
    </div>
  );
}
