"use client";

import { useEffect, useState } from "react";
import { createLocus } from "../lib/api";
import { Locus } from "../lib/types";

export function DropLocus({
  seedCoords,
  onCreated,
}: {
  seedCoords?: { lat: number; lng: number };
  onCreated: (locus: Locus) => void;
}) {
  const fallback = seedCoords ?? { lat: -1.2864, lng: 36.8172 };
  const [message, setMessage] = useState("");
  const [coords, setCoords] = useState(fallback);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [info, setInfo] = useState<string | null>(null);

  useEffect(() => {
    if (!navigator?.geolocation) return;
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        setCoords({ lat: pos.coords.latitude, lng: pos.coords.longitude });
      },
      () => {
        /* ignore errors, fallback remains */
      },
      { enableHighAccuracy: false, timeout: 3000 }
    );
  }, []);

  const handleUseCurrent = () => {
    if (!navigator?.geolocation) {
      setError("Geolocation not available");
      return;
    }
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        setCoords({ lat: pos.coords.latitude, lng: pos.coords.longitude });
        setError(null);
        setInfo("Pinned to your current location");
      },
      () => setError("Unable to get your location"),
      { enableHighAccuracy: true, timeout: 5000 }
    );
  };

  const handleSubmit = async () => {
    setBusy(true);
    setError(null);
    setInfo(null);
    try {
      if (!message.trim()) {
        setError("Message is required");
        return;
      }
      const locus = await createLocus(message.trim(), coords);
      setMessage("");
      onCreated(locus);
      setInfo("Dropped!");
    } catch (e) {
      const msg = e instanceof Error ? e.message : "Failed to drop locus";
      setError(msg);
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="panel" style={{ display: "grid", gap: 10 }}>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <div>
           <div style={{ fontWeight: 700 }}>Drop a message on the map</div>
          <div className="muted" style={{ fontSize: 13 }}>
            Message + location. You need to be logged in.
          </div>
        </div>
        <button className="badge" onClick={handleUseCurrent} style={{ border: "1px solid rgba(255,255,255,0.12)", background: "transparent" }}>
          Use my location
        </button>
      </div>
      <textarea
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        placeholder="What's happening here?"
        style={{
          width: "100%",
          minHeight: 80,
          borderRadius: 12,
          padding: "12px 14px",
          background: "#fff",
          border: "1px solid var(--border)",
          color: "var(--text)",
        }}
      />
      <div className="pill-row">
        <span className="stat">Lat: {coords.lat.toFixed(4)}</span>
        <span className="stat">Lng: {coords.lng.toFixed(4)}</span>
      </div>
      {error && <div style={{ color: "var(--danger)", fontWeight: 600 }}>{error}</div>}
      {info && !error && <div style={{ color: "var(--accent)", fontWeight: 600 }}>{info}</div>}
      <button className="primary-btn" onClick={handleSubmit} disabled={busy}>
        {busy ? "Sending..." : "Drop locus"}
      </button>
    </div>
  );
}
