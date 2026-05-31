import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import { useLang } from "../contexts/LanguageContext";
import "../styles/pricing.css";

export default function PricingPage() {
  const { t } = useLang();
  return (
    <div className="pricing-page">
      <Navbar />
      <main className="pricing-content">
        <div className="pricing-hero">
          <h1 className="pricing-title">{t("nav.pricing")}</h1>
          <p className="pricing-desc">{t("pricing.desc")}</p>
        </div>
      </main>
      <Footer />
    </div>
  );
}
