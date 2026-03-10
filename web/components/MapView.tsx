"use client";

import { CircleMarker, MapContainer, TileLayer, Tooltip } from "react-leaflet";
import { Locus } from "../lib/types";

function computeIntensity(viewCount: number) {
  const base = 0.25 + Math.log10(viewCount + 1) * 0.35;
  return Math.min(Math.max(base, 0.2), 1);
}

export function MapView({ loci }: { loci: Locus[] }) {
  const first = loci[0]?.coords;
  const center = first && (first.lat !== 0 || first.lng !== 0) ? first : { lat: 30.2672, lng: -97.7431 };

  return (
    <MapContainer
      center={[center.lat, center.lng]}
      zoom={14}
      style={{ height: 420, width: "100%", borderRadius: 14, overflow: "hidden" }}
      attributionControl={false}
    >
      <TileLayer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        attribution="&copy; OpenStreetMap contributors"
      />
      {loci.map((locus) => {
        const intensity = computeIntensity(locus.viewCount);
        const radius = 8 + intensity * 12;
        const color = "#2f7bff";
        const fillOpacity = 0.3 + intensity * 0.4;
        return (
          <CircleMarker
            key={locus.id}
            center={[locus.coords.lat, locus.coords.lng]}
            radius={radius}
            pathOptions={{
              color,
              fillColor: color,
              fillOpacity,
              weight: 1,
            }}
          >
            <Tooltip direction="top" offset={[0, -4]} opacity={0.9} sticky>
              <div style={{ maxWidth: 240 }}>
                <div style={{ fontWeight: 700 }}>{locus.message}</div>
                <div style={{ fontSize: 12 }}>👁 {locus.viewCount} · 💬 {locus.repliesCount}</div>
              </div>
            </Tooltip>
          </CircleMarker>
        );
      })}
    </MapContainer>
  );
}
