import { useCallback, useEffect, useState } from "react";
import { ArrowUpRight, CheckCircle2, CircleAlert, RefreshCw, Server, Star } from "lucide-react";
import { Button } from "./components/ui/button";
import { fetchHealth, fetchWoodfish, knockWoodfish } from "./lib/api";

const githubURL = "https://github.com/xiaolinbenben/fullstack-app-template";

type HealthState = "loading" | "online" | "offline";

export default function App() {
  const [health, setHealth] = useState<HealthState>("loading");
  const [checkedAt, setCheckedAt] = useState("");
  const [count, setCount] = useState<number | null>(null);
  const [knocking, setKnocking] = useState(false);
  const [knockError, setKnockError] = useState("");

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

  const loadCount = useCallback(async () => {
    try {
      const result = await fetchWoodfish();
      setCount(result.data.count);
      setKnockError("");
    } catch (error) {
      setKnockError(error instanceof Error ? error.message : "木鱼次数暂时读不到");
    }
  }, []);

  useEffect(() => {
    void checkHealth();
    void loadCount();
  }, [checkHealth, loadCount]);

  const knock = async () => {
    setKnocking(true);
    setKnockError("");
    try {
      const result = await knockWoodfish();
      setCount(result.data.count);
    } catch (error) {
      setKnockError(error instanceof Error ? error.message : "这一下没有记上");
    } finally {
      setKnocking(false);
    }
  };

  const isOnline = health === "online";

  return (
    <main className="min-h-screen bg-[#f7f8fa] text-slate-900">
      <header className="border-b border-slate-200/80 bg-white/85 backdrop-blur">
        <div className="mx-auto flex h-16 max-w-6xl items-center justify-between gap-3 px-4 sm:px-6">
          <a className="flex min-w-0 items-center gap-2.5 font-semibold tracking-tight" href="/">
            <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-slate-900 text-white">
              <Server size={18} />
            </span>
            <span className="truncate">Fullstack Template</span>
          </a>
          <nav className="flex shrink-0 items-center gap-3 text-sm text-slate-500">
            <a className="inline-flex items-center gap-1 whitespace-nowrap transition hover:text-slate-900" href={githubURL} target="_blank" rel="noreferrer">
              <Star size={15} />
              <span className="sr-only sm:not-sr-only">给项目点 Star</span>
            </a>
            <a className="inline-flex items-center gap-1 whitespace-nowrap transition hover:text-slate-900" href="/admin/">
              管理端 <ArrowUpRight size={15} className="shrink-0" />
            </a>
          </nav>
        </div>
      </header>

      <div className="mx-auto max-w-6xl px-6 py-16 lg:py-24">
        <section className="max-w-3xl">
          <p className="mb-5 text-xs font-semibold uppercase tracking-[0.22em] text-teal-600">Starter workspace</p>
          <h1 className="text-4xl font-semibold tracking-tight text-slate-950 sm:text-6xl">这是一个全栈模板。</h1>
          <p className="mt-6 max-w-2xl text-lg leading-8 text-slate-500">
            公共前端、管理端和 Go 服务打在同一个进程里。下面的木鱼次数存在 PostgreSQL 中，刷新之后还在。
          </p>
        </section>

        <section className="mt-12 grid gap-5 md:grid-cols-2" aria-label="模板入口">
          <div className="rounded-xl border border-slate-200 bg-white p-6 shadow-[0_10px_35px_rgba(15,23,42,0.04)]">
            <p className="text-sm font-medium text-slate-500">参观管理端</p>
            <h2 className="mt-2 text-2xl font-semibold tracking-tight">账号已经写在这里</h2>
            <dl className="mt-6 space-y-3 text-sm">
              <div className="flex items-center justify-between gap-4 border-b border-slate-100 pb-3">
                <dt className="text-slate-500">账号</dt>
                <dd className="font-mono text-slate-950">admin</dd>
              </div>
              <div className="flex items-center justify-between gap-4">
                <dt className="text-slate-500">密码</dt>
                <dd className="font-mono text-slate-950">admin123</dd>
              </div>
            </dl>
            <a className="mt-6 inline-flex items-center gap-1 text-sm font-medium text-slate-950" href="/admin/">
              打开管理端 <ArrowUpRight size={15} />
            </a>
          </div>

          <div className="rounded-xl border border-slate-200 bg-slate-950 p-6 text-slate-100">
            <p className="text-sm font-medium text-slate-400">木鱼</p>
            <p className="mt-6 font-mono text-5xl font-semibold tracking-tight">{count === null ? "—" : count}</p>
            <p className="mt-2 text-sm text-slate-400">一共被人敲过的次数</p>
            <Button className="mt-8 bg-white text-slate-950 hover:bg-slate-200" onClick={() => void knock()} disabled={knocking}>
              {knocking ? "记下这一下…" : "敲一下"}
            </Button>
            {knockError && <p className="mt-4 text-sm text-amber-300">{knockError}</p>}
          </div>
        </section>

        <section className="mt-5 grid gap-5 md:grid-cols-[1.35fr_1fr]" aria-label="服务状态">
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

          <div className="rounded-xl border border-slate-200 bg-white p-6">
            <p className="text-sm font-medium text-slate-500">Endpoints</p>
            <div className="mt-5 space-y-4 font-mono text-sm text-slate-700">
              <div className="flex items-center justify-between gap-4 border-b border-slate-100 pb-3"><span className="text-slate-400">woodfish</span><span>GET / POST</span></div>
              <div className="flex items-center justify-between gap-4 border-b border-slate-100 pb-3"><span className="text-slate-400">admin</span><span>/admin/</span></div>
              <div className="flex items-center justify-between gap-4"><span className="text-slate-400">health</span><span>/healthz</span></div>
            </div>
          </div>
        </section>
      </div>
    </main>
  );
}
