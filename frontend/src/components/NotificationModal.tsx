import { useRef, useEffect, type MouseEvent } from "react";
import { useLang } from "../contexts/LanguageContext";
import "../styles/notification-modal.css";

interface Props {
  onClose: () => void;
}

export default function NotificationModal({ onClose }: Props) {
  const { t } = useLang();
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) onClose();
    };
    document.addEventListener("mousedown", handler as any);
    return () => document.removeEventListener("mousedown", handler as any);
  }, [onClose]);

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    document.addEventListener("keydown", handler);
    return () => document.removeEventListener("keydown", handler);
  }, [onClose]);

  const handleOverlay = (e: MouseEvent) => {
    if (e.target === e.currentTarget) onClose();
  };

  return (
    <div className="notif-modal" ref={ref} onClick={handleOverlay}>
      <div className="notif-modal-card">
        <div className="notif-modal-arrow" />
        <div className="notif-modal-header">
          {t("notif.title")}
          <button className="modal-close-btn" onClick={onClose} aria-label="Close">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <line x1="18" y1="6" x2="6" y2="18" />
              <line x1="6" y1="6" x2="18" y2="18" />
            </svg>
          </button>
        </div>
        <div className="notif-modal-body">
          <div className="notif-placeholder">
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
              <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
              <path d="M13.73 21a2 2 0 0 1-3.46 0" />
            </svg>
            <p>{t("notif.empty")}</p>
          </div>
        </div>
      </div>
    </div>
  );
}
