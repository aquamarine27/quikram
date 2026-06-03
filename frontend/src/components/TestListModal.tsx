import { useEffect } from "react";
import { useLang } from "../contexts/LanguageContext";
import { type Document } from "../api/documents";

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
  doc: Document;
  onClose: () => void;
}

export default function TestListModal({ doc, onClose }: Props) {
  const { t } = useLang();
  const icon = getDocIcon(doc);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => { if (e.key === "Escape") onClose(); };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [onClose]);

  useEffect(() => {
    document.body.style.overflow = "hidden";
    return () => { document.body.style.overflow = ""; };
  }, []);

  const backdropClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) onClose();
  };

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
        <div className="testlist-body">
          <div className="testlist-empty">
            <div className="testlist-empty-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" strokeLinejoin="round">
                <polyline points="9 11 12 14 22 4"/>
                <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/>
              </svg>
            </div>
            <h2 className="testlist-empty-title">{t("courses.tests.empty") || "No tests yet"}</h2>
            <p className="testlist-empty-desc">{t("courses.tests.emptyDesc") || "Create your first test based on this material."}</p>
          </div>
        </div>
      </div>
    </div>
  );
}
