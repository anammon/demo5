// frontend/src/pages/apps/Translator.tsx
import { useState } from "react";
import Layout from "../../components/Layout";
import api from "../../services/api";

export default function Translator() {
  const [inputText, setInputText] = useState("");
  const [result, setResult] = useState("");

  const handleTranslate = async () => {
    if (!inputText.trim()) return;
    try {
      // 如果你已经实现了后端 /translate 或 /api/translate，请改成正确路径
      const res = await api.post("/translate", { texts: [inputText], targetLang: "en" });
      // 假设后端返回 { translations: ["..."] }
      if (res?.data?.translations) setResult(res.data.translations[0]);
      else setResult("EN: " + inputText); // fallback
    } catch (err) {
      // fallback local mock
      setResult("EN: " + inputText);
    }
  };

  return (
    <Layout>
      <div className="max-w-3xl mx-auto">
        <h2 className="text-2xl font-bold mb-4">翻译器</h2>

        <textarea
          rows={6}
          value={inputText}
          onChange={(e) => setInputText(e.target.value)}
          className="w-full p-3 border rounded mb-4"
          placeholder="输入要翻译的文本或粘贴网页 OCR 结果"
        />

        <button onClick={handleTranslate} className="px-6 py-2 bg-green-600 text-white rounded mb-4">
          翻译
        </button>

        {result && (
          <div className="p-4 bg-gray-50 rounded">
            <h3 className="font-semibold mb-2">翻译结果</h3>
            <p>{result}</p>
          </div>
        )}
      </div>
    </Layout>
  );
}
