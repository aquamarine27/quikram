import { useState, useRef, useEffect } from "react";
import { useLang } from "../contexts/LanguageContext";
import type { Lang } from "../i18n/translations";
import "../styles/lang-switcher.css";

// ─── Options ───
const langs: { code: Lang; label: string }[] = [
  { code: "ru", label: "RU" },
  { code: "en", label: "EN" },
];

// ─── Component ───
export default function LangSwitcher() {
  const { lang, setLang } = useLang();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  // Close on outside click
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", handler);
    return () => document.removeEventListener("mousedown", handler);
  }, []);

  return (
    <div className="lang-switcher" ref={ref}>
      <button className="lang-btn" onClick={() => setOpen((v) => !v)}>
        <span>{langs.find((l) => l.code === lang)?.label}</span>
        <Chevron open={open} />
      </button>

      <div className={`lang-dropdown ${open ? "lang-dropdown-open" : ""}`}>
        {langs.map((l) => (
          <button
            key={l.code}
            className={`lang-opt ${l.code === lang ? "lang-opt-active" : ""}`}
            onClick={() => { setLang(l.code); setOpen(false); }}
          >
            {l.label}
          </button>
        ))}
      </div>
    </div>
  );
}

// ─── Chevron icon ───
function Chevron({ open }: { open: boolean }) {
  return (
    <svg
      className={`lang-chevron ${open ? "lang-chevron-open" : ""}`}
      width="10" height="10" viewBox="0 0 24 24"
      fill="none" stroke="currentColor" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round"
    >
      <polyline points="6 9 12 15 18 9" />
    </svg>
  );
}
