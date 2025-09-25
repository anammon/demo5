// frontend/src/components/Layout.tsx
import type { ReactNode } from "react";

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen w-full bg-gradient-to-br from-blue-50 via-white to-blue-100 flex items-start justify-center py-10">
      <div className="w-full max-w-5xl bg-white rounded-2xl shadow-lg p-8">
        {children}
      </div>
    </div>
  );
}
