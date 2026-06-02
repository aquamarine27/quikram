import { useEffect, useRef, useState } from "react";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import { useLang } from "../contexts/LanguageContext";
import "../styles/about.css";

function AnimatedNumber({ target, suffix = "" }: { target: number; suffix?: string }) {
  const [count, setCount] = useState(0);
  const ref = useRef<HTMLSpanElement>(null);
  const [started, setStarted] = useState(false);

  useEffect(() => {
    const el = ref.current;
    if (!el) return;
    const observer = new IntersectionObserver(
      ([entry]) => { if (entry.isIntersecting) setStarted(true); },
      { threshold: 0.3 },
    );
    observer.observe(el);
    return () => observer.disconnect();
  }, []);

  useEffect(() => {
    if (!started) return;
    let raf: number;
    const startTime = performance.now();
    const duration = 3000;
    const animate = (now: number) => {
      const elapsed = now - startTime;
      const progress = Math.min(elapsed / duration, 1);
      const eased = 1 - Math.pow(1 - progress, 4);
      setCount(Math.floor(eased * target));
      if (progress < 1) raf = requestAnimationFrame(animate);
    };
    raf = requestAnimationFrame(animate);
    return () => cancelAnimationFrame(raf);
  }, [started, target]);

  return <span ref={ref}>{count.toLocaleString()}{suffix}</span>;
}

const stats = [
  { key: "about.stats.users", target: 2500, suffix: "+" },
  { key: "about.stats.subjects", target: 1200, suffix: "+" },
  { key: "about.stats.summaries", target: 8000, suffix: "+" },
  { key: "about.stats.tests", target: 15000, suffix: "+" },
];

export default function AboutPage() {
  const { t } = useLang();
  useEffect(() => { document.title = `Quikram — ${t("nav.about")}`; }, [t]);

  return (
    <div className="about-page">
      <Navbar />
      <main className="about-content">
        <section className="about-hero">
          <span className="about-hero-label">Quikram</span>
          <h1 className="about-hero-title">{t("nav.about")}</h1>
          <p className="about-hero-desc">{t("about.desc")}</p>
        </section>

        <section className="about-stats">
          <div className="about-stats-grid">
            {stats.map((s) => (
              <div key={s.key} className="about-stat-card">
                <span className="about-stat-num">
                  <AnimatedNumber target={s.target} suffix={s.suffix} />
                </span>
                <span className="about-stat-label">{t(s.key)}</span>
              </div>
            ))}
          </div>
        </section>

        <section className="about-contact">
          <h2 className="about-contact-title">{t("about.contact")}</h2>
          <div className="about-contact-items">
            <a href="mailto:quikram@support.com" className="about-contact-item">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <rect x="2" y="4" width="20" height="16" rx="2" />
                <path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7" />
              </svg>
              <span>quikram@support.com</span>
            </a>
            <div className="about-contact-item">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z" />
                <circle cx="12" cy="10" r="3" />
              </svg>
              <span>г. Москва, ул. Тверская, д. 1</span>
            </div>
            <a href="tel:+70000000000" className="about-contact-item">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z" />
              </svg>
              <span>+7 (000) 000-00-00</span>
            </a>
          </div>
        </section>
      </main>
      <Footer />
    </div>
  );
}
