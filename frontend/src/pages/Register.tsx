import { useState } from "react";
import api from "../services/api";
import { useNavigate } from "react-router-dom";

export default function Register() {
  const [form, setForm] = useState({
    account: "",
    password: "",
    email: "",
    name: "",
    birthday: "",
  });
  const navigate = useNavigate();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await api.post("/user/register", form);

      alert("注册成功！欢迎 " + res.data.user.name);
      navigate("/");
    } catch (err: any) {
      alert("注册失败：" + (err.response?.data?.error || "未知错误"));
    }
  };

  return (
    <div className="w-screen h-screen flex items-center justify-center bg-gradient-to-br from-blue-100 via-blue-50 to-white">
      <div className="bg-white rounded-2xl shadow-2xl p-12 w-full max-w-md">
        <h1 className="text-3xl font-bold text-center text-blue-600 mb-8">
          注册新账号
        </h1>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              账户
            </label>
            <input
              type="text"
              name="account"
              value={form.account}
              onChange={handleChange}
              className="w-full px-4 py-3 border rounded-xl focus:ring-2 focus:ring-blue-400 focus:outline-none"
              placeholder="请输入账户名"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              密码
            </label>
            <input
              type="password"
              name="password"
              value={form.password}
              onChange={handleChange}
              className="w-full px-4 py-3 border rounded-xl focus:ring-2 focus:ring-blue-400 focus:outline-none"
              placeholder="请输入密码"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              邮箱
            </label>
            <input
              type="email"
              name="email"
              value={form.email}
              onChange={handleChange}
              className="w-full px-4 py-3 border rounded-xl focus:ring-2 focus:ring-blue-400 focus:outline-none"
              placeholder="请输入邮箱"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              姓名
            </label>
            <input
              type="text"
              name="name"
              value={form.name}
              onChange={handleChange}
              className="w-full px-4 py-3 border rounded-xl focus:ring-2 focus:ring-blue-400 focus:outline-none"
              placeholder="请输入姓名"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              生日
            </label>
            <input
              type="date"
              name="birthday"
              value={form.birthday}
              onChange={handleChange}
              className="w-full px-4 py-3 border rounded-xl focus:ring-2 focus:ring-blue-400 focus:outline-none"
              required
            />
          </div>
          <button
            type="submit"
            className="w-full py-3 px-4 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition"
          >
            注册
          </button>
        </form>
        <p className="mt-6 text-center text-sm text-gray-600">
          已有账号？{" "}
          <a href="/" className="text-blue-600 hover:underline">
            去登录
          </a>
        </p>
      </div>
    </div>
  );
}
