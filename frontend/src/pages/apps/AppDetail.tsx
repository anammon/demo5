// frontend/src/pages/apps/AppDetail.tsx
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import Layout from "../../components/Layout";
import { fetchAppById, getAppLikes, likeApp } from "../../services/app";

export default function AppDetail() {
  const { id } = useParams<{ id: string }>();
  const [app, setApp] = useState<any | null>(null);
  const [likes, setLikes] = useState<number>(0);

  useEffect(() => {
    if (!id) return;
    fetchAppById(id).then((d) => setApp(d)).catch(() => {});
    getAppLikes(id).then((r: any) => setLikes(Number(r.likes || 0))).catch(()=>{});
  }, [id]);

  const handleLike = async () => {
    if (!id) return;
    await likeApp(id);
    setLikes((s) => s + 1);
  };

  if (!app) return <Layout>加载中...</Layout>;

  return (
    <Layout>
      <div className="max-w-3xl mx-auto text-center">
        {app.icon ? (
          <img src={app.icon} alt={app.name} className="w-24 h-24 mx-auto mb-4" />
        ) : (
          <div className="w-24 h-24 bg-gray-200 rounded-full mx-auto mb-4" />
        )}
        <h1 className="text-3xl font-bold mb-2">{app.name}</h1>
        <p className="text-gray-600 mb-6">{app.description}</p>

        <button onClick={handleLike} className="px-6 py-2 bg-indigo-600 text-white rounded-lg">
          👍 点赞 ({likes})
        </button>
      </div>
    </Layout>
  );
}
