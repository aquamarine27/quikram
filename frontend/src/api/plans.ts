import client from "./client";

export interface PlanFeature {
  subjects: number;
  uploads_per_month: number;
  ai_summary: boolean;
  compression: boolean;
  basic_tests: boolean;
  advanced_tests: boolean;
  difficulty: boolean;
  analytics: boolean;
  weak_spots: boolean;
  export: boolean;
  ai_chat: boolean;
}

export interface Plan {
  id: string;
  price: number;
  period: string;
  badge: string;
  highlighted: boolean;
  features: PlanFeature;
}

export async function getPlans(): Promise<Plan[]> {
  const res = await client.get<Plan[]>("/plans");
  return res.data;
}
