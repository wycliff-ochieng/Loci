import { Locus, Reply } from "./types";

type AuthPayload = {
  user: any;
  accessToken: string;
  refreshToken?: string;
};

const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:3000";

const authHeaders = () => {
  if (typeof window === "undefined") return {};
  const token = localStorage.getItem("accessToken");
  return token ? { Authorization: `Bearer ${token}` } : {};
};

const persistAuth = (auth: AuthPayload) => {
  if (typeof window === "undefined") return;
  localStorage.setItem("accessToken", auth.accessToken);
  if (auth.refreshToken) localStorage.setItem("refreshToken", auth.refreshToken);
  localStorage.setItem("user", JSON.stringify(auth.user));
};

export const clearAuth = () => {
  if (typeof window === "undefined") return;
  localStorage.removeItem("accessToken");
  localStorage.removeItem("refreshToken");
  localStorage.removeItem("user");
};

export async function login(email: string, password: string, username = "") {
  const res = await fetch(`${API_BASE}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password, username }),
  });

  if (!res.ok) {
    const msg = await res.text();
    throw new Error(msg || "Login failed");
  }

  const data = await res.json();
  persistAuth({
    user: data.User ?? data.user ?? {},
    accessToken: data.AccessToken ?? data.accessToken ?? "",
    refreshToken: data.RefreshToken ?? data.refreshToken ?? "",
  });
  return data;
}

export async function register(payload: { email: string; password: string; username: string; firstname?: string; lastname?: string }) {
  const res = await fetch(`${API_BASE}/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      email: payload.email,
      password: payload.password,
      username: payload.username,
      firstname: payload.firstname ?? payload.username,
      lastname: payload.lastname ?? payload.username,
    }),
  });

  if (!res.ok) {
    const msg = await res.text();
    throw new Error(msg || "Registration failed");
  }

  const data = await res.json();
  return data;
}

const toLocus = (row: any): Locus => ({
  id: row.id ?? row.ID ?? row.lociid ?? row.LociID ?? "",
  userId: row.user_id ?? row.UserID ?? "",
  message: row.message ?? row.Message ?? "",
  createdAt: (row.created_at ?? row.CreatedAt ?? new Date()).toString(),
  coords: {
    lat: row.lat ?? row.Lat ?? row.location?.lat ?? row.Location?.Lat ?? row.location?.Lat ?? 0,
    lng: row.long ?? row.Long ?? row.location?.long ?? row.Location?.Long ?? row.location?.Long ?? 0,
  },
  viewCount: row.view_count ?? row.ViewCount ?? 0,
  repliesCount: row.replies_count ?? row.RepliesCount ?? 0,
});

export async function fetchLoci(): Promise<Locus[]> {
  // Broad box to get everything; adjust once client supplies bounds
  const url = `${API_BASE}/api/get/loci?SouthWestLat=-90&SouthWestLong=-180&NorthEastLat=90&NorthEastLong=180`;
  const res = await fetch(url, { headers: { ...authHeaders() } });
  if (!res.ok) {
    const msg = await res.text();
    throw new Error(msg || `Failed to fetch messages from ${API_BASE}`);
  }
  const data = await res.json();
  return Array.isArray(data) ? data.map(toLocus) : [];
}

export async function fetchLocus(id: string): Promise<{ locus: Locus; replies: Reply[] }> {
  // No dedicated endpoint; fetch all then pick one
  const loci = await fetchLoci();
  const locus = loci.find((l) => String(l.id) === id) ?? loci[0];
  if (!locus) {
    throw new Error("Message not found or API unavailable.");
  }
  const replies = await fetchReplies(id);
  return { locus, replies };
}

export async function fetchReplies(locusId: string): Promise<Reply[]> {
  const res = await fetch(`${API_BASE}/api/loci/${locusId}/replies`, {
    headers: { ...authHeaders() },
  });
  if (!res.ok) return [];
  const data = await res.json();
  return Array.isArray(data)
    ? data.map((r: any) => ({
        id: r.replyid ?? r.ReplyID ?? r.id ?? r.ID ?? "",
        locusId: r.locusid ?? r.LocusID ?? r.locus_id ?? "",
        userId: r.userid ?? r.UserID ?? r.user_id ?? "",
        content: r.content ?? "",
        createdAt: (r.createdat ?? r.CreatedAT ?? r.created_at ?? new Date()).toString(),
      }))
    : [];
}

export async function createReply(locusId: string, content: string) {
  const res = await fetch(`${API_BASE}/api/loci/${locusId}/reply`, {
    method: "POST",
    headers: { "Content-Type": "application/json", ...authHeaders() },
    body: JSON.stringify({ content }),
  });
  if (!res.ok) throw new Error("Failed to create reply");
  const r = await res.json();
  return {
    id: r.replyid ?? r.ReplyID ?? r.id ?? r.ID ?? "",
    locusId: r.locusid ?? r.LocusID ?? r.locus_id ?? "",
    userId: r.userid ?? r.UserID ?? r.user_id ?? "",
    content: r.content ?? "",
    createdAt: (r.createdat ?? r.CreatedAT ?? r.created_at ?? new Date()).toString(),
  } as Reply;
}

export async function registerView(locusId: string) {
  if (typeof window !== "undefined" && !localStorage.getItem("accessToken")) {
    return { ok: false };
  }
  const res = await fetch(`${API_BASE}/api/loci/${locusId}/view`, {
    method: "POST",
    headers: { ...authHeaders() },
  });
  return { ok: res.ok };
}

export async function createLocus(message: string, coords: { lat: number; lng: number }): Promise<Locus> {
  const res = await fetch(`${API_BASE}/api/post/loci`, {
    method: "POST",
    headers: { "Content-Type": "application/json", ...authHeaders() },
    body: JSON.stringify({
      message,
      location: { lat: coords.lat, long: coords.lng },
    }),
  });
  if (!res.ok) {
    const msg = await res.text();
    throw new Error(msg || "Failed to create locus (are you logged in?)");
  }
  const data = await res.json();
  const row = Array.isArray(data) ? data[0] : data;
  return toLocus(row);
}
