import client from "./client";

export interface Review {
  text: string;
  name: string;
  author: string;
  badge: string;
}

export async function getReviews(lang: string): Promise<Review[]> {
  const res = await client.get("/reviews", { params: { lang } });
  return res.data;
}
