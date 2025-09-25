// frontend/src/pages/apps/CreateApp.tsx
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { createApp } from "../../services/app";

export default function CreateApp() {
  const navigate = useNavigate();
  const [form, setForm] = useState({
    name: "",
    description: "",
    icon: "",
    type: "",
    tags: "",
    status: "published",
  });
  const [loading, setLoading] = useState(false);

  const onChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) =>
    setForm({ ...form, [e.target.name]: e.target.value });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const payload = {
        name: form.name.trim(),
        description: form.description.trim(),
        icon: form.icon.trim(),
        type: form.type.trim(),
        tags: form.tags.split(",").map((t) => t.trim()).filter(Boolean),
        status: form.status,
      };
      await createApp(payload);
      navigate("/home");
    } catch (err: any) {
      alert("创建失败: " + (err.response?.data?.error || err.message || "未知错误"));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="relative w-screen h-screen flex items-center justify-center overflow-hidden">
      <div className="absolute inset-0 bg-gradient-to-br from-blue-200 via-blue-100 to-white"></div>
      <div className="relative z-10 bg-white shadow-2xl rounded-2xl p-8 w-full max-w-2xl">
        <h2 className="text-2xl font-bold mb-4">创建新应用</h2>
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
            {loading ? "创建中..." : "创建应用"}
          </button>
        </form>
      </div>
    </div>
  );
}
