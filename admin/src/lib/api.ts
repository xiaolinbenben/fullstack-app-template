export type Envelope<T> = {
  success: boolean;
  message: string;
  data: T;
};

export async function fetchHealth(): Promise<Envelope<{ status: string }>> {
  const response = await fetch("/healthz");
  const payload = (await response.json()) as Envelope<{ status: string }>;
  if (!response.ok || !payload.success) {
    throw new Error(payload.message || "服务不可用");
  }
  return payload;
}
