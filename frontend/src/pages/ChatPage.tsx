import { useEffect } from "react";
import { Link } from "react-router-dom";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import { useAuth } from "../hooks/useAuth";
import { useLang } from "../contexts/LanguageContext";
import "../styles/stub.css";

export default function ChatPage() {
  const { t } = useLang();
  const { user, isAuthenticated } = useAuth();
  const canAccess = isAuthenticated && user?.plan === "proai";
  useEffect(() => { document.title = `Quikram — ${t("nav.chat")}`; }, [t]);

  return (
    <div className="stub-page">
      <Navbar />
      <main className="stub-content">
        {canAccess ? (
          <div className="stub-hero">
            <h1 className="stub-title">{t("nav.chat")}</h1>
            <p className="stub-desc">{t("chat.desc")}</p>
          </div>
        ) : (
          <div className="stub-hero">
            <h1 className="stub-title">{t("nav.chat")}</h1>
            <p className="stub-error-msg">{t("chat.restricted")}</p>
            <Link to="/pricing" className="stub-btn">{t("chat.upgrade")}</Link>
          </div>
        )}
      </main>
      <Footer />
    </div>
  );
}
