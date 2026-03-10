"use client";

import Link from "next/link";
import { useEffect, useRef, useState } from "react";
import { enqueueView } from "../lib/viewReporter";
import { Locus } from "../lib/types";

export function LocusCard({ locus }: { locus: Locus }) {
  const created = new Date(locus.createdAt).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  const intensity = Math.min(0.25 + Math.log10(locus.viewCount + 1) * 0.35, 1);
  const glow = `0 0 0 ${4 + intensity * 10}px rgba(47, 123, 255, ${0.12 + intensity * 0.2}), 0 10px 24px rgba(0,0,0,0.06)`;

  const cardRef = useRef<HTMLDivElement | null>(null);
  const [hasReported, setHasReported] = useState(false);

  useEffect(() => {
    const el = cardRef.current;
    if (!el || hasReported) return;

    let dwellTimer: ReturnType<typeof setTimeout> | null = null;

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            if (!dwellTimer) {
              dwellTimer = setTimeout(() => {
                enqueueView(String(locus.id));
                setHasReported(true);
              }, 1200);
            }
          } else if (dwellTimer) {
            clearTimeout(dwellTimer);
            dwellTimer = null;
          }
        });
      },
      { threshold: 0.6 }
    );

    observer.observe(el);

    return () => {
      if (dwellTimer) clearTimeout(dwellTimer);
      observer.disconnect();
    };
  }, [locus.id, hasReported]);

  return (
    <div
      ref={cardRef}
      className="locus-card"
      style={{
        boxShadow: glow,
        borderColor: `rgba(47, 123, 255, ${0.08 + intensity * 0.25})`,
      }}
    >
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <div className="badge">{created}</div>
        <div className="pill-row">
          <span className="stat">👁 {locus.viewCount}</span>
          <span className="stat">💬 {locus.repliesCount}</span>
        </div>
      </div>
      <div style={{ margin: "10px 0", fontWeight: 700 }}>{locus.message}</div>
      <div className="muted" style={{ fontSize: 12 }}>
        {locus.coords.lat.toFixed(4)} / {locus.coords.lng.toFixed(4)}
      </div>
      <div style={{ marginTop: 10 }}>
        <Link href={`/loci/${locus.id}`} className="badge" style={{ border: "1px solid var(--accent)", color: "var(--accent)" }}>
          Open thread
        </Link>
      </div>
    </div>
  );
}
