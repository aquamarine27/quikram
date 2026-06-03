import client from "./client";

export interface Subject {
  id: string;
  user_id: string;
  title: string;
  category: string;
  created_at: string;
  updated_at: string;
}

export async function getSubjects(): Promise<Subject[]> {
  const res = await client.get("/subjects");
  return res.data;
}

export async function createSubject(title: string, category = ""): Promise<Subject> {
  const res = await client.post("/subjects", { title, category });
  return res.data;
}
