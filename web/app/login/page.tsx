"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { TopBar } from "../../components/TopBar";
import { login } from "../../lib/api";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const router = useRouter();

  const handleLogin = async () => {
    setBusy(true);
    setError("");
    try {
      if (!email || !password) {
        setError("Email and password are required");
        return;
      }
      await login(email, password);
      router.push("/map");
    } catch (e) {
      const msg = e instanceof Error ? e.message : "Login failed. Try again.";
      setError(msg);
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="layout-shell" style={{ minHeight: "100vh" }}>
      <TopBar />
      <main
        style={{
          display: "grid",
          placeItems: "center",
          padding: "32px 18px 64px",
          background: "radial-gradient(circle at 50% 10%, rgba(47,123,255,0.08), transparent 42%), var(--bg)",
        }}
      >
        <div
          className="panel"
          style={{ maxWidth: 420, width: "100%", display: "grid", gap: 12, padding: 22, boxShadow: "0 20px 50px rgba(42,18,47,0.08)" }}
        >
          <div>
            <div className="badge" style={{ border: "1px solid var(--accent)", color: "var(--accent)" }}>
              Welcome back
            </div>
            <h2 style={{ margin: "6px 0 0" }}>Sign in</h2>
            <p className="muted" style={{ margin: 0 }}>Drop messages, view threads, stay close to what's around you.</p>
          </div>
          <label style={{ display: "grid", gap: 6 }}>
            <span className="muted">Email</span>
            <input
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
              style={inputStyle}
              type="email"
            />
          </label>
          <label style={{ display: "grid", gap: 6 }}>
            <span className="muted">Password</span>
            <input
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              style={inputStyle}
              type="password"
            />
          </label>
          {error && <div style={{ color: "var(--danger)", fontWeight: 600 }}>{error}</div>}
          <button className="primary-btn" onClick={handleLogin} disabled={busy}>
            {busy ? "Checking..." : "Login"}
          </button>
          <div className="muted" style={{ fontSize: 14 }}>
            New here? <Link href="/signup" className="badge" style={{ border: "1px solid var(--accent)", color: "var(--accent)" }}>Create an account</Link>
          </div>
        </div>
      </main>
    </div>
  );
}

const inputStyle: React.CSSProperties = {
  background: "#fff",
  border: "1px solid var(--border)",
  borderRadius: 12,
  padding: "12px 14px",
  color: "var(--text)",
};
