import client from "./client";

export interface Document {
  id: string;
  subject_id: string;
  user_id: string;
  filename: string;
  file_size: number;
  mime_type: string;
  status: string;
  created_at: string;
}

export async function getDocuments(subjectId: string): Promise<Document[]> {
  const res = await client.get(`/subjects/${subjectId}/documents`);
  return res.data;
}

export async function uploadDocument(subjectId: string, file: File): Promise<Document> {
  const form = new FormData();
  form.append("file", file);
  const res = await client.post(`/subjects/${subjectId}/documents`, form, {
    headers: { "Content-Type": "multipart/form-data" },
  });
  return res.data;
}

export async function deleteDocument(subjectId: string, docId: string): Promise<void> {
  await client.delete(`/subjects/${subjectId}/documents/${docId}`);
}

export async function getDocument(subjectId: string, docId: string): Promise<Document> {
  const res = await client.get(`/subjects/${subjectId}/documents/${docId}`);
  return res.data;
}

export interface Summary {
  id: string;
  document_id: string;
  subject_id: string;
  user_id: string;
  content_short: string;
  content_medium: string;
  content_long: string;
  created_at: string;
}

export async function getSummary(subjectId: string, docId: string): Promise<Summary> {
  const res = await client.get(`/subjects/${subjectId}/documents/${docId}/summary`);
  return res.data;
}

export async function createSummary(subjectId: string, docId: string): Promise<{ status?: string; message?: string } & Partial<Summary>> {
  const res = await client.post(`/subjects/${subjectId}/documents/${docId}/summary`);
  return res.data;
}
