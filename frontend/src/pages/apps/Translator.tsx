import { useState } from "react";

export default function Translator() {
  const [inputText, setInputText] = useState("");
  const [result, setResult] = useState("");

  const handleTranslate = async () => {
    if (!inputText.trim()) return;

    try {
      // mock 模拟翻译
      setResult("EN: " + inputText);
    } catch (err) {
      setResult("翻译失败，请检查 API");
    }
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-green-50 via-green-100 to-white">
      <div className="bg-white shadow-2xl rounded-2xl p-8 w-full max-w-lg">
        <h1 className="text-2xl font-bold text-center text-green-600 mb-6">
          AssistApp 翻译器
        </h1>
        <textarea
          className="w-full p-3 border rounded-xl focus:ring-2 focus:ring-green-400 mb-4"
          rows={5}
          value={inputText}
          onChange={(e) => setInputText(e.target.value)}
          placeholder="请输入要翻译的文本..."
        />
        <button
          onClick={handleTranslate}
          className="w-full py-2 px-4 bg-green-600 text-white rounded-xl hover:bg-green-700 transition"
        >
          翻译
        </button>
        {result && (
          <div className="mt-6 p-4 bg-gray-50 rounded-xl border">
            <h2 className="text-lg font-semibold text-gray-700 mb-2">翻译结果：</h2>
            <p className="text-gray-800">{result}</p>
          </div>
        )}
      </div>
    </div>
  );
}
