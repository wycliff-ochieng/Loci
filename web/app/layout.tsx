import type { Metadata } from "next";
import "./globals.css";
import "leaflet/dist/leaflet.css";

export const metadata: Metadata = {
  title: "Loci | Geo Messaging",
  description: "Map-first messaging to drop, view, and reply to loci near you.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="layout-shell">{children}</body>
    </html>
  );
}
