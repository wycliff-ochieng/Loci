import { UUID } from "crypto";

export type Locus = {
  id: UUID | string;
  userId: UUID | string;
  message: string;
  createdAt: string;
  coords: { lat: number; lng: number };
  viewCount: number;
  repliesCount: number;
};

export type Reply = {
  id: UUID | string;
  locusId: UUID | string;
  userId: UUID | string;
  content: string;
  createdAt: string;
};

export type ViewEvent = {
  locusId: UUID | string;
  userId: UUID | string;
  viewedAt: string;
};

export type GeoBounds = {
  sw: { lat: number; lng: number };
  ne: { lat: number; lng: number };
};
