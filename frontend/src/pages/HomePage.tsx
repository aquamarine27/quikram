import { useEffect } from "react";
import { useAuth } from "../hooks/useAuth";
import { useLang } from "../contexts/LanguageContext";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import "../styles/home.css";

export default function HomePage() {
  const { user } = useAuth();
  const { t } = useLang();
  useEffect(() => { document.title = `Quikram — ${t("nav.home")}`; }, [t]);

  return (
    <div className="home-page">
      <Navbar />
      <main className="home-content">
        <div className="home-hero">
          <h1 className="home-title">
            {user ? t("home.welcome", { name: user.name }) : t("home.title")}
          </h1>
          <p className="home-desc">{t("home.desc")}</p>
        </div>
      </main>
      <Footer />
    </div>
  );
}
