import { useEffect, useState } from "react";
import { useAuth } from "../hooks/useAuth";
import { useLang } from "../contexts/LanguageContext";
import { useToast } from "../hooks/useToast";
import { getPlans, type Plan } from "../api/plans";
import { changePlan } from "../api/auth";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import "../styles/pricing.css";

const featureRows: { key: string; featureKey: keyof Plan["features"] }[] = [
  { key: "pricing.subjects", featureKey: "subjects" },
  { key: "pricing.uploads", featureKey: "uploads_per_month" },
  { key: "pricing.aiSummary", featureKey: "ai_summary" },
  { key: "pricing.compression", featureKey: "compression" },
  { key: "pricing.basicTests", featureKey: "basic_tests" },
  { key: "pricing.advancedTests", featureKey: "advanced_tests" },
  { key: "pricing.difficulty", featureKey: "difficulty" },
  { key: "pricing.analytics", featureKey: "analytics" },
  { key: "pricing.weakSpots", featureKey: "weak_spots" },
  { key: "pricing.export", featureKey: "export" },
  { key: "pricing.aiChat", featureKey: "ai_chat" },
];

const cardFeatureKeys: { key: string; featureKey: keyof Plan["features"] }[] = [
  { key: "pricing.subjects", featureKey: "subjects" },
  { key: "pricing.uploads", featureKey: "uploads_per_month" },
  { key: "pricing.aiSummary", featureKey: "ai_summary" },
  { key: "pricing.basicTests", featureKey: "basic_tests" },
  { key: "pricing.advancedTests", featureKey: "advanced_tests" },
  { key: "pricing.analytics", featureKey: "analytics" },
  { key: "pricing.export", featureKey: "export" },
  { key: "pricing.aiChat", featureKey: "ai_chat" },
];

function fmtFeature(key: string, plan: Plan, t: (k: string, p?: any) => string): string {
  const val = plan.features[featureRows.find((r) => r.key === key)!.featureKey];
  if (typeof val === "number") {
    if (val === -1) return t("pricing.unlimited");
    if (key === "pricing.uploads" && val === 15) return "15";
    return String(val);
  }
  return val ? t("pricing.yes") : t("pricing.no");
}

function cardFeatureIncluded(featureKey: keyof Plan["features"], plan: Plan): boolean {
  const val = plan.features[featureKey];
  return typeof val === "number" ? val !== 0 : val;
}

function useButton(
  plan: Plan, isAuth: boolean, currentPlan: string | null, t: (k: string, p?: any) => string,
): { text: string; disabled: boolean } {
  if (!isAuth) {
    if (plan.id === "free") return { text: t("pricing.start"), disabled: false };
    if (plan.id === "pro") return { text: t("pricing.tryPro"), disabled: false };
    return { text: t("pricing.tryProAI"), disabled: false };
  }
  if (currentPlan === plan.id) return { text: t("pricing.current"), disabled: true };
  return { text: t("pricing.upgrade", { name: t(`pricing.${plan.id}`) }), disabled: false };
}

export default function PricingPage() {
  const { t } = useLang();
  const { user, isAuthenticated, updateUser } = useAuth();
  const { showToast } = useToast();
  const [plans, setPlans] = useState<Plan[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    document.title = `Quikram — ${t("nav.pricing")}`;
  }, [t]);

  useEffect(() => {
    getPlans()
      .then(setPlans)
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const handleClick = async (planId: string) => {
    if (!isAuthenticated) {
      showToast(t("pricing.comingSoon"), "success");
      return;
    }
    try {
      const updated = await changePlan(planId);
      updateUser(updated);
      showToast(t("pricing.planChanged"), "success");
    } catch (err) {
      showToast(err instanceof Error ? err.message : t("toast.error"), "error");
    }
  };

  return (
    <div className="pricing-page">
      <Navbar />
      <main className="pricing-content">
        <div className="pricing-hero">
          <h1 className="pricing-hero-title">{t("pricing.desc")}</h1>
          <p className="pricing-hero-sub">{t("pricing.subtitle")}</p>
        </div>

        {loading ? (
          <div className="pricing-loading">{t("auth.loading")}</div>
        ) : (
          <>
            <div className="pricing-cards">
              {plans.map((plan) => {
                const btn = useButton(plan, isAuthenticated, user?.plan ?? null, t);
                const isCurrent = isAuthenticated && user?.plan === plan.id;
                const cardClasses = [
                  "pricing-card",
                  plan.highlighted ? "pricing-card--highlighted" : "",
                  isCurrent ? "pricing-card--current" : "",
                ].filter(Boolean).join(" ");
                const btnClasses = [
                  "pricing-btn",
                  btn.disabled ? "pricing-btn--disabled" : "",
                  isCurrent && plan.id === "free" ? "pricing-btn--current-free" : "",
                ].filter(Boolean).join(" ");

                const badgeKey: Record<string, string> = {
                  popular: "pricing.proBadge",
                  ai_chat: "pricing.proaiBadge",
                };

                return (
                  <div key={plan.id} className={cardClasses}>
                    {plan.badge && <span className="pricing-badge">{t(badgeKey[plan.badge] || plan.badge)}</span>}

                    <div className="pricing-card-header">
                      <h2 className="pricing-card-name">{t(`pricing.${plan.id}`)}</h2>
                      <div className="pricing-card-price">
                        {plan.price === 0 ? t("pricing.freePrice") : `${plan.price / 100} ₽`}
                      </div>
                      <div className="pricing-card-period">
                        {plan.period === "forever" ? t("pricing.freePeriod") : t("pricing.month")}
                      </div>
                    </div>

                    <div className="pricing-card-features">
                      {cardFeatureKeys.map((f) => {
                        const included = cardFeatureIncluded(f.featureKey, plan);
                        return (
                          <div key={f.key} className={`pricing-feature${included ? "" : " pricing-feature--off"}`}>
                            <span className="pricing-feature-icon">{included ? "✓" : "✗"}</span>
                            <span>{t(f.key)}</span>
                          </div>
                        );
                      })}
                    </div>

                    <button
                      className={btnClasses}
                      disabled={btn.disabled}
                      onClick={btn.disabled ? undefined : () => handleClick(plan.id)}
                    >
                      {btn.text}
                    </button>
                  </div>
                );
              })}
            </div>

            <div className="pricing-compare">
              <h2 className="pricing-compare-title">{t("pricing.compTitle")}</h2>
              <div className="pricing-table-wrap">
                <table className="pricing-table">
                  <thead>
                    <tr>
                      <th>{t("pricing.features")}</th>
                      {plans.map((p) => {
                        const isCurrent = isAuthenticated && user?.plan === p.id;
                        return (
                          <th key={p.id} className={isCurrent ? "col-highlighted" : ""}>
                            {t(`pricing.${p.id}`)}
                          </th>
                        );
                      })}
                    </tr>
                  </thead>
                  <tbody>
                    {featureRows.map((row) => (
                      <tr key={row.key}>
                        <td className="pricing-td-label">{t(row.key)}</td>
                        {plans.map((p) => {
                          const isCurrent = isAuthenticated && user?.plan === p.id;
                          return (
                            <td key={p.id} className={`pricing-td-val${isCurrent ? " col-highlighted" : ""}`}>
                              {fmtFeature(row.key, p, t)}
                            </td>
                          );
                        })}
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </>
        )}
      </main>
      <Footer />
    </div>
  );
}
