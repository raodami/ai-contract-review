import './globals.css';
import type { Metadata } from 'next';
import Link from 'next/link';

export const metadata: Metadata = {
  title: 'AI Contract Review - Professional Contract Analysis',
  description: 'AI-powered contract analysis with risk detection and compliance checking',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body style={{ background: '#0a0f1a', color: '#e2e8f0', minHeight: '100vh' }}>
        <nav style={{ padding: '16px 24px', borderBottom: '1px solid rgba(255,255,255,0.05)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Link href="/" style={{ color: '#533afd', fontSize: 20, fontWeight: 700, textDecoration: 'none' }}>
            ⚖️ ContractAI
          </Link>
          <div style={{ display: 'flex', gap: 20 }}>
            <Link href="/" style={{ color: '#8899a6', textDecoration: 'none' }}>Home</Link>
            <Link href="/pricing" style={{ color: '#8899a6', textDecoration: 'none' }}>Pricing</Link>
            <Link href="/dashboard" style={{ color: '#533afd', textDecoration: 'none', fontWeight: 500 }}>Dashboard</Link>
            <Link href="/login" style={{ background: '#533afd', color: '#fff', padding: '8px 20px', borderRadius: 6, textDecoration: 'none' }}>Login</Link>
          </div>
        </nav>
        {children}
        <footer style={{ padding: '20px', textAlign: 'center', color: '#8899a6', fontSize: 13, borderTop: '1px solid rgba(255,255,255,0.05)', marginTop: 40 }}>
          © 2026 AI Contract Review. All rights reserved.
        </footer>
      </body>
    </html>
  );
}
