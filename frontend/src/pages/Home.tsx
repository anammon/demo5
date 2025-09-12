import { useEffect } from "react";
import { useNavigate, Link } from "react-router-dom";

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
    <div className="relative min-h-screen w-screen overflow-hidden">
      {/* 渐变背景填满 */}
      <div className="absolute inset-0 bg-gradient-to-br from-blue-50 via-blue-100 to-white"></div>
      {/* 顶部导航 */}
      <header className="relative z-10 w-full flex justify-between items-center px-8 py-4 bg-white shadow-md">
        <h1 className="text-xl font-bold text-blue-700">AssistApp 应用广场</h1>
        <button
          onClick={handleLogout}
          className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition"
        >
          退出登录
        </button>
      </header>

      {/* 应用网格内容居中显示 */}
      <main className="relative z-10 flex-1 w-full flex justify-center items-start pt-8">
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
          {/* 翻译器卡片 */}
          <Link
            to="/apps/translator"
            className="bg-white rounded-xl shadow-lg flex flex-col justify-center items-center w-48 h-48 mx-auto hover:shadow-2xl transition"
          >
            <div className="w-16 h-16 bg-green-100 flex items-center justify-center rounded-full mb-4">
              🌐
            </div>
            <h2 className="text-lg font-semibold text-gray-800">翻译器</h2>
            <p className="text-sm text-gray-500 mt-2 text-center">
              实时翻译文本，保留原始布局。
            </p>
          </Link>

          {/* 预留应用卡片 */}
          <div className="bg-gray-100 rounded-xl flex flex-col justify-center items-center w-48 h-48 mx-auto text-gray-400">
            <div className="w-16 h-16 bg-gray-200 flex items-center justify-center rounded-full mb-4">
              🚧
            </div>
            <h2 className="text-lg font-semibold">更多应用</h2>
            <p className="text-sm mt-2 text-center">敬请期待...</p>
          </div>

          {/* 你可以继续加更多应用 */}
        </div>
      </main>
    </div>
  );
}
