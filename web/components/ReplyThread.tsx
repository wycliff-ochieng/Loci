"use client";

import { useEffect, useState } from "react";
import { createReply, fetchReplies, registerView } from "../lib/api";
import { Reply } from "../lib/types";
import { Composer } from "./Composer";

export function ReplyThread({ locusId }: { locusId: string }) {
  const [replies, setReplies] = useState<Reply[]>([]);

  useEffect(() => {
    fetchReplies(locusId).then(setReplies);
    registerView(locusId).catch(() => undefined);
  }, [locusId]);

  const handleSubmit = async (value: string) => {
    const created = await createReply(locusId, value);
    setReplies((prev) => [created, ...prev]);
  };

  return (
    <div className="panel" style={{ display: "grid", gap: 12 }}>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <div>
          <div style={{ fontWeight: 700 }}>Thread</div>
          <div className="muted" style={{ fontSize: 13 }}>
            New replies land at the top.
          </div>
        </div>
        <span className="badge">Live</span>
      </div>
      <Composer placeholder="Drop a reply" onSubmit={handleSubmit} />
      <div className="feed-list">
        {replies.map((reply) => (
          <div key={reply.id} className="reply-card">
            <div style={{ fontWeight: 700, marginBottom: 6 }}>{reply.content}</div>
            <div className="muted" style={{ fontSize: 12 }}>
              {new Date(reply.createdAt).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
