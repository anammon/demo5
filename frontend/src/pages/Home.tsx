import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

export default function Home() {
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) navigate("/login");
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/login");
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-blue-50 via-blue-100 to-white">
      <h1 className="text-3xl font-bold text-blue-700 mb-6">欢迎使用 AssistApp</h1>
      <p className="mb-4 text-gray-700">你已经成功登录！</p>
      <button
        onClick={handleLogout}
        className="py-2 px-6 bg-red-500 text-white rounded-lg hover:bg-red-600 transition"
      >
        退出登录
      </button>
    </div>
  );
}
