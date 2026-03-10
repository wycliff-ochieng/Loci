// Client-side view batching and fire-and-forget delivery.
// Collects locus IDs viewed on screen and flushes them periodically.

import { registerView } from "./api";

const pending = new Set<string>();
const FLUSH_INTERVAL_MS = 5000;
const MAX_BATCH_BEFORE_IMMEDIATE_FLUSH = 5;
let flushTimer: ReturnType<typeof setTimeout> | null = null;
let flushing = false;

async function flush() {
  if (pending.size === 0 || flushing) return;
  flushing = true;
  const ids = Array.from(pending);
  pending.clear();

  try {
    // Fire-and-forget semantics: best effort; ignore errors.
    await Promise.all(ids.map((id) => registerView(id)));
  } catch {
    // Swallow errors; next flush can retry new items.
  } finally {
    flushing = false;
  }
}

function scheduleFlush() {
  if (flushTimer) return;
  flushTimer = setTimeout(() => {
    flushTimer = null;
    flush();
  }, FLUSH_INTERVAL_MS);
}

export function enqueueView(locusId: string) {
  pending.add(locusId);
  if (pending.size >= MAX_BATCH_BEFORE_IMMEDIATE_FLUSH) {
    flush();
    if (flushTimer) {
      clearTimeout(flushTimer);
      flushTimer = null;
    }
  } else {
    scheduleFlush();
  }
}
