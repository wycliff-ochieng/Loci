"use client";

import { useState } from "react";

export function Composer({
  placeholder,
  onSubmit,
}: {
  placeholder: string;
  onSubmit: (value: string) => Promise<void> | void;
}) {
  const [value, setValue] = useState("");
  const [busy, setBusy] = useState(false);

  const handleSend = async () => {
    if (!value.trim()) return;
    setBusy(true);
    await onSubmit(value.trim());
    setValue("");
    setBusy(false);
  };

  return (
    <div className="composer">
      <textarea
        value={value}
        onChange={(e) => setValue(e.target.value)}
        placeholder={placeholder}
      />
      <button className="primary-btn" onClick={handleSend} disabled={busy}>
        {busy ? "Sending..." : "Send"}
      </button>
    </div>
  );
}
