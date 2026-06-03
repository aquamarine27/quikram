import { useEffect, useState, useRef } from "react";
import { useLang } from "../contexts/LanguageContext";
import { getSummary, createSummary, type Document, type Summary } from "../api/documents";

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

function getDocIcon(doc: Document): { label: string; cls: string } {
  const mt = doc.mime_type;
  const ext = doc.filename.split(".").pop()?.toLowerCase();
  if (mt === "application/pdf" || ext === "pdf") return { label: "PDF", cls: "pdf" };
  if (mt.startsWith("image/") || ["png", "jpg", "jpeg", "gif", "webp"].includes(ext || "")) return { label: "IMG", cls: "image" };
  if (mt.includes("word") || ["doc", "docx"].includes(ext || "")) return { label: "DOC", cls: "doc" };
  if (mt.startsWith("text/") || ["txt", "md", "csv"].includes(ext || "")) return { label: "TXT", cls: "text" };
  return { label: "FILE", cls: "other" };
}

interface Props {
  subjectId: string;
  doc: Document;
  onClose: () => void;
  generatingRef?: React.MutableRefObject<boolean>;
}

type State =
  | { type: "loading" }
  | { type: "processing" }
  | { type: "plan_restricted"; message: string }
  | { type: "ready"; summary: Summary }
  | { type: "placeholder" }
  | { type: "error"; message: string };

function renderSummaryBlocks(text: string): React.ReactNode[] {
  const blocks = text.split("\n\n");
  const result: React.ReactNode[] = [];
  blocks.forEach((block, i) => {
    const lines = block.split("\n").map(l => l.trim()).filter(l => l);
    if (lines.length === 0) return;
    if (lines.length > 1 && lines[0].length < 120) {
      result.push(<h3 key={`h${i}`} className="summary-section-header">{lines[0]}</h3>);
      result.push(<p key={`p${i}`} className="summary-section-text">{lines.slice(1).join("\n")}</p>);
    } else if (lines.length === 1 && lines[0].length < 120) {
      result.push(<h3 key={`h${i}`} className="summary-section-header">{lines[0]}</h3>);
    } else {
      result.push(<p key={`p${i}`} className="summary-section-text">{block}</p>);
    }
  });
  return result;
}

