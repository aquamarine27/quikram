import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AuthProvider } from "./contexts/AuthContext";
import { ThemeProvider } from "./contexts/ThemeContext";
import { LangProvider } from "./contexts/LanguageContext";
import { useAuth } from "./hooks/useAuth";
import { ToastProvider } from "./contexts/ToastContext";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";
import ProfilePage from "./pages/ProfilePage";
import "./styles/themes.css";

// ─── Root redirect ───
function RootRedirect() {
  const { isAuthenticated, isLoading } = useAuth();
  if (isLoading) return <div style={{ minHeight: "100vh", background: "var(--page-bg)" }} />;
  return isAuthenticated ? <Navigate to="/profile" replace /> : <Navigate to="/register" replace />;
}

// ─── App ───
export default function App() {
  return (
    <BrowserRouter>
      <ThemeProvider>
        <LangProvider>
          <ToastProvider>
            <AuthProvider>
              <Routes>
                <Route path="/" element={<RootRedirect />} />
                <Route path="/login" element={<LoginPage />} />
                <Route path="/register" element={<RegisterPage />} />
                <Route path="/profile" element={<ProfilePage />} />
                <Route path="*" element={<Navigate to="/" replace />} />
              </Routes>
            </AuthProvider>
          </ToastProvider>
        </LangProvider>
      </ThemeProvider>
    </BrowserRouter>
  );
}
