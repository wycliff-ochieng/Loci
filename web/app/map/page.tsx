"use client";

import { useEffect, useState } from "react";
import { DropLocus } from "../../components/DropLocus";
import { LociFeed } from "../../components/LociFeed";
import { MapPanel } from "../../components/MapPanel";
import { TopBar } from "../../components/TopBar";
import { fetchLoci } from "../../lib/api";
import { Locus } from "../../lib/types";

export default function MapPage() {
  const [loci, setLoci] = useState<Locus[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const data = await fetchLoci();
        if (!cancelled) {
          setLoci(data);
          setError(null);
        }
      } catch (err: any) {
        if (!cancelled) setError(err?.message ?? "Failed to fetch messages");
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <div className="layout-shell">
      <TopBar />
      <div className="page-shell" style={{ paddingTop: 20 }}>
        <h1 className="section-title" style={{ marginTop: 12 }}>Live map + feed</h1>
        <p className="section-subtitle">Drop a message, see it on the map, and follow the replies.</p>
        {error && (
          <div className="panel" style={{ color: "var(--danger)", fontWeight: 600 }}>
            {error} — ensure the API is running at your NEXT_PUBLIC_API_BASE and you are logged in for posting.
          </div>
        )}
        {loading && !error && <div className="panel">Loading nearby messages…</div>}
        {!loading && (
          <div className="live-grid">
            <MapPanel loci={loci} />
            <div style={{ display: "grid", gap: 14 }}>
              <DropLocus seedCoords={loci[0]?.coords} onCreated={(l) => setLoci((prev) => [l, ...prev])} />
              <LociFeed loci={loci} />
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
