import Link from 'next/link';

export default function Hero() {
  return (
    <section style={{ padding: '80px 20px', textAlign: 'center' }}>
      <h1 style={{ fontSize: 48, fontWeight: 700, color: '#f8fafc', marginBottom: 20, lineHeight: 1.2 }}>
        AI-Powered Contract<br />Analysis & Review
      </h1>
      <p style={{ fontSize: 18, color: '#8899a6', maxWidth: 600, margin: '0 auto 40px', lineHeight: 1.6 }}>
        Upload your contracts and get instant risk analysis, compliance checks, and AI-powered recommendations.
      </p>
      <div style={{ display: 'flex', gap: 16, justifyContent: 'center' }}>
        <Link href="/dashboard" style={{ padding: '14px 32px', background: '#533afd', color: '#fff', borderRadius: 8, textDecoration: 'none', fontWeight: 500 }}>
          Start Analyzing
        </Link>
        <Link href="/pricing" style={{ padding: '14px 32px', background: 'rgba(255,255,255,0.05)', color: '#e2e8f0', borderRadius: 8, textDecoration: 'none', border: '1px solid rgba(255,255,255,0.1)' }}>
          View Pricing
        </Link>
      </div>
    </section>
  );
}
