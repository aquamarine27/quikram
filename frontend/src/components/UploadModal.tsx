import { useState, useRef, useEffect, useCallback } from "react";
import { useLang } from "../contexts/LanguageContext";

interface Props {
  uploading: boolean;
  onUpload: (file: File) => void;
  onClose: () => void;
}

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

function getFileIcon(file: File): { label: string; cls: string } {
  const name = file.name.toLowerCase();
  const type = file.type;
  if (type === "application/pdf" || name.endsWith(".pdf")) return { label: "PDF", cls: "pdf" };
  if (type.startsWith("image/") || /\.(png|jpg|jpeg|gif|webp)$/.test(name)) return { label: "IMG", cls: "image" };
  if (type.includes("word") || /\.(doc|docx)$/.test(name)) return { label: "DOC", cls: "doc" };
  if (type.startsWith("text/") || /\.(txt|md|csv)$/.test(name)) return { label: "TXT", cls: "text" };
  return { label: "FILE", cls: "other" };
}

export default function UploadModal({ uploading, onUpload, onClose }: Props) {
  const { t } = useLang();
  const [file, setFile] = useState<File | null>(null);
  const [dragging, setDragging] = useState(false);
  const fileRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => { if (e.key === "Escape") onClose(); };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [onClose]);

  const backdropClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) onClose();
  };

  const onDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setDragging(false);
    const f = e.dataTransfer.files[0];
    if (f) setFile(f);
  }, []);

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const f = e.target.files?.[0];
    if (f) setFile(f);
    e.target.value = "";
  };

  const icon = file ? getFileIcon(file) : null;

  return (
    <div className="courses-modal-overlay" onClick={backdropClick}>
      <div className="courses-modal-backdrop" />
      <div className="courses-modal" style={{ perspective: "800px" }}>
        <button className="courses-modal-close" onClick={onClose}>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>

        <h2>{t("courses.modal.title")}</h2>

        {!file ? (
          <div
            className={`courses-modal-zone${dragging ? " dragging" : ""}`}
            onDragOver={(e) => { e.preventDefault(); setDragging(true); }}
            onDragLeave={() => setDragging(false)}
            onDrop={onDrop}
            onClick={() => fileRef.current?.click()}
          >
            <div className="courses-modal-zone-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="17 8 12 3 7 8"/>
                <line x1="12" y1="3" x2="12" y2="15"/>
              </svg>
            </div>
            <p>{t("courses.upload.zone")}</p>
            <span>{t("courses.upload.maxSize")}</span>
            <input ref={fileRef} type="file" onChange={onChange} accept=".pdf,.doc,.docx,.txt,.png,.jpg,.jpeg,.gif,.webp,.md,.csv" />
          </div>
        ) : (
          <>
            <div className="courses-modal-file">
              <div className={`courses-modal-file-icon ${icon!.cls}`}>{icon!.label}</div>
              <div className="courses-modal-file-info">
                <p className="courses-modal-file-name">{file.name}</p>
                <span className="courses-modal-file-size">{formatSize(file.size)}</span>
              </div>
              <button className="courses-modal-file-remove" onClick={() => setFile(null)}>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
              </button>
            </div>

            <button
              className="courses-modal-add-btn"
              onClick={() => { if (!uploading) onUpload(file); }}
              disabled={uploading}
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="17 8 12 3 7 8"/>
                <line x1="12" y1="3" x2="12" y2="15"/>
              </svg>
              {t("courses.modal.add")}
            </button>
          </>
        )}
      </div>
    </div>
  );
}
