import type { Metadata } from 'next';
import '@/styles/globals.css';

export const metadata: Metadata = {
  metadataBase: new URL('https://distill.dotbrains.io'),
  title: 'distill — AI-Powered Knowledge Compactor for Agents',
  description: 'Compact technical books into agent-optimized markdown. Template-driven AI compaction with incremental updates, multi-source ingestion, and pluggable agents.',
  openGraph: {
    title: 'distill — AI-Powered Knowledge Compactor for Agents',
    description: 'Compact technical books into agent-optimized markdown. Template-driven AI compaction with incremental updates, multi-source ingestion, and pluggable agents.',
    url: 'https://distill.dotbrains.io',
    siteName: 'distill',
    images: [
      {
        url: '/og-image.svg',
        width: 1200,
        height: 630,
        alt: 'distill — AI-Powered Knowledge Compactor for Agents',
      },
    ],
    locale: 'en_US',
    type: 'website',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'distill — AI-Powered Knowledge Compactor for Agents',
    description: 'Compact technical books into agent-optimized markdown.',
    images: ['/og-image.svg'],
  },
  icons: {
    icon: [
      {
        url: '/favicon.svg',
        type: 'image/svg+xml',
      },
    ],
    apple: [
      {
        url: '/favicon.svg',
        type: 'image/svg+xml',
      },
    ],
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <head>
        <meta charSet="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      </head>
      <body>{children}</body>
    </html>
  );
}
