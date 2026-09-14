'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

export default function Pricing() {
  const router = useRouter();
  const [plans, setPlans] = useState<any[]>([]);
  const [user, setUser] = useState<any>(null);
  const [subscribeLoading, setSubscribeLoading] = useState<string | null>(null);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      fetch('/api/auth/me', { headers: { Authorization: `Bearer ${token}` } })
        .then(r => r.json())
        .then(data => setUser(data))
        .catch(() => {});
    }
    fetch('/api/payment/plans')
      .then(r => r.json())
      .then(data => setPlans(data.plans || []))
      .catch(() => setPlans([
        { id: 'free', name: 'Free', price: 0, minutes: 30, interval: 'monthly', description: '30 minutes/month analysis' },
        { id: 'pro-monthly', name: 'Pro Monthly', price: 990, minutes: 500, interval: 'monthly', description: '500 minutes/month analysis' },
        { id: 'pro-yearly', name: 'Pro Yearly', price: 9900, minutes: 500, interval: 'yearly', description: '500 minutes/month, billed annually - Save 17%' },
      ]));
  }, []);

  const handleSubscribe = async (planId: string) => {
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/login?redirect=/pricing');
      return;
    }
    setSubscribeLoading(planId);
    try {
      const res = await fetch('/api/payment/checkout', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ plan_id: planId }),
      });
      const data = await res.json();
      if (data.checkout_url) {
        window.location.href = data.checkout_url;
      } else {
        // Demo mode - subscribe directly
        const subRes = await fetch('/api/user/subscribe', {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({ plan: planId }),
        });
        if (subRes.ok) {
          setUser({ ...user, is_pro: true });
          alert(`Successfully subscribed to ${planId}!`);
        }
      }
    } catch (e) {
      console.error('Subscription failed', e);
    } finally {
      setSubscribeLoading(null);
    }
  };

  const formatPrice = (price: number) => {
    if (price === 0) return 'Free';
    return `$${(price / 100).toFixed(2)}`;
  };

  return (
    <div style={{ padding: '40px 20px', maxWidth: 900, margin: '0 auto' }}>
      <h1 style={{ marginBottom: 8, fontSize: 32, color: '#f8fafc' }}>Pricing</h1>
      <p style={{ color: '#8899a6', marginBottom: 40 }}>
        Choose the plan that works for you
      </p>

      <div style={{ display: 'grid', gap: 20 }}>
        {plans.map((plan) => (
          <div
            key={plan.id}
            style={{
              padding: 24,
              borderRadius: 12,
              background: plan.id === 'pro-yearly' ? 'linear-gradient(135deg, rgba(83,58,253,0.15) 0%, rgba(6,27,49,0.3) 100%)' : 'rgba(255,255,255,0.03)',
              border: plan.id === 'pro-yearly' ? '1px solid rgba(83,58,253,0.5)' : '1px solid rgba(255,255,255,0.05)',
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              flexWrap: 'wrap',
              gap: 16,
            }}
          >
            <div>
              <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                <h3 style={{ margin: 0, color: '#f8fafc', fontSize: 20 }}>{plan.name}</h3>
                {plan.id === 'pro-yearly' && (
                  <span style={{
                    padding: '2px 10px',
                    borderRadius: 20,
                    fontSize: 11,
                    background: 'rgba(83,58,253,0.3)',
                    color: '#a5b4fc',
                  }}>
                    Save 17%
                  </span>
                )}
              </div>
              <div style={{ margin: '8px 0', color: '#8899a6', fontSize: 14 }}>{plan.description}</div>
              <div style={{ fontSize: 32, fontWeight: 700, color: '#f8fafc' }}>
                {formatPrice(plan.price)}
                <span style={{ fontSize: 14, color: '#8899a6', fontWeight: 400 }}>
                  /{plan.interval === 'yearly' ? 'year' : 'month'}
                </span>
              </div>
            </div>
            <button
              onClick={() => handleSubscribe(plan.id)}
              disabled={subscribeLoading === plan.id}
              style={{
                padding: '12px 28px',
                borderRadius: 8,
                border: 'none',
                background: plan.price === 0 ? 'rgba(255,255,255,0.1)' : '#533afd',
                color: '#fff',
                cursor: subscribeLoading ? 'not-allowed' : 'pointer',
                fontSize: 14,
                fontWeight: 500,
                minWidth: 120,
              }}
            >
              {subscribeLoading === plan.id ? 'Processing...' :
               plan.price === 0 ? 'Get Started' : 'Subscribe'}
            </button>
          </div>
        ))}
      </div>

      {user && (
        <div style={{ marginTop: 32, padding: 16, background: 'rgba(255,255,255,0.03)', borderRadius: 8 }}>
          <div style={{ color: '#8899a6', fontSize: 14 }}>
            Current Plan:{' '}
            <span style={{ color: user.is_pro ? '#533afd' : '#e2e8f0' }}>
              {user.is_pro ? 'Pro' : 'Free'}
            </span>
            {' · '}
            <span style={{ color: '#8899a6' }}>
              Usage: {user.usage_minutes || 0} minutes
            </span>
          </div>
        </div>
      )}

      <div style={{ marginTop: 48, textAlign: 'center' }}>
        <button
          onClick={() => router.push('/dashboard')}
          style={{
            padding: '12px 32px',
            background: 'rgba(255,255,255,0.05)',
            border: '1px solid rgba(255,255,255,0.1)',
            borderRadius: 8,
            color: '#e2e8f0',
            cursor: 'pointer',
          }}
        >
          Back to Dashboard
        </button>
      </div>
    </div>
  );
}
