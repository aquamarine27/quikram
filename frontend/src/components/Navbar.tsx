import { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";
import { useLang } from "../contexts/LanguageContext";
import SettingsModal from "./SettingsModal";
import NotificationModal from "./NotificationModal";
import "../styles/navbar.css";

const navItems = [
  { key: "nav.home", path: "/home" },
  { key: "nav.courses", path: "/courses" },
  { key: "nav.tests", path: "/tests" },
  { key: "nav.chat", path: "/chat" },
  { key: "nav.pricing", path: "/pricing" },
  { key: "nav.about", path: "/about" },
];

export default function Navbar() {
  const [open, setOpen] = useState(false);
  const [settingsOpen, setSettingsOpen] = useState(false);
  const [notifOpen, setNotifOpen] = useState(false);
  const { user } = useAuth();
  const { t } = useLang();
  const location = useLocation();

  const activeItem = navItems.find((i) => location.pathname === i.path);

  return (
    <nav className="navbar">
      <div className="navbar-main">
        <button className="navbar-burger" onClick={() => setOpen((v) => !v)} aria-label="Menu">
          <span className={`burger-line ${open ? "burger-line--open" : ""}`} />
          <span className={`burger-line ${open ? "burger-line--open" : ""}`} />
          <span className={`burger-line ${open ? "burger-line--open" : ""}`} />
        </button>

        <button className="navbar-brand" onClick={() => setOpen((v) => !v)}>
          <span className="navbar-logo">Quikram</span>
          <span className={`navbar-arrow ${open ? "navbar-arrow--open" : ""}`}>▼</span>
        </button>

        <span className="navbar-brand-mobile">Quikram</span>

        <div className="navbar-actions">
          <div className="navbar-notif-wrap">
            <button className="navbar-icon-btn" aria-label="Notifications" onClick={() => setNotifOpen((v) => !v)}>
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
                <path d="M13.73 21a2 2 0 0 1-3.46 0" />
              </svg>
            </button>
            {notifOpen && <NotificationModal onClose={() => setNotifOpen(false)} />}
          </div>

          <div className="navbar-settings-wrap">
            <button className="navbar-icon-btn" aria-label="Settings" onClick={() => setSettingsOpen((v) => !v)}>
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <circle cx="12" cy="12" r="3" />
                <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z" />
              </svg>
            </button>
            {settingsOpen && <SettingsModal onClose={() => setSettingsOpen(false)} />}
          </div>

          <Link to="/profile" className="navbar-profile">
            <span className="navbar-avatar">
              {(user?.name ?? "U").charAt(0).toUpperCase()}
            </span>
            <span className="navbar-username">{user?.name ?? "User"}</span>
          </Link>
        </div>
      </div>

      <div className={`navbar-dropdown ${open ? "navbar-dropdown--open" : ""}`}>
        <div className="navbar-nav">
          {navItems.map((item) => (
            <Link
              key={item.path}
              to={item.path}
              className={`navbar-link ${item === activeItem ? "navbar-link--active" : ""}`}
            >
              {t(item.key)}
            </Link>
          ))}
        </div>
      </div>
    </nav>
  );
}