export default function SummaryModal({ subjectId, doc, onClose, generatingRef }: Props) {
  const { t } = useLang();
  const icon = getDocIcon(doc);
  const [state, setState] = useState<State>({ type: "loading" });
  const pollRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const [retryCount, setRetryCount] = useState(0);
  const [version, setVersion] = useState<"short" | "medium" | "long">("medium");
  const versions = ["short", "medium", "long"] as const;
  const switcherRef = useRef<HTMLDivElement>(null);
  const [indicator, setIndicator] = useState({ left: 0, width: 0 });

  useEffect(() => {
    if (!switcherRef.current) return;
    const active = switcherRef.current.querySelector(".summary-version-btn.active") as HTMLElement;
    if (active) setIndicator({ left: active.offsetLeft, width: active.offsetWidth });
  }, [version, state]);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => { if (e.key === "Escape") onClose(); };
    window.addEventListener("keydown", onKey);
    return () => {
      window.removeEventListener("keydown", onKey);
      if (pollRef.current) clearTimeout(pollRef.current);
      if (generatingRef) generatingRef.current = false;
    };
  }, [onClose, generatingRef]);

  useEffect(() => {
    document.body.style.overflow = "hidden";
    return () => { document.body.style.overflow = ""; };
  }, []);

  useEffect(() => {
    if (!generatingRef) return;
    if (state.type === "processing") {
      generatingRef.current = true;
    } else if (state.type === "ready" || state.type === "error" || state.type === "placeholder" || state.type === "plan_restricted") {
      generatingRef.current = false;
    }
  }, [state, generatingRef]);

  useEffect(() => {
    (async () => {
      setState({ type: "loading" });

      // 1 — try to fetch existing summary
      try {
        const s = await getSummary(subjectId, doc.id);
        if (s.content_short || s.content_medium || s.content_long) {
          setState({ type: "ready", summary: s });
          return;
        }
        // Summary exists but all empty — still generating, fall through to create
      } catch (err: any) {
        if (err?.response?.status !== 404) {
          setState({ type: "error", message: "Failed to load summary" });
          return;
        }
      }

      // 2 — 404 or empty summary, need to create/retry
      try {
        const result = await createSummary(subjectId, doc.id);
        if (result.content_short || result.content_medium || result.content_long) {
          setState({ type: "ready", summary: result as Summary });
          return;
        }
      } catch (err: any) {
        const msg = err?.response?.data?.error || "";
        if (msg.includes("not available") || err?.response?.status === 402) {
          setState({ type: "plan_restricted", message: t("courses.summary.planRestricted") || "AI summaries require Pro or ProAI plan" });
        } else {
          setState({ type: "placeholder" });
        }
        return;
      }

      // 3 — processing, check once after delay, then manual only
      setState({ type: "processing" });
      pollRef.current = setTimeout(async () => {
        try {
          const s = await getSummary(subjectId, doc.id);
          if (s.content_short || s.content_medium || s.content_long) {
            setState({ type: "ready", summary: s });
          }
        } catch {
          // stay on processing — user clicks check button to retry
        }
      }, 15000);
    })();
  }, [subjectId, doc.id, t, retryCount]);

  const retry = () => {
    if (pollRef.current) clearTimeout(pollRef.current);
    setRetryCount((c) => c + 1);
  };

  const backdropClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) onClose();
  };

  const uploaded = doc.created_at
    ? new Date(doc.created_at).toLocaleDateString(undefined, { day: "numeric", month: "short", year: "numeric" })
    : "";

  let body: React.ReactNode;

  switch (state.type) {
    case "loading":
      body = (
        <div className="summary-body-center">
          <div className="summary-body-spinner" />
        </div>
      );
      break;

    case "processing":
      body = (
        <div className="summary-body-center">
          <div className="summary-processing-ring">
            <svg viewBox="0 0 96 96" fill="none">
              <circle className="summary-processing-ring-bg" cx="48" cy="48" r="42" />
              <circle className="summary-processing-arc" cx="48" cy="48" r="42"
                strokeDasharray="264" strokeDashoffset="264" />
              <circle className="summary-processing-arc2" cx="48" cy="48" r="42"
                strokeDasharray="264" strokeDashoffset="264" />
            </svg>
            <div className="summary-processing-icon-inner">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/>
              </svg>
            </div>
          </div>
          <h2 className="summary-body-title">{t("courses.summary.processing") || "Обработка..."}</h2>
          <p className="summary-body-desc">{t("courses.summary.processingDesc") || "Анализируем ваш файл. Это может занять некоторое время."}</p>
          <div className="summary-processing-dots">
            <span className="summary-dot" />
            <span className="summary-dot" />
            <span className="summary-dot" />
          </div>
          <button className="summary-body-btn" onClick={retry}>
            {t("courses.summary.checkStatus") || "Проверить статус"}
          </button>
        </div>
      );
      break;

    case "plan_restricted":
      body = (
        <div className="summary-body-center">
          <div className="summary-body-lock-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
            </svg>
          </div>
          <h2 className="summary-body-title">{t("courses.summary.upgradeRequired") || "Upgrade Required"}</h2>
          <p className="summary-body-desc">{state.message}</p>
          <a href="/pricing" className="summary-body-btn primary" onClick={onClose}>
            {t("courses.summary.upgrade") || "View Plans"}
          </a>
        </div>
      );
      break;

    case "placeholder":
      body = (
        <div className="summary-body-center">
          <div className="summary-body-fail-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" strokeLinejoin="round">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
          </div>
          <h2 className="summary-body-title">{t("courses.summary.failed") || "Не удалось создать конспект"}</h2>
          <p className="summary-body-desc">{t("courses.summary.retryDesc") || "Попробуйте позже или загрузите файл заново"}</p>
          <button className="summary-body-btn" onClick={retry}>
            {t("courses.summary.retry") || "Повторить"}
          </button>
        </div>
      );
      break;

    case "error":
      body = (
        <div className="summary-body-center">
          <div className="summary-body-error-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
              <circle cx="12" cy="12" r="10"/>
              <line x1="15" y1="9" x2="9" y2="15"/>
              <line x1="9" y1="9" x2="15" y2="15"/>
            </svg>
          </div>
          <h2 className="summary-body-title">{t("courses.summary.error") || "Something went wrong"}</h2>
          <p className="summary-body-desc">{state.message}</p>
          <button className="summary-body-btn" onClick={retry}>
            {t("courses.summary.retry") || "Try again"}
          </button>
        </div>
      );
      break;

    case "ready":
      const contentMap = {
        short: state.summary.content_short,
        medium: state.summary.content_medium,
        long: state.summary.content_long,
      };
      const currentContent = contentMap[version];
      const partialLoading = !currentContent && (contentMap.short || contentMap.medium || contentMap.long);
      body = (
        <div className="summary-body-content">
          <div className="summary-content-sizer">
            <div className="summary-content-header">
              <div className="summary-content-header-left">
                <div className={`summary-content-header-icon ${icon.cls}`}>{icon.label}</div>
                <div>
                  <h2 className="summary-content-title">{t("courses.summary.title")}</h2>
                </div>
              </div>
              <div className="summary-content-header-right">
                <div className="summary-content-meta-item">
                  <span className="summary-content-meta-label">{t("courses.summary.uploadedAt") || "Uploaded"}</span>
                  <span className="summary-content-meta-value">{uploaded}</span>
                </div>
                <div className="summary-content-meta-item">
                  <span className="summary-content-meta-label">{t("courses.summary.size") || "Size"}</span>
                  <span className="summary-content-meta-value">{formatSize(doc.file_size)}</span>
                </div>
              </div>
            </div>
            <div className="summary-version-switcher" ref={switcherRef} onMouseMove={(e) => { const r = e.currentTarget.getBoundingClientRect(); e.currentTarget.style.setProperty("--vx", `${((e.clientX - r.left) / r.width) * 100}%`); }}>
              <div className="summary-version-indicator" style={{ left: indicator.left, width: indicator.width }} />
              {versions.map((v) => (
                <button key={v} className={`summary-version-btn ${version === v ? "active" : ""}`}
                  onClick={() => setVersion(v)}>
                  {v === "short" ? (t("courses.summary.short") || "Краткий") : v === "medium" ? (t("courses.summary.medium") || "Средний") : (t("courses.summary.long") || "Длинный")}
                </button>
              ))}
            </div>
            <div className="summary-content-divider" />
            <div className="summary-content-text" key={version}>
              <div className="summary-content-text-inner">
                {partialLoading ? (t("courses.summary.generating") || "еще генерируется...") : currentContent ? renderSummaryBlocks(currentContent) : (t("courses.summary.emptyContent") || "No content")}
              </div>
            </div>
          </div>
        </div>
      );
      break;
  }

  return (
    <div className="summary-overlay" onClick={backdropClick}>
      <div className="summary-backdrop" />
      <div className="summary-modal-full">
        <div className="summary-modal-full-header">
          <div className="summary-modal-full-header-left">
            <div className={`summary-mini-icon ${icon.cls}`}>{icon.label}</div>
            <div className="summary-mini-info">
              <span className="summary-mini-name">{doc.filename}</span>
              <span className="summary-mini-meta">{formatSize(doc.file_size)}</span>
            </div>
          </div>
          <button className="summary-close-btn" onClick={onClose}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div className="summary-modal-full-body">
          {body}
        </div>
      </div>
    </div>
  );
}
