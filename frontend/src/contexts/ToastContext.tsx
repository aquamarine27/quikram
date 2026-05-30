import { createContext, useCallback, useRef, useState, type ReactNode } from "react";
import "../styles/toast.css";

interface ToastItem {
  id: number;
  message: string;
  type: "success" | "error";
  leaving: boolean;
}

interface ToastContextType {
  showToast: (message: string, type?: "success" | "error") => void;
}

export const ToastContext = createContext<ToastContextType | null>(null);

let nextId = 0;

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastItem[]>([]);
  const timers = useRef<Map<number, ReturnType<typeof setTimeout>>>(new Map());

  const removeToast = useCallback((id: number) => {
    setToasts((prev) => prev.map((t) => (t.id === id ? { ...t, leaving: true } : t)));
    setTimeout(() => {
      setToasts((prev) => prev.filter((t) => t.id !== id));
    }, 300);
  }, []);

  const showToast = useCallback(
    (message: string, type: "success" | "error" = "success") => {
      const id = nextId++;
      setToasts((prev) => [...prev, { id, message, type, leaving: false }]);
      const timer = setTimeout(() => removeToast(id), 3500);
      timers.current.set(id, timer);
    },
    [removeToast],
  );

  return (
    <ToastContext.Provider value={{ showToast }}>
      {children}
      <div className="toast-container">
        {toasts.map((t) => (
          <div key={t.id} className={`toast toast-${t.type} ${t.leaving ? "toast-out" : ""}`}>
            <span className="toast-icon">
              {t.type === "success" ? SuccessIcon : ErrorIcon}
            </span>
            <span className="toast-message">{t.message}</span>
            <div className="toast-bar" />
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
}

const SuccessIcon = (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
    <circle cx="12" cy="12" r="10" />
    <polyline points="16 8 10 16 7 13" />
  </svg>
);

const ErrorIcon = (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
    <circle cx="12" cy="12" r="10" />
    <line x1="15" y1="9" x2="9" y2="15" />
    <line x1="9" y1="9" x2="15" y2="15" />
  </svg>
);
