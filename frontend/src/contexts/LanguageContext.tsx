import { createContext, useCallback, useContext, useEffect, useState, type ReactNode } from "react";
import { translations, type Lang } from "../i18n/translations";

const LANG_KEY = "quikram_lang";

interface LangContextType {
  lang: Lang;
  setLang: (lang: Lang) => void;
  t: (key: string, params?: Record<string, string>) => string;
}

export const LangContext = createContext<LangContextType | null>(null);

export function LangProvider({ children }: { children: ReactNode }) {
  const [lang, setLangState] = useState<Lang>(() => {
    if (typeof window !== "undefined") {
      const stored = localStorage.getItem(LANG_KEY);
      if (stored === "ru" || stored === "en") return stored;
    }
    return "ru";
  });

  useEffect(() => {
    document.documentElement.setAttribute("lang", lang);
    localStorage.setItem(LANG_KEY, lang);
  }, [lang]);

  const setLang = useCallback((l: Lang) => setLangState(l), []);

  const t = useCallback(
    (key: string, params?: Record<string, string>): string => {
      let val = translations[lang]?.[key] ?? key;
      if (params) {
        for (const [k, v] of Object.entries(params)) {
          val = val.replace(`{${k}}`, v);
        }
      }
      return val;
    },
    [lang],
  );

  return (
    <LangContext.Provider value={{ lang, setLang, t }}>
      {children}
    </LangContext.Provider>
  );
}

export function useLang() {
  const ctx = useContext(LangContext);
  if (!ctx) throw new Error("useLang must be used within LangProvider");
  return ctx;
}
