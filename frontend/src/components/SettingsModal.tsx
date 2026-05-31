import { useRef, useEffect, type MouseEvent } from "react";
import { useLang } from "../contexts/LanguageContext";
import ThemeToggle from "./ThemeToggle";
import LangSwitcher from "./LangSwitcher";
import "../styles/settings-modal.css";

interface Props {
  onClose: () => void;
}

export default function SettingsModal({ onClose }: Props) {
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
    <div className="settings-modal" ref={ref} onClick={handleOverlay}>
      <div className="settings-modal-card">
        <div className="settings-modal-arrow" />
        <div className="settings-modal-header">
          {t("settings.title")}
          <button className="modal-close-btn" onClick={onClose} aria-label="Close">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <line x1="18" y1="6" x2="6" y2="18" />
              <line x1="6" y1="6" x2="18" y2="18" />
            </svg>
          </button>
        </div>
        <div className="settings-modal-body">
          <div className="settings-row">
            <span className="settings-label">{t("settings.theme")}</span>
            <ThemeToggle />
          </div>
          <div className="settings-row">
            <span className="settings-label">{t("settings.language")}</span>
            <LangSwitcher />
          </div>
        </div>
      </div>
    </div>
  );
}
