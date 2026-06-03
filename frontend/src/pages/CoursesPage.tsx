import { useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import { useLang } from "../contexts/LanguageContext";
import "../styles/courses.css";

export default function CoursesPage() {
  const { t } = useLang();
  const navigate = useNavigate();
  const existingRef = useRef<HTMLDivElement>(null);
  const mineRef = useRef<HTMLDivElement>(null);

  useEffect(() => { document.title = `Quikram — ${t("nav.courses")}`; }, [t]);

  const onMouseMove = (e: React.MouseEvent, ref: React.RefObject<HTMLDivElement | null>) => {
    if (!ref.current) return;
    const rect = ref.current.getBoundingClientRect();
    const x = ((e.clientX - rect.left) / rect.width) * 100;
    const y = ((e.clientY - rect.top) / rect.height) * 100;
    ref.current.style.setProperty("--mx", `${x}%`);
    ref.current.style.setProperty("--my", `${y}%`);
  };

  return (
    <div className="courses-page">
      <Navbar />
      <div className="courses-content">
        <div className="courses-header">
          <h1>{t("nav.courses")}</h1>
        </div>

        <div className="courses-grid">
          <div
            ref={existingRef}
            className="courses-card courses-card-existing"
            onMouseMove={(e) => onMouseMove(e, existingRef)}
            style={{ cursor: "default" }}
          >
            <div className="courses-card-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                <circle cx="12" cy="12" r="10"/>
                <line x1="2" y1="12" x2="22" y2="12"/>
                <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
              </svg>
            </div>
            <h2>{t("courses.card.existing")}</h2>
            <p>{t("courses.card.existingDesc")}</p>
            <span className="courses-card-badge">{t("pricing.comingSoon")}</span>
          </div>

          <div
            ref={mineRef}
            className="courses-card courses-card-mine"
            onMouseMove={(e) => onMouseMove(e, mineRef)}
            onClick={() => navigate("/courses/mine")}
          >
            <div className="courses-card-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14 2 14 8 20 8"/>
                <line x1="12" y1="18" x2="12" y2="12"/>
                <line x1="9" y1="15" x2="15" y2="15"/>
              </svg>
            </div>
            <h2>{t("courses.card.mine")}</h2>
            <p>{t("courses.card.mineDesc")}</p>
            <div className="courses-card-arrow">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                <line x1="5" y1="12" x2="19" y2="12"/>
                <polyline points="12 5 19 12 12 19"/>
              </svg>
            </div>
          </div>
        </div>
      </div>
      <Footer />
    </div>
  );
}
