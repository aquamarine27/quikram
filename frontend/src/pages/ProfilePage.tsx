import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";
import { useLang } from "../contexts/LanguageContext";
import ThemeToggle from "../components/ThemeToggle";
import LangSwitcher from "../components/LangSwitcher";
import "../styles/profile.css";

// ─── Component ───
export default function ProfilePage() {
  const { user, isAuthenticated, isLoading, logout } = useAuth();
  const { t } = useLang();
  const navigate = useNavigate();

  // ─── Effects ───
  useEffect(() => {
    document.title = `Quikram — ${t("profile.title")}`;
    if (!isLoading && !isAuthenticated) navigate("/login", { replace: true });
  }, [isAuthenticated, isLoading, navigate, t]);

  if (isLoading) return <div className="profile-loading">{t("auth.loading")}</div>;
  if (!user) return null;

  // ─── Handlers ───
  const handleLogout = async () => {
    await logout();
    navigate("/login", { replace: true });
  };

  // ─── Render ───
  return (
    <div className="profile-page">
      <div className="profile-card">
        <div className="profile-top">
          <ThemeToggle />
          <LangSwitcher />
        </div>

        <div className="profile-avatar">{user.name.charAt(0).toUpperCase()}</div>
        <h1 className="profile-name">{user.name}</h1>
        <p className="profile-email">{user.email}</p>
        <span className="profile-plan">{user.plan === "free" ? t("profile.free") : t("profile.premium")}</span>

        <div className="profile-divider">
          <div className="profile-field">
            <div className="profile-field-label">{t("profile.userId")}</div>
            <p className="profile-field-value">{user.id}</p>
          </div>
        </div>

        <button className="profile-logout" onClick={handleLogout}>
          {t("profile.logout")}
        </button>
      </div>
    </div>
  );
}
