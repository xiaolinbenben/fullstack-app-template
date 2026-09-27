import { useCallback, useEffect, useState } from "react";
import { ArrowUpRight, CheckCircle2, CircleAlert, RefreshCw, Server } from "lucide-react";
import { Button } from "./components/ui/button";
import { fetchHealth } from "./lib/api";

type HealthState = "loading" | "online" | "offline";

export default function App() {
  const [health, setHealth] = useState<HealthState>("loading");
  const [checkedAt, setCheckedAt] = useState("");

  const checkHealth = useCallback(async () => {
    setHealth("loading");
    try {
      const result = await fetchHealth();
      setHealth(result.data.status === "ok" ? "online" : "offline");
      setCheckedAt(new Date().toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" }));
    } catch {
      setHealth("offline");
      setCheckedAt(new Date().toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" }));
    }
  }, []);

  useEffect(() => {
    void checkHealth();
  }, [checkHealth]);

  const isOnline = health === "online";

  return (
    <main className="min-h-screen bg-[#f7f8fa] text-slate-900">
      <header className="border-b border-slate-200/80 bg-white/85 backdrop-blur">
        <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
          <a className="flex items-center gap-3 font-semibold tracking-tight" href="/">
            <span className="flex h-9 w-9 items-center justify-center rounded-lg bg-slate-900 text-white">
              <Server size={18} />
            </span>
            Fullstack Template
          </a>
          <a className="inline-flex items-center gap-1 text-sm text-slate-500 transition hover:text-slate-900" href="/admin/">
            管理端 <ArrowUpRight size={15} />
          </a>
        </div>
      </header>

      <div className="mx-auto max-w-6xl px-6 py-16 lg:py-24">
        <section className="max-w-3xl">
          <p className="mb-5 text-xs font-semibold uppercase tracking-[0.22em] text-teal-600">Starter workspace</p>
          <h1 className="text-4xl font-semibold tracking-tight text-slate-950 sm:text-6xl">准备好开始构建。</h1>
          <p className="mt-6 max-w-2xl text-lg leading-8 text-slate-500">
            公共前端、管理端和 Gin 服务已经连接在同一个可部署的工程结构中。
          </p>
        </section>

        <section className="mt-16 grid gap-5 md:grid-cols-[1.35fr_1fr]" aria-label="服务状态">
          <div className="rounded-xl border border-slate-200 bg-white p-6 shadow-[0_10px_35px_rgba(15,23,42,0.04)]">
            <div className="flex items-start justify-between gap-4">
              <div>
                <p className="text-sm font-medium text-slate-500">Runtime status</p>
                <h2 className="mt-2 text-2xl font-semibold tracking-tight">后端服务</h2>
              </div>
              <span className={`inline-flex items-center gap-2 rounded-full px-3 py-1 text-xs font-medium ${isOnline ? "bg-emerald-50 text-emerald-700" : "bg-amber-50 text-amber-700"}`}>
                {isOnline ? <CheckCircle2 size={14} /> : <CircleAlert size={14} />}
                {health === "loading" ? "检查中" : isOnline ? "运行正常" : "暂不可用"}
              </span>
            </div>
            <div className="mt-10 flex items-center justify-between border-t border-slate-100 pt-5 text-sm text-slate-500">
              <span>{checkedAt ? `最近检查 ${checkedAt}` : "等待检查"}</span>
              <Button variant="outline" size="sm" onClick={() => void checkHealth()} disabled={health === "loading"}>
                <RefreshCw size={15} className={health === "loading" ? "animate-spin" : ""} />
                刷新
              </Button>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-slate-950 p-6 text-slate-100">
            <p className="text-sm font-medium text-slate-400">Endpoints</p>
            <div className="mt-5 space-y-4 font-mono text-sm">
              <div className="flex items-center justify-between gap-4 border-b border-white/10 pb-3"><span className="text-slate-400">public</span><span>/</span></div>
              <div className="flex items-center justify-between gap-4 border-b border-white/10 pb-3"><span className="text-slate-400">admin</span><span>/admin/</span></div>
              <div className="flex items-center justify-between gap-4"><span className="text-slate-400">health</span><span>/healthz</span></div>
            </div>
          </div>
        </section>
      </div>
    </main>
  );
}
