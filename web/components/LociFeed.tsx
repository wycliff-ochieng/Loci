import { Locus } from "../lib/types";
import { LocusCard } from "./LocusCard";

export function LociFeed({ loci }: { loci: Locus[] }) {
  return (
    <div className="panel">
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <div>
          <div style={{ fontWeight: 700 }}>Nearby messages</div>
          <div className="muted" style={{ fontSize: 13 }}>
            Sorted by recency. Tap to open and reply.
          </div>
        </div>
        <span className="badge">Auto-refreshing</span>
      </div>
      <div className="feed-list" style={{ marginTop: 12 }}>
        {loci.map((locus) => (
          <LocusCard key={locus.id} locus={locus} />
        ))}
      </div>
    </div>
  );
}
