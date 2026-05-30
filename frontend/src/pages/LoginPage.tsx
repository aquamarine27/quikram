import { useEffect, useMemo, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useAuth } from "../hooks/useAuth";
import { useToast } from "../hooks/useToast";
import { useLang } from "../contexts/LanguageContext";
import ThemeToggle from "../components/ThemeToggle";
import LangSwitcher from "../components/LangSwitcher";
import "../styles/login.css";

// ─── Form type ───
type FormData = {
  email: string;
  password: string;
};

// ─── Component ───
export default function LoginPage() {
  const { login, isAuthenticated, isLoading } = useAuth();
  const { showToast } = useToast();
  const { t } = useLang();
  const navigate = useNavigate();

  const [expanded, setExpanded] = useState(false);
  const [formError, setFormError] = useState(false);
  const [showPwd, setShowPwd] = useState(false);

  // Schema with live translation
  const schema = useMemo(() => z.object({
    email: z.string().email(t("val.invalidEmail")),
    password: z.string().min(6, t("val.minChars")),
  }), [t]);

  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<FormData>({
    resolver: zodResolver(schema),
  });

  // ─── Effects ───
  useEffect(() => { document.title = `Quikram — ${t("auth.login")}`; }, [t]);
  useEffect(() => {
    if (!isLoading && isAuthenticated) navigate("/profile", { replace: true });
  }, [isAuthenticated, isLoading, navigate]);

  // ─── Handlers ───
  const onFocus = () => {
    if (!expanded) setExpanded(true);
    if (formError) setFormError(false);
  };

  const onSubmit = async (data: FormData) => {
    try {
      setFormError(false);
      await login(data);
      showToast(t("toast.welcomeBack"), "success");
      navigate("/profile", { replace: true });
    } catch (err: unknown) {
      setFormError(true);
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error
        || (err as Error).message
        || t("toast.error");
      showToast(msg, "error");
    }
  };

  const onError = () => setFormError(true);

  const cardClass = [
    "login-card",
    expanded ? "login-card-expanded" : "",
    formError ? "login-card-error" : "",
  ].filter(Boolean).join(" ");

  // ─── Render ───
  return (
    <div className="login-page">
      <div className={cardClass}>
        <div className="login-controls">
          <ThemeToggle />
          <LangSwitcher />
        </div>

        <h1 className="login-title">{t("auth.login")}</h1>
        <p className="login-subtitle">{t("auth.welcome")}</p>

        <form className="login-form" onSubmit={handleSubmit(onSubmit, onError)}>
          <div className="field">
            <div className={`field-inner${errors.email ? " field-inner-error" : ""}`}>
              <input className="field-input" placeholder=" " {...register("email")} onFocus={onFocus} />
              <label className="field-label">{t("auth.email")}</label>
            </div>
            {errors.email && <p className="field-error">{errors.email.message}</p>}
          </div>

          <div className="field">
            <div className={`field-inner field-password${errors.password ? " field-inner-error" : ""}`}>
              <input
                className="field-input"
                type={showPwd ? "text" : "password"}
                placeholder=" "
                {...register("password")}
                onFocus={onFocus}
              />
              <label className="field-label">{t("auth.password")}</label>
              <button type="button" className="field-eye" onClick={() => setShowPwd((v) => !v)} aria-label="Toggle password visibility">
                {showPwd ? EyeOffIcon : EyeIcon}
              </button>
            </div>
            {errors.password && <p className="field-error">{errors.password.message}</p>}
          </div>

          <button className="login-btn" type="submit" disabled={isSubmitting}>
            {isSubmitting ? t("auth.loading") : t("auth.signin")}
          </button>
        </form>

        <p className="login-footer">
          {t("auth.noAccount")} <Link to="/register">{t("auth.register")}</Link>
        </p>

        <div className="login-divider">
          <span className="login-divider-text">{t("auth.loginVia")}</span>
        </div>

        <div className="login-social">
          <button className="login-social-btn" aria-label="Yandex" title="Yandex">{YandexIcon}</button>
          <button className="login-social-btn" aria-label="Google" title="Google">{GoogleIcon}</button>
          <button className="login-social-btn" aria-label="Apple" title="Apple">{AppleIcon}</button>
        </div>
      </div>
    </div>
  );
}

// ─── Icons ───
const EyeIcon = (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
    <circle cx="12" cy="12" r="3" />
  </svg>
);

const EyeOffIcon = (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94" />
    <path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19" />
    <line x1="1" y1="1" x2="23" y2="23" />
  </svg>
);

const YandexIcon = (
  <svg width="24" height="24" viewBox="0 0 24 24">
    <rect width="24" height="24" rx="6" fill="#FC3F1D" />
    <text x="12" y="17" textAnchor="middle" fill="#fff" fontSize="14" fontWeight="700" fontFamily="Inter, sans-serif">Я</text>
  </svg>
);

const GoogleIcon = (
  <svg width="24" height="24" viewBox="0 0 24 24">
    <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" />
    <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
    <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" />
    <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
  </svg>
);

const AppleIcon = (
  <svg width="24" height="24" viewBox="0 0 24 24">
    <path fill="currentColor" d="M17.05 20.28c-.98.95-2.05.88-3.08.4-1.09-.5-2.08-.48-3.15 0-1.35.66-2.13.48-3.1-.4C4.62 17.28 4.1 12.93 5.98 9.62c1.05-1.84 2.63-2.78 4.3-2.78 1.18 0 2.16.46 2.95.46.78 0 2.01-.55 3.36-.47 1.22.04 2.35.58 3.2 1.53-2.67 1.68-2.24 4.7.46 5.74-.45 1.2-1.03 2.4-1.76 3.46-.57.77-1.1 1.54-1.44 2.28zM12.03 7.25c-.26-2.14 1.77-3.97 3.73-4.25.16 2.03-1.78 3.98-3.73 4.25z" />
  </svg>
);
