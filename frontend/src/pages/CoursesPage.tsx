import { useEffect } from "react";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import { useLang } from "../contexts/LanguageContext";
import "../styles/stub.css";

export default function CoursesPage() {
  const { t } = useLang();
  useEffect(() => { document.title = `Quikram — ${t("nav.courses")}`; }, [t]);
  return (
    <div className="stub-page">
      <Navbar />
      <main className="stub-content">
        <div className="stub-hero">
          <h1 className="stub-title">{t("nav.courses")}</h1>
          <p className="stub-desc">{t("courses.desc")}</p>
        </div>
      </main>
      <Footer />
    </div>
  );
}
