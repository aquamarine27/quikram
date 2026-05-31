import { useEffect } from "react";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import { useLang } from "../contexts/LanguageContext";
import "../styles/about.css";

export default function AboutPage() {
  const { t } = useLang();
  useEffect(() => { document.title = `Quikram — ${t("nav.about")}`; }, [t]);
  return (
    <div className="about-page">
      <Navbar />
      <main className="about-content">
        <div className="about-hero">
          <h1 className="about-title">{t("nav.about")}</h1>
          <p className="about-desc">{t("about.desc")}</p>
        </div>
      </main>
      <Footer />
    </div>
  );
}
