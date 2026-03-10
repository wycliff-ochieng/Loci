"use client";

import dynamic from "next/dynamic";
import { Locus } from "../lib/types";

const MapView = dynamic(() => import("./MapView").then((mod) => mod.MapView), { ssr: false });

export function MapPanel({ loci }: { loci: Locus[] }) {
  return (
    <div className="panel" style={{ minHeight: 480, display: "grid", gap: 14 }}>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <div>
          <div style={{ fontWeight: 700 }}>Live map</div>
          <div className="muted" style={{ fontSize: 13 }}>
            Drop a message, see who is nearby, glance at density.
          </div>
        </div>
        <div className="pill-row">
          <span className="badge">Heat: {loci.length} active</span>
          <span className="badge">Visibility: Public</span>
        </div>
      </div>
      <div className="map-shell" style={{ padding: 0 }}>
        <MapView loci={loci} />
      </div>
    </div>
  );
}
