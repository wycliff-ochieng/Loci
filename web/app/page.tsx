import Link from "next/link";
import { TopBar } from "../components/TopBar";

const FEATURES = [
  { title: "One map inbox", desc: "All nearby chats in a single, map-first timeline. No switching apps to see what is happening around you.", icon: "🗺️" },
  { title: "Secure by location", desc: "Messages are tied to coordinates, not threads lost in DMs. Your view stays relevant to where you are.", icon: "🔒" },
  { title: "Realtime replies", desc: "Threads update instantly so you can coordinate meetups, alerts, and quick check-ins without lag.", icon: "⚡" },
  { title: "Works everywhere", desc: "Built for desktop and mobile. Sign in from any device and keep the same live map experience.", icon: "📱" },
];

export default function Home() {
  return (
    <div className="layout-shell">
      <TopBar />
      <div className="page-shell">
        <section className="hero">
          <div>
            <div className="pill">Geo messaging · Map inbox</div>
            <h1>One map inbox for all your chats.</h1>
            <p>
              Drop and read location-aware messages in a single view. Everything nearby shows up together, with replies
              that stay anchored to the place they matter.
            </p>
            <div className="cta-row">
              <Link href="/signup" className="primary-btn">Get started</Link>
              <Link href="/map" className="outline-btn">Open live map</Link>
              <Link href="/login" className="outline-btn">Login</Link>
            </div>
            <div className="pill-row" style={{ marginTop: 14 }}>
              <span className="badge">👁 View counts stay local</span>
              <span className="badge">💬 Threads pinned to place</span>
              <span className="badge">📡 Ready for realtime fan-out</span>
            </div>
          </div>
          <div className="hero-card">
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
              <div>
                <div className="badge" style={{ border: "1px solid var(--accent)", color: "var(--accent)" }}>
                  Live snapshot
                </div>
                <div style={{ fontWeight: 700, marginTop: 6 }}>What you see on the ground</div>
              </div>
              <div className="pill" aria-hidden>
                ⏱ Instant
              </div>
            </div>
            <div style={{ display: "grid", gap: 8 }}>
              <div className="stat">Pinned to place</div>
              <div className="stat" style={{ background: "#eef4ff", color: "var(--accent)" }}>
                Drop and reply without losing context
              </div>
              <div className="muted" style={{ fontSize: 14 }}>
                Head to the live map to see real threads, drop your own, and keep conversations tied to where they matter.
              </div>
            </div>
          </div>
        </section>

        <h2 className="section-title">All your nearby chats in one place</h2>
        <p className="section-subtitle">Designed to keep conversations grounded to where they happen.</p>
        <div className="feature-grid">
          {FEATURES.map((f) => (
            <div key={f.title} className="feature-card">
              <div className="pill" aria-hidden>
                {f.icon}
              </div>
              <div style={{ fontWeight: 700, marginTop: 6 }}>{f.title}</div>
              <div className="muted" style={{ fontSize: 14 }}>{f.desc}</div>
            </div>
          ))}
        </div>

        <div className="cta-banner">
          <div>
            <h3>Ready to drop your first message?</h3>
            <div style={{ opacity: 0.9 }}>Create an account, pick a spot on the map, and keep every chat tied to place.</div>
          </div>
          <div className="cta-row" style={{ margin: 0, justifyContent: "flex-end" }}>
            <Link href="/signup" className="primary-btn" style={{ boxShadow: "0 12px 30px rgba(0,0,0,0.1)" }}>
              Get started free
            </Link>
            <Link href="/login" className="outline-btn" style={{ background: "rgba(255,255,255,0.2)", color: "#fff", borderColor: "rgba(255,255,255,0.4)" }}>
              Login
            </Link>
          </div>
        </div>

        <section>
          <h2 className="section-title">About GeoMessages</h2>
          <p className="section-subtitle" style={{ maxWidth: 820, margin: "0 auto 24px" }}>
            GeoMessages is a standalone, map-first messaging experience. We keep chats anchored to the real places they
            belong, so teams, neighbors, and communities can see what matters around them without sifting through endless
            feeds. Use it from any device—your live map inbox follows you everywhere.
          </p>
        </section>
      </div>
    </div>
  );
}
