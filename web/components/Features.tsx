export default function Features() {
  const features = [
    { icon: '🔍', title: 'Risk Detection', desc: 'Automatically identify high-risk clauses and potential legal issues in your contracts.' },
    { icon: '⚖️', title: 'Compliance Checks', desc: 'Ensure your contracts meet legal requirements with automated compliance validation.' },
    { icon: '📊', title: 'Smart Analysis', desc: 'AI-powered classification of contract types and detailed clause analysis.' },
    { icon: '💡', title: 'Actionable Insights', desc: 'Get clear recommendations to improve contract terms and reduce risk.' },
  ];

  return (
    <section style={{ padding: '60px 20px', maxWidth: 900, margin: '0 auto' }}>
      <h2 style={{ textAlign: 'center', marginBottom: 48, color: '#f8fafc', fontSize: 32 }}>Features</h2>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: 24 }}>
        {features.map((f, i) => (
          <div key={i} style={{ padding: 24, background: 'rgba(255,255,255,0.03)', borderRadius: 12, border: '1px solid rgba(255,255,255,0.05)' }}>
            <div style={{ fontSize: 32, marginBottom: 16 }}>{f.icon}</div>
            <h3 style={{ color: '#f8fafc', marginBottom: 8 }}>{f.title}</h3>
            <p style={{ color: '#8899a6', fontSize: 14, lineHeight: 1.6 }}>{f.desc}</p>
          </div>
        ))}
      </div>
    </section>
  );
}
