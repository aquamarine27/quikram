import { useEffect, useState, useRef } from "react";
import { useLang } from "../contexts/LanguageContext";
import { useParallax } from "../hooks/useParallax";
import { useSmoothSnap } from "../hooks/useSmoothSnap";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import { getReviews, type Review } from "../api/reviews";
import "../styles/home.css";

const SLIDE_MS = 3000;
const PAUSE_MS = 6000;
const CYCLE_MS = SLIDE_MS + PAUSE_MS;

function layoutForWidth(w: number) {
  if (w <= 480) return { gap: 10, visible: 2 };
  if (w <= 768) return { gap: 14, visible: 3 };
  return { gap: 20, visible: 3 };
}

export default function HomePage() {
  const { t, lang } = useLang();
  useEffect(() => { document.title = `Quikram — ${t("nav.home")}`; }, [t]);

  const snapRef = useSmoothSnap();
  const heroParallax = useParallax(0.08, 15);
  const photo1Parallax = useParallax(0.06, 12);
  const photo2Parallax = useParallax(0.09, 12);
  const whyParallax = useParallax(0.06, 15);
  const stackBackParallax = useParallax(0.05, 10);
  const stackFrontParallax = useParallax(0.08, 10);
  const reviewsParallax = useParallax(0.06, 15);

  const [reviews, setReviews] = useState<Review[]>([]);
  const [current, setCurrent] = useState(0);
  const [direction, setDirection] = useState(1);
  const [gap, setGap] = useState(20);
  const [visibleCount, setVisibleCount] = useState(3);
  const progressFillRef = useRef<HTMLDivElement>(null);
  const phaseStartRef = useRef(Date.now());
  const dirRef = useRef(direction);
  dirRef.current = direction;
  const prevVisibleRef = useRef(visibleCount);
  prevVisibleRef.current = visibleCount;
  const progressMsRef = useRef(0);
  const resizeRef = useRef<number>();

  const reviewCount = reviews.length;
  const batchCount = Math.ceil(reviewCount / visibleCount) || 1;
  progressMsRef.current = batchCount * CYCLE_MS;

  useEffect(() => {
    getReviews(lang).then((data) => setReviews(data));
  }, [lang]);

  useEffect(() => {
    const onResize = () => {
      const layout = layoutForWidth(window.innerWidth);
      setGap(layout.gap);
      setVisibleCount(layout.visible);
      if (layout.visible !== prevVisibleRef.current) {
        setCurrent(0);
        setDirection(1);
        phaseStartRef.current = Date.now();
      }
    };
    onResize();
    clearTimeout(resizeRef.current);
    resizeRef.current = window.setTimeout(onResize, 100);
    window.addEventListener("resize", onResize);
    return () => {
      window.removeEventListener("resize", onResize);
      clearTimeout(resizeRef.current);
    };
  }, []);

  useEffect(() => {
    if (reviewCount === 0) return;
    const timer = setTimeout(() => {
      let next = current + direction * visibleCount;
      let newDir = direction;

      if (next < 0) {
        newDir = 1;
        next = visibleCount;
        phaseStartRef.current = Date.now();
      } else if (next + visibleCount > reviewCount) {
        newDir = -1;
        next = reviewCount - visibleCount - visibleCount;
        phaseStartRef.current = Date.now();
      }

      if (newDir !== direction) setDirection(newDir);
      setCurrent(next);
    }, CYCLE_MS);
    return () => clearTimeout(timer);
  }, [current, direction, visibleCount, reviewCount]);

  useEffect(() => {
    let raf: number;
    const tick = () => {
      const dir = dirRef.current;
      const phaseElapsed = Date.now() - phaseStartRef.current;
      const total = progressMsRef.current;
      const pct = total > 0 ? Math.min(phaseElapsed / total, 1) * 100 : 0;
      const display = dir === 1 ? pct : 100 - pct;

      if (progressFillRef.current) {
        progressFillRef.current.style.width = `${display}%`;
      }

      raf = requestAnimationFrame(tick);
    };
    raf = requestAnimationFrame(tick);
    return () => cancelAnimationFrame(raf);
  }, []);

  return (
    <div className="home-page" ref={snapRef}>
      <Navbar />

      <section className="home-section home-section--hero">
        <div className="home-section-bg" ref={heroParallax.ref} style={heroParallax.style} />
        <div className="home-section-inner">
          <div className="home-hero-left">
            <span className="home-hero-badge">{t("home.welcomeBadge")}</span>
            <h1 className="home-hero-title">{t("home.heroTitle")}</h1>
            <p className="home-hero-sub">{t("home.heroSubtitle")}</p>
            <div className="home-search">
              <input className="home-search-input" type="text" placeholder={t("home.searchPlaceholder")} />
              <button className="home-search-btn" aria-label="Search">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                  <line x1="5" y1="12" x2="19" y2="12" />
                  <polyline points="12 5 19 12 12 19" />
                </svg>
              </button>
            </div>
          </div>

          <div className="home-hero-right">
            <div className="home-collage">
              <div className="home-collage-photo home-collage-photo--1" ref={photo1Parallax.ref} style={photo1Parallax.style}>
                <div className="home-collage-inner-photo home-collage-inner-photo--1" />
              </div>
              <div className="home-collage-photo home-collage-photo--2" ref={photo2Parallax.ref} style={photo2Parallax.style}>
                <div className="home-collage-inner-photo home-collage-inner-photo--2" />
              </div>
              <div className="home-collage-badge home-collage-badge--members">
                <span className="home-collage-badge-num">{t("home.membersBadge")}</span>
                <div className="home-collage-avatars">
                  <span className="home-collage-avatar" style={{ background: "#ef4444" }} />
                  <span className="home-collage-avatar" style={{ background: "#f59e0b" }} />
                  <span className="home-collage-avatar" style={{ background: "#22c55e" }} />
                </div>
              </div>
              <div className="home-collage-badge home-collage-badge--reviews">
                <span className="home-collage-badge-num">{t("home.reviewsBadge")}</span>
              </div>
              <span className="home-dot home-dot--green" />
              <span className="home-dot home-dot--red" />
              <span className="home-dot home-dot--blue" />
            </div>
          </div>
        </div>
      </section>

      <section className="home-section home-section--why">
        <div className="home-section-bg" ref={whyParallax.ref} style={whyParallax.style} />
        <div className="home-section-inner">
          <div className="home-why-left">
            <span className="home-why-badge">{t("home.whyBadge")}</span>
            <h2 className="home-why-title">{t("home.whyTitle")}</h2>
            <p className="home-why-desc">{t("home.whyDesc")}</p>
            <div className="home-features">
              <div className="home-feature home-feature--top-left">
                <span className="home-feature-check">✓</span>
                <div>
                  <strong className="home-feature-title">{t("home.feature1.title")}</strong>
                  <span className="home-feature-desc">{t("home.feature1.desc")}</span>
                </div>
              </div>
              <div className="home-feature home-feature--top-right">
                <span className="home-feature-check">✓</span>
                <div>
                  <strong className="home-feature-title">{t("home.feature2.title")}</strong>
                  <span className="home-feature-desc">{t("home.feature2.desc")}</span>
                </div>
              </div>
              <div className="home-feature home-feature--bottom">
                <span className="home-feature-check">✓</span>
                <div>
                  <strong className="home-feature-title">{t("home.feature3.title")}</strong>
                  <span className="home-feature-desc">{t("home.feature3.desc")}</span>
                </div>
              </div>
            </div>
          </div>

          <div className="home-why-right">
            <div className="home-stack">
              <div className="home-stack-photo home-stack-photo--back" ref={stackBackParallax.ref} style={stackBackParallax.style}>
                <div className="home-stack-inner home-stack-inner--back" />
              </div>
              <div className="home-stack-photo home-stack-photo--front" ref={stackFrontParallax.ref} style={stackFrontParallax.style}>
                <div className="home-stack-inner home-stack-inner--front" />
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="home-section home-section--reviews">
        <div className="home-section-bg" ref={reviewsParallax.ref} style={reviewsParallax.style}>
          <div className="home-orb home-orb--1" />
          <div className="home-orb home-orb--2" />
          <div className="home-orb home-orb--3" />
          <div className="home-orb home-orb--4" />
        </div>
        <div className="home-section-inner home-section-inner--reviews">
          <h2 className="home-reviews-title">{t("home.reviewsTitle")}</h2>
          <span className="home-reviews-badge">{t("home.reviewsBadgeInner")}</span>

          <div className="home-reviews-track">
            <div
              className="home-reviews-slider"
              style={{
                gap: `${gap}px`,
                "--card-gap": `${gap}px`,
                "--visible-count": visibleCount,
                transform: `translateX(calc(-${current / visibleCount} * 100% - ${(current / visibleCount) * gap}px))`,
                transition: `transform ${SLIDE_MS}ms cubic-bezier(0.88, 0, 0.12, 1)`,
              } as React.CSSProperties}
            >
              {reviews.map((r, i) => (
                <div key={i} className="home-review-card">
                  <div className="home-review-top">
                    <div className="home-review-avatar" style={{ background: r.badge }}>
                      {r.name[0]}
                    </div>
                    <div className="home-review-info">
                      <span className="home-review-name">{r.name}</span>
                      <span className="home-review-meta">{r.author}</span>
                    </div>
                  </div>
                  <p className="home-review-text">"{r.text}"</p>
                  <span className="home-review-badge" style={{ background: r.badge }} />
                </div>
              ))}
            </div>
          </div>

          <div className="home-reviews-progress">
            <div className="home-reviews-progress-fill" ref={progressFillRef} />
          </div>
        </div>
      </section>

      <Footer />
    </div>
  );
}
