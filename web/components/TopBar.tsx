"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { clearAuth } from "../lib/api";

export function TopBar() {
  const pathname = usePathname();
  const router = useRouter();
  const [isAuthed, setIsAuthed] = useState(false);

  useEffect(() => {
    const checkAuth = () => setIsAuthed(!!localStorage.getItem("accessToken"));
    checkAuth();
    window.addEventListener("storage", checkAuth);
    return () => window.removeEventListener("storage", checkAuth);
  }, []);

  const links = [
    { href: "/", label: "Home" },
    { href: "/map", label: "Live map" },
    ...(!isAuthed ? [{ href: "/login", label: "Login" }, { href: "/signup", label: "Sign up" }] : []),
  ];

  const handleLogout = () => {
    clearAuth();
    setIsAuthed(false);
    router.push("/login");
  };
  return (
    <header
      style={{
        padding: "12px 18px",
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
        borderBottom: "1px solid var(--border)",
        backdropFilter: "blur(10px)",
        background: "rgba(255,255,255,0.78)",
        position: "sticky",
        top: 0,
        zIndex: 10,
      }}
    >
      <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
        <div
          style={{
            width: 34,
            height: 34,
            borderRadius: 10,
            background: "linear-gradient(135deg, var(--accent), var(--accent-2))",
            display: "grid",
            placeItems: "center",
            color: "#fff",
            fontWeight: 800,
          }}
        >
          G
        </div>
        <div>
          <div style={{ fontWeight: 700, letterSpacing: 0.2 }}>GeoMessages</div>
          <div className="muted" style={{ fontSize: 12 }}>
            Geo messaging in real-time
          </div>
        </div>
      </div>
      <nav style={{ display: "flex", gap: 12 }}>
        {links.map((link) => {
          const active = pathname === link.href;
          return (
            <Link
              key={link.href}
              href={link.href}
              className="badge"
              style={{
                border: active ? "1px solid var(--accent)" : "1px solid var(--border)",
                color: active ? "var(--accent)" : "var(--muted)",
              }}
            >
              {link.label}
            </Link>
          );
        })}
        {isAuthed && (
          <button
            className="badge"
            onClick={handleLogout}
            style={{ border: "1px solid var(--border)", color: "var(--muted)", background: "#fff" }}
          >
            Logout
          </button>
        )}
      </nav>
      <div className="pill-row">
        <span className="stat">Now: 123 online</span>
        <span className="stat" style={{ background: "#fff0f6", color: "var(--accent)" }}>
          Safe radius: 2km
        </span>
      </div>
    </header>
  );
}
