import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";
import { useLang } from "../contexts/LanguageContext";
import { useToast } from "../hooks/useToast";
import { getAnalytics, type AnalyticsData } from "../api/auth";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import ChangePasswordModal from "../components/ChangePasswordModal";
import "../styles/profile.css";

export default function ProfilePage() {
  const { user, isAuthenticated, isLoading, logout } = useAuth();
  const { t, lang } = useLang();
  const { showToast } = useToast();
  const navigate = useNavigate();
  const [analytics, setAnalytics] = useState<AnalyticsData | null>(null);
  const [cpOpen, setCpOpen] = useState(false);

  useEffect(() => {
    document.title = `Quikram — ${t("profile.title")}`;
    if (!isLoading && !isAuthenticated) navigate("/login", { replace: true });
  }, [isAuthenticated, isLoading, navigate, t]);

  useEffect(() => {
    if (isAuthenticated) {
      getAnalytics().then(({ data }) => setAnalytics(data)).catch(() => {});
    }
  }, [isAuthenticated]);

  if (isLoading)
    return (
      <div className="profile-page">
        <Navbar />
        <div className="profile-loading">{t("auth.loading")}</div>
        <Footer />
      </div>
    );
  if (!user) return null;

  const handleLogout = async () => {
    await logout();
    navigate("/login", { replace: true });
  };

  const createdDate = new Date(user.created_at);
  const daysSince = Math.floor(
    (Date.now() - createdDate.getTime()) / (1000 * 60 * 60 * 24),
  );
  const memberSince = createdDate.toLocaleDateString(lang === "ru" ? "ru-RU" : "en-US", {
    month: "long",
    year: "numeric",
  });
  const uploadsLimit = user.plan === "pro" || user.plan === "proai" ? -1 : 10;
  const uploadPercent =
    uploadsLimit === -1
      ? 0
      : Math.min(100, (user.uploads_this_month / uploadsLimit) * 100);

  const statItems = [
    { value: analytics?.total_tests ?? 0, label: t("profile.tests") },
    { value: analytics?.subjects_count ?? 0, label: t("profile.courses") },
    {
      value:
        analytics?.average_score != null
          ? `${Math.round(analytics.average_score)}%`
          : "—",
      label: t("profile.success"),
    },
    { value: daysSince > 0 ? daysSince.toString() : "—", label: t("profile.days") },
  ];

  return (
    <div className="profile-page">
      <Navbar />
      <main className="profile-content">
        <div className="profile-card">
          <div className="profile-banner" />

          <div className="profile-avatar-wrapper">
            <div className="profile-avatar">
              {user.name.charAt(0).toUpperCase()}
            </div>
          </div>

          <div className="profile-body">
            <div className="profile-head">
              <h1 className="profile-name">{user.name}</h1>
              <span className="profile-plan">
                {user.plan === "free" ? t("profile.free") : t("profile.premium")}
              </span>
            </div>
            <p className="profile-email">{user.email}</p>

            <div className="profile-stats">
              {statItems.map((s) => (
                <div key={s.label} className="profile-stat">
                  <span className="profile-stat-value">{s.value}</span>
                  <span className="profile-stat-label">{s.label}</span>
                </div>
              ))}
            </div>

            <div className="profile-details">
              <div className="profile-detail profile-detail-full">
                <span className="profile-detail-label">{t("profile.userId")}</span>
                <span className="profile-detail-value">{user.id}</span>
              </div>

              <div className="profile-detail">
                <span className="profile-detail-label">{t("profile.uploads")}</span>
                <span className="profile-detail-value">
                  {user.uploads_this_month}
                  {uploadsLimit === -1 ? " / ∞" : ` / ${uploadsLimit}`}
                </span>
              </div>

              <div className="profile-detail">
                <span className="profile-detail-label">{t("profile.memberSince")}</span>
                <span className="profile-detail-value">{memberSince}</span>
              </div>

              {uploadsLimit !== -1 && (
                <div className="profile-upload-bar-wrap">
                  <div className="profile-upload-bar">
                    <div
                      className="profile-upload-fill"
                      style={{ width: `${uploadPercent}%` }}
                    />
                  </div>
                </div>
              )}
            </div>

            <button className="profile-logout" onClick={handleLogout}>
              {t("profile.logout")}
            </button>

            <button className="profile-cp-btn" onClick={() => setCpOpen(true)}>
              {t("profile.changePass")}
            </button>
          </div>
        </div>
      </main>
      {cpOpen && <ChangePasswordModal onClose={() => setCpOpen(false)} onSuccess={(msg) => showToast(msg)} />}
      <Footer />
    </div>
  );
}
