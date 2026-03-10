"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { TopBar } from "../../../components/TopBar";
import { MapView } from "../../../components/MapView";
import { ReplyThread } from "../../../components/ReplyThread";
import { fetchLocus } from "../../../lib/api";
import { Locus, Reply } from "../../../lib/types";

export default function LocusDetail({ params }: { params: { id: string } }) {
  const [locus, setLocus] = useState<Locus | null>(null);
  const [replies, setReplies] = useState<Reply[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const data = await fetchLocus(params.id);
        if (!cancelled) {
          setLocus(data.locus);
          setReplies(data.replies);
        }
      } catch (err: any) {
        if (!cancelled) setError(err?.message ?? "Failed to load locus");
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [params.id]);

  if (error) {
    return (
      <div className="layout-shell">
        <TopBar />
        <div className="page-shell">
          <div className="panel" style={{ color: "var(--danger)", margin: 0 }}>
            {error}
          </div>
          <Link href="/map" className="badge" style={{ border: "1px solid var(--accent)", marginTop: 12, display: "inline-flex", width: "fit-content" }}>
            ← Back to map
          </Link>
        </div>
      </div>
    );
  }

  if (!locus) {
    return (
      <div className="layout-shell">
        <TopBar />
        <div className="page-shell">
          <div className="panel">Loading message…</div>
        </div>
      </div>
    );
  }

  return (
    <div className="layout-shell">
      <TopBar />
      <header style={{ padding: "12px 18px" }}>
        <Link href="/map" className="badge" style={{ border: "1px solid var(--accent)" }}>
          ← Back to map
        </Link>
      </header>
      <main className="main-grid">
        <div className="panel" style={{ display: "grid", gap: 10 }}>
          <div className="pill-row">
            <span className="badge">👁 {locus.viewCount} views</span>
            <span className="badge">💬 {locus.repliesCount} replies</span>
            <span className="badge">{new Date(locus.createdAt).toLocaleString()}</span>
          </div>
          <h2 style={{ margin: 0 }}>{locus.message}</h2>
          <div className="muted" style={{ fontSize: 13 }}>
            {locus.coords.lat.toFixed(4)} / {locus.coords.lng.toFixed(4)}
          </div>
          <div className="map-shell" style={{ minHeight: 260, padding: 0 }}>
            <MapView loci={[locus]} />
          </div>
          <div className="pill-row">
            <span className="stat">Reach: 2km</span>
            <span className="stat">Signal: Strong</span>
          </div>
        </div>
        <ReplyThread locusId={locus.id as string} key={locus.id as string} />
      </main>
    </div>
  );
}
