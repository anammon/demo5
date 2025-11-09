import { useEffect, useState, useRef } from "react";
import api from "../../services/api";

export default function Bottle() {
  const [thrownBottles, setThrownBottles] = useState<any[]>([]);
  const [pickedBottles, setPickedBottles] = useState<any[]>([]);
  const [message, setMessage] = useState("");
  const [isAnonymous, setIsAnonymous] = useState(true);
  
  // 使用 useRef 来获取最新的 isAnonymous 值
  const isAnonymousRef = useRef(isAnonymous);
  
  // 同步 ref 和 state
  useEffect(() => {
    isAnonymousRef.current = isAnonymous;
  }, [isAnonymous]);

  // 扔瓶子 - 使用 ref 获取最新值
  const throwBottle = async () => {
    try {
      const currentIsAnonymous = isAnonymousRef.current;
      
      const requestData = {
        content: message,
        is_anonymous: currentIsAnonymous
      };
      
      console.log("🎯 准备发送的数据:", requestData);
      console.log("🎯 当前匿名状态:", currentIsAnonymous);
      
      const response = await api.post("/app/bottle", requestData);
      
      console.log("✅ 服务器响应:", response.data);
      alert(`瓶子扔出去啦！${currentIsAnonymous ? "（匿名）" : "（显示身份）"}`);
      
      setMessage("");
      load();
    } catch (err: any) {
      console.error("❌ 扔瓶子失败", err);
      alert("扔瓶子失败，请重试");
    }
  };

  // 捡瓶子
  const pickBottle = async () => {
    try {
      const res = await api.get("/app/bottle/pick");
      if (res.data) {
        if (res.data.is_system) {
          alert("系统消息：" + res.data.content);
        } else {
          let alertMsg = "你捡到了：" + res.data.content;
          if (!res.data.is_anonymous && res.data.throw_user_info) {
            alertMsg += `\n\n👤 来自用户：${res.data.throw_user_info}`;
          } else {
            alertMsg += "\n\n🎭 匿名瓶子";
          }
          alert(alertMsg);
        }
      }
      load();
    } catch (err: any) {
      console.error("捡瓶子失败", err);
      if (err.response?.data?.error) {
        alert(err.response.data.error);
      } else {
        alert("捡瓶子失败，请重试");
      }
    }
  };

  // 加载数据
  const load = async () => {
    try {
      const res1 = await api.get("/app/bottle/my/thrown");
      setThrownBottles(Array.isArray(res1.data) ? res1.data : []);

      const res2 = await api.get("/app/bottle/my/picked");
      setPickedBottles(Array.isArray(res2.data) ? res2.data : []);
    } catch (err) {
      console.error("加载瓶子失败", err);
      setThrownBottles([]);
      setPickedBottles([]);
    }
  };

  useEffect(() => {
    load();
  }, []);

  return (
    <div className="p-6">
      <h2 className="text-xl font-bold mb-4">漂流瓶</h2>

      <div className="mb-6">
        <div className="flex gap-2 mb-3">
          <input
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="写下你的心里话..."
            className="flex-1 border px-3 py-2 rounded-md"
            maxLength={600}
          />
          <button
            onClick={throwBottle}
            disabled={!message.trim()}
            className="px-4 py-2 bg-blue-600 text-white rounded-md disabled:bg-gray-400"
          >
            扔出去
          </button>
          <button
            onClick={pickBottle}
            className="px-4 py-2 bg-green-600 text-white rounded-md"
          >
            捡瓶子
          </button>
        </div>
        
        {/* 匿名选择 */}
        <div className="flex items-center gap-3 bg-gray-50 p-3 rounded-lg border">
          <label className="flex items-center gap-2 cursor-pointer">
            <input 
              type="checkbox" 
              checked={isAnonymous}
              onChange={(e) => {
                const newValue = e.target.checked;
                console.log("🔄 设置匿名状态为:", newValue);
                setIsAnonymous(newValue);
              }}
              className="w-5 h-5 text-blue-600 rounded focus:ring-blue-500"
            />
            <span className="font-medium text-gray-800">匿名发送</span>
          </label>
          <span className="text-sm text-gray-600">
            {isAnonymous ? "🎭 别人不会知道是你写的" : "👤 捡到瓶子的人可以看到你的身份"}
          </span>
        </div>
        
        {/* 显示当前状态 */}
        <div className="mt-2 text-sm font-medium">
          当前模式: <span className={isAnonymous ? "text-gray-600" : "text-blue-600"}>
            {isAnonymous ? "🔒 匿名模式" : "🔓 公开模式"}
          </span>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-6">
        <div>
          <h3 className="font-semibold mb-3">我扔的瓶子 ({thrownBottles.length})</h3>
          <ul className="space-y-3">
            {thrownBottles.map((bottle, idx) => (
              <li key={idx} className="p-3 border rounded-md bg-white shadow-sm">
                <div className="flex justify-between items-center mb-2">
                  <span className={`text-xs px-2 py-1 rounded ${
                    bottle.is_anonymous 
                      ? 'bg-gray-100 text-gray-600 border border-gray-200' 
                      : 'bg-blue-100 text-blue-700 border border-blue-200'
                  }`}>
                    {bottle.is_anonymous ? "🎭 匿名" : "👤 公开身份"}
                  </span>
                  <span className={`text-xs px-2 py-1 rounded ${
                    bottle.is_picked 
                      ? 'bg-green-100 text-green-700 border border-green-200' 
                      : 'bg-yellow-100 text-yellow-700 border border-yellow-200'
                  }`}>
                    {bottle.is_picked ? "✅ 已被捡" : "⏳ 漂流中"}
                  </span>
                </div>
                <div className="text-gray-800 mb-2">{bottle.content}</div>
                <div className="text-xs text-gray-500">
                  {new Date(bottle.created_at).toLocaleString()}
                </div>
              </li>
            ))}
          </ul>
        </div>

        <div>
          <h3 className="font-semibold mb-3">我捡到的瓶子 ({pickedBottles.length})</h3>
          <ul className="space-y-3">
            {pickedBottles.map((bottle, idx) => (
              <li key={idx} className="p-3 border rounded-md bg-white shadow-sm">
                <div className="flex justify-between items-center mb-2">
                  <span className={`text-xs px-2 py-1 rounded ${
                    bottle.is_anonymous 
                      ? 'bg-gray-100 text-gray-600 border border-gray-200' 
                      : 'bg-green-100 text-green-700 border border-green-200'
                  }`}>
                    {bottle.is_anonymous ? "🎭 匿名瓶子" : "👤 实名瓶子"}
                  </span>
                  {!bottle.is_anonymous && (
                    <span className="text-xs bg-orange-100 text-orange-700 px-2 py-1 rounded border border-orange-200">
                      用户ID: {bottle.throw_user_id}
                    </span>
                  )}
                </div>
                <div className="text-gray-800 mb-2">{bottle.content}</div>
                <div className="text-xs text-gray-500">
                  {new Date(bottle.created_at).toLocaleString()}
                </div>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
  );
}