import { useEffect, useState, type MouseEvent } from "react";
import { useLang } from "../contexts/LanguageContext";
import { changePassword } from "../api/auth";
import "../styles/settings-modal.css";

interface Props {
  onClose: () => void;
  onSuccess: (msg: string) => void;
}

export default function ChangePasswordModal({ onClose, onSuccess }: Props) {
  const { t } = useLang();
  const [oldPass, setOldPass] = useState("");
  const [newPass, setNewPass] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

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

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    if (!oldPass || !newPass) { setError(t("val.invalidEmail")); return; }
    if (newPass.length < 6) { setError(t("val.minChars")); return; }
    setLoading(true);
    try {
      await changePassword({ old_password: oldPass, new_password: newPass });
      onSuccess(t("profile.passChanged"));
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : t("toast.error"));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="cp-overlay" onClick={handleOverlay}>
      <div className="settings-modal-card" onClick={(e) => e.stopPropagation()}>
        <div className="settings-modal-header">
          {t("profile.changePassTitle")}
          <button className="modal-close-btn" onClick={onClose} aria-label="Close">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <line x1="18" y1="6" x2="6" y2="18" />
              <line x1="6" y1="6" x2="18" y2="18" />
            </svg>
          </button>
        </div>
        <form className="settings-modal-body" onSubmit={handleSubmit}>
          <div className="cp-field">
            <label className="cp-label">{t("profile.oldPass")}</label>
            <input
              className="cp-input"
              type="password"
              value={oldPass}
              onChange={(e) => setOldPass(e.target.value)}
              autoComplete="current-password"
            />
          </div>
          <div className="cp-field">
            <label className="cp-label">{t("profile.newPass")}</label>
            <input
              className="cp-input"
              type="password"
              value={newPass}
              onChange={(e) => setNewPass(e.target.value)}
              autoComplete="new-password"
            />
          </div>
          {error && <p className="cp-error">{error}</p>}
          <button className="cp-submit" type="submit" disabled={loading}>
            {loading ? t("auth.loading") : t("profile.changePass")}
          </button>
        </form>
      </div>
    </div>
  );
}
