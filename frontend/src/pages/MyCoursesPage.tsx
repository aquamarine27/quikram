import { useEffect, useState, useCallback, useContext, useRef } from "react";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import UploadModal from "../components/UploadModal";
import SummaryModal from "../components/SummaryModal";
import TestListModal from "../components/TestListModal";
import { useLang } from "../contexts/LanguageContext";
import { AuthContext } from "../contexts/AuthContext";
import { ToastContext } from "../contexts/ToastContext";
import { getSubjects, createSubject } from "../api/subjects";
import { getDocuments, uploadDocument, deleteDocument, type Document } from "../api/documents";
import "../styles/courses.css";

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

export default function MyCoursesPage() {
  const { t } = useLang();
  const { user } = useContext(AuthContext)!;
  const { showToast } = useContext(ToastContext)!;

  const [subjectId, setSubjectId] = useState<string | null>(null);
  const [docs, setDocs] = useState<Document[]>([]);
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [summaryDoc, setSummaryDoc] = useState<Document | null>(null);
  const [testDoc, setTestDoc] = useState<Document | null>(null);
  const generatingRef = useRef(false);

  useEffect(() => { document.title = `Quikram — ${t("courses.mine.title")}`; }, [t]);

  const init = useCallback(async () => {
    setLoading(true);
    try {
      const list = await getSubjects();
      const defaultTitle = t("courses.mine.title");
      let my = list.find((s) => s.title === "My Materials" || s.title === "Мои материалы" || s.title === defaultTitle) || list[0];
      if (!my) {
        try {
          my = await createSubject(defaultTitle);
        } catch (e: any) {
          showToast(e?.response?.data?.error || "Failed to create subject", "error");
        }
      }
      if (my) {
        setSubjectId(my.id);
        const d = await getDocuments(my.id);
        setDocs(d);
      }
    } catch {
      showToast("Failed to load courses", "error");
    } finally {
      setLoading(false);
    }
  }, [t, showToast]);

  useEffect(() => {
    if (user) init();
    else setLoading(false);
  }, [user, init]);



  const onUpload = useCallback(async (file: File) => {
    if (!subjectId) return;
    setUploading(true);
    try {
      const doc = await uploadDocument(subjectId, file);
      setDocs((prev) => [doc, ...prev]);
      setModalOpen(false);
      showToast(t("courses.upload.success") || "Uploaded", "success");
    } catch {
      showToast(t("courses.upload.error") || "Upload failed", "error");
    } finally {
      setUploading(false);
    }
  }, [subjectId, t, showToast]);

  const handleDelete = useCallback(async (docId: string) => {
    if (!subjectId) return;
    try {
      await deleteDocument(subjectId, docId);
      setDocs((prev) => prev.filter((d) => d.id !== docId));
    } catch {
      showToast("Delete failed", "error");
    }
  }, [subjectId, showToast]);

  const statusLabel = (s: string): string => {
    const key = `courses.document.${s}`;
    return t(key) || s;
  };

  return (
    <div className="courses-page">
      <Navbar />
      <div className="courses-content">
        <div className="courses-mine-top">
          <h1>{t("courses.mine.title")}</h1>
          {docs.length > 0 && (
            <button className="courses-mine-create" onClick={() => setModalOpen(true)}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <line x1="12" y1="5" x2="12" y2="19"/>
                <line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
              {t("courses.mine.uploadNew")}
            </button>
          )}
        </div>

        {loading ? (
          <div className="courses-spinner">
            <div className="courses-spinner-ring" />
          </div>
        ) : docs.length === 0 ? (
          <div className="courses-mine-empty">
            <div className="courses-mine-empty-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14 2 14 8 20 8"/>
                <line x1="12" y1="18" x2="12" y2="12"/>
                <line x1="9" y1="15" x2="15" y2="15"/>
              </svg>
            </div>
            <h2>{t("courses.mine.empty")}</h2>
            <p>{t("courses.mine.emptyDesc")}</p>
            <button className="courses-mine-create" onClick={() => setModalOpen(true)}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <line x1="12" y1="5" x2="12" y2="19"/>
                <line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
              {t("courses.mine.create")}
            </button>
          </div>
        ) : (
          <div className="courses-doc-grid">
              {docs.map((doc) => {
              const icon = getDocIcon(doc);
              const uploadedDate = doc.created_at
                ? new Date(doc.created_at).toLocaleDateString(undefined, { day: "numeric", month: "short", year: "numeric" })
                : "";
              return (
                <div key={doc.id} className="courses-doc-card" onMouseMove={(e) => { const r = e.currentTarget.getBoundingClientRect(); e.currentTarget.style.setProperty("--mx", `${((e.clientX - r.left) / r.width) * 100}%`); e.currentTarget.style.setProperty("--my", `${((e.clientY - r.top) / r.height) * 100}%`); }}>
                  <div className="courses-doc-top">
                    <div className={`courses-doc-icon ${icon.cls}`}>{icon.label}</div>
                    <div className="courses-doc-info">
                      <p className="courses-doc-name">{doc.filename}</p>
                      <div className="courses-doc-meta">
                        <span>{formatSize(doc.file_size)}</span>
                        <span className="courses-doc-sep">·</span>
                        <span>{uploadedDate}</span>
                      </div>
                    </div>
                    <div className="courses-doc-top-right">
                      <button className="courses-doc-delete" onClick={(e) => { e.stopPropagation(); handleDelete(doc.id); }} title="Delete">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                          <polyline points="3 6 5 6 21 6"/>
                          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                        </svg>
                      </button>
                      <span className={`courses-doc-status ${doc.status}`}>{statusLabel(doc.status)}</span>
                    </div>
                  </div>
                  <div className="courses-doc-test-count">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                      <polyline points="9 11 12 14 22 4"/>
                      <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/>
                    </svg>
                    <span>0 {t("courses.mine.tests") || "tests"}</span>
                  </div>
                  <div className="courses-doc-actions">
                    <button className="courses-doc-action summary-btn" onClick={(e) => { e.stopPropagation(); if (generatingRef.current) { showToast("Сейчас генерируется конспект для другого файла, подождите", "error"); return; } setSummaryDoc(doc); }}>
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                        <polyline points="14 2 14 8 20 8"/>
                        <line x1="16" y1="13" x2="8" y2="13"/>
                        <line x1="16" y1="17" x2="8" y2="17"/>
                      </svg>
                      {t("courses.modal.summary")}
                    </button>
                    <button className="courses-doc-action test-btn" onClick={(e) => { e.stopPropagation(); setTestDoc(doc); }}>
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <polyline points="9 11 12 14 22 4"/>
                        <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/>
                      </svg>
                      {t("courses.modal.test")}
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {summaryDoc && subjectId && (
        <SummaryModal subjectId={subjectId} doc={summaryDoc} onClose={() => setSummaryDoc(null)} generatingRef={generatingRef} />
      )}

      {testDoc && (
        <TestListModal doc={testDoc} onClose={() => setTestDoc(null)} />
      )}

      {modalOpen && (
        <UploadModal
          uploading={uploading}
          onUpload={onUpload}
          onClose={() => setModalOpen(false)}
        />
      )}

      <Footer />
    </div>
  );
}
