import { useState } from "react";
import api from "../../services/api";

function parseMatrixInput(text: string): number[][] {
  // 支持以换行分隔行，逗号或空格分隔列
  const rows = text
    .trim()
    .split(/\r?\n/)
    .map((r) => r.trim())
    .filter((r) => r.length > 0)
    .map((r) =>
      r
        .split(/[ ,]+/)
        .filter((c) => c.length > 0)
        .map((c) => Number(c))
    );
  return rows;
}

export default function Matrix() {
  const [aText, setAText] = useState("1 2\n3 4");
  const [bText, setBText] = useState("5 6\n7 8");
  const [result, setResult] = useState<number[][] | null>(null);
  const [error, setError] = useState<string | null>(null);

  const buildPayload = () => {
    const A = parseMatrixInput(aText);
    const B = parseMatrixInput(bText);
    return { a: { rows: A.length, cols: A[0]?.length ?? 0, data: A }, b: { rows: B.length, cols: B[0]?.length ?? 0, data: B } };
  };

  const doOp = async (op: "addition" | "subtraction" | "multiplication") => {
    setError(null);
    setResult(null);
    const payload = buildPayload();

    // basic validation
    if (payload.a.rows === 0 || payload.a.cols === 0) {
      setError("矩阵 A 无效");
      return;
    }
    if (payload.b.rows === 0 || payload.b.cols === 0) {
      setError("矩阵 B 无效");
      return;
    }

    try {
      const res = await api.post(`/app/matrix/${op}`, payload);
      setResult(res.data.data);
    } catch (e: any) {
      console.error(e);
      setError(e.response?.data?.error || e.message || "请求失败");
    }
  };

  return (
    <div className="p-6">
      <h2 className="text-xl font-bold mb-4">矩阵计算</h2>

      <div className="grid grid-cols-2 gap-6 mb-4">
        <div>
          <h3 className="font-semibold mb-2">矩阵 A</h3>
          <textarea
            value={aText}
            onChange={(e) => setAText(e.target.value)}
            className="w-full h-40 p-2 border rounded"
          />
          <div className="text-sm text-gray-500 mt-2">格式示例：每行一行数据，空格或逗号分隔列，例如：<code>1 2\n3 4</code></div>
        </div>

        <div>
          <h3 className="font-semibold mb-2">矩阵 B</h3>
          <textarea
            value={bText}
            onChange={(e) => setBText(e.target.value)}
            className="w-full h-40 p-2 border rounded"
          />
          <div className="text-sm text-gray-500 mt-2">示例：<code>5 6\n7 8</code></div>
        </div>
      </div>

      <div className="flex gap-3 mb-4">
        <button onClick={() => doOp("addition")} className="px-4 py-2 bg-blue-600 text-white rounded">相加</button>
        <button onClick={() => doOp("subtraction")} className="px-4 py-2 bg-yellow-600 text-white rounded">相减</button>
        <button onClick={() => doOp("multiplication")} className="px-4 py-2 bg-green-600 text-white rounded">相乘</button>
      </div>

      {error && <div className="text-red-600 mb-4">错误：{error}</div>}

      {result && (
        <div>
          <div className="flex items-center justify-between mb-2">
            <h3 className="font-semibold">结果</h3>
            <div className="flex items-center gap-2">
              <button
                onClick={() => {
                  try {
                    navigator.clipboard.writeText(JSON.stringify(result));
                    alert("已复制结果到剪贴板");
                  } catch (e) {
                    console.error(e);
                  }
                }}
                className="px-3 py-1 bg-gray-200 rounded text-sm"
              >复制结果</button>
              <button
                onClick={() => {
                  const csv = result.map((r) => r.join(",")).join("\n");
                  const blob = new Blob([csv], { type: "text/csv;charset=utf-8;" });
                  const url = URL.createObjectURL(blob);
                  const a = document.createElement("a");
                  a.href = url;
                  a.download = "matrix-result.csv";
                  a.click();
                  URL.revokeObjectURL(url);
                }}
                className="px-3 py-1 bg-gray-200 rounded text-sm"
              >导出 CSV</button>
            </div>
          </div>

          <div className="inline-block p-3 bg-white border rounded">
            <div className="grid gap-1"
              style={{ gridTemplateColumns: `repeat(${result[0].length}, minmax(48px, 1fr))` }}>
              {result.map((row, i) =>
                row.map((cell, j) => (
                  <div key={`${i}-${j}`} className="border bg-gray-50 text-center p-2 text-sm">
                    {cell}
                  </div>
                ))
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
