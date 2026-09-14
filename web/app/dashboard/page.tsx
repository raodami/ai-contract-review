'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

export default function Dashboard() {
  const router = useRouter();
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);
  const [jobs, setJobs] = useState([]);
  const [result, setResult] = useState(null);
  const [selectedJob, setSelectedJob] = useState(null);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/login');
      return;
    }
    fetch('/api/auth/me', { headers: { Authorization: `Bearer ${token}` } })
      .then(r => r.json())
      .then(data => {
        setUser(data);
        loadJobs();
      })
      .catch(() => router.push('/login'))
      .finally(() => setLoading(false));
  }, []);

  const loadJobs = async () => {
    const token = localStorage.getItem('token');
    try {
      const res = await fetch('/api/contract/jobs', {
        headers: { Authorization: `Bearer ${token}` }
      });
      if (res.ok) {
        const data = await res.json();
        setJobs(data.jobs || []);
      }
    } catch (e) {
      console.error('Failed to load jobs', e);
    }
  };

  const handleUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    const token = localStorage.getItem('token');
    const formData = new FormData();
    formData.append('file', file);

    setUploading(true);
    try {
      const res = await fetch('/api/contract/upload', {
        method: 'POST',
        headers: { Authorization: `Bearer ${token}` },
        body: formData,
      });
      const data = await res.json();
      if (res.ok) {
        loadJobs();
        alert('Upload started! Check results in Job History.');
      } else {
        alert('Error: ' + (data.error || 'Unknown error'));
      }
    } catch (err) {
      alert('Upload failed: ' + err.message);
    } finally {
      setUploading(false);
    }
  };

  const fetchJobResult = async (jobId) => {
    const token = localStorage.getItem('token');
    try {
      const res = await fetch(`/api/contract/jobs/${jobId}`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      if (res.ok) {
        const job = await res.json();
        setSelectedJob(job);
        if (job.result) {
          setResult(JSON.parse(job.result));
        }
      }
    } catch (e) {
      console.error('Failed to fetch job', e);
    }
  };

  if (loading) {
    return (
      <div style={{ padding: 40, textAlign: 'center', color: '#8899a6' }}>
        Loading...
      </div>
    );
  }

  return (
    <div style={{ padding: 40, maxWidth: 900, margin: '0 auto' }}>
      <h1 style={{ marginBottom: 32, fontSize: 28 }}>Contract Review Dashboard</h1>

      {/* Upload Section */}
      <div style={{ background: 'rgba(255,255,255,0.05)', borderRadius: 12, padding: 24, marginBottom: 24 }}>
        <h3 style={{ marginBottom: 16 }}>Upload Contract (.docx)</h3>
        <input
          type="file"
          accept=".docx"
          onChange={handleUpload}
          disabled={uploading}
          style={{ marginBottom: 16, display: 'block', color: '#8899a6' }}
        />
        {uploading && <p style={{ color: '#533afd' }}>Processing...</p>}
        <p style={{ color: '#8899a6', fontSize: 14 }}>
          Max 10MB. AI analysis will run automatically.
        </p>
      </div>

      {/* Result Section */}
      {result && (
        <div style={{ background: 'rgba(255,255,255,0.05)', borderRadius: 12, padding: 24, marginBottom: 24 }}>
          <h3 style={{ marginBottom: 16 }}>Analysis Result</h3>
          <div style={{ marginBottom: 16 }}>
            <span style={{
              fontSize: 48,
              fontWeight: 700,
              color: result.score >= 80 ? '#22c55e' : result.score >= 50 ? '#f59e0b' : '#ef4444',
            }}>
              {result.score}/100
            </span>
            <span style={{ color: '#8899a6', marginLeft: 12 }}>Risk Score</span>
          </div>
          {result.summary && (
            <div style={{ marginBottom: 16, whiteSpace: 'pre-wrap', color: '#e2e8f0' }}>
              {result.summary}
            </div>
          )}
          {result.risks && result.risks.length > 0 && (
            <div>
              <h4 style={{ marginBottom: 8, color: '#f8fafc' }}>Detected Risks:</h4>
              {result.risks.map((risk, i) => (
                <div key={i} style={{
                  padding: 12,
                  marginBottom: 8,
                  borderRadius: 8,
                  background: risk.risk_level === 'high' ? 'rgba(239,68,68,0.1)' :
                              risk.risk_level === 'medium' ? 'rgba(245,158,11,0.1)' : 'rgba(34,197,94,0.1)',
                  borderLeft: `4px solid ${risk.risk_level === 'high' ? '#ef4444' : risk.risk_level === 'medium' ? '#f59e0b' : '#22c55e'}`,
                }}>
                  <strong>{risk.clause}</strong>
                  <p style={{ margin: '4px 0', fontSize: 13, color: '#8899a6' }}>{risk.explanation}</p>
                </div>
              ))}
            </div>
          )}
          {result.suggestions && result.suggestions.length > 0 && (
            <div>
              <h4 style={{ marginBottom: 8, color: '#f8fafc' }}>Suggestions:</h4>
              <ul style={{ color: '#e2e8f0', paddingLeft: 20 }}>
                {result.suggestions.map((s, i) => <li key={i}>{s}</li>)}
              </ul>
            </div>
          )}
        </div>
      )}

      {/* Job History */}
      <div style={{ background: 'rgba(255,255,255,0.05)', borderRadius: 12, padding: 24 }}>
        <h3 style={{ marginBottom: 16 }}>Job History</h3>
        {jobs.length === 0 ? (
          <p style={{ color: '#8899a6' }}>No contracts analyzed yet. Upload one above!</p>
        ) : (
          <div>
            {jobs.map(job => (
              <div
                key={job.id}
                onClick={() => fetchJobResult(job.id)}
                style={{
                  padding: 12,
                  marginBottom: 8,
                  borderRadius: 8,
                  background: 'rgba(255,255,255,0.03)',
                  cursor: 'pointer',
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                }}
              >
                <div>
                  <div style={{ color: '#e2e8f0', fontWeight: 500 }}>{job.file_name}</div>
                  <div style={{ color: '#8899a6', fontSize: 12 }}>
                    {new Date(job.created_at * 1000).toLocaleString()} · {job.file_size} bytes
                  </div>
                </div>
                <span style={{
                  padding: '4px 12px',
                  borderRadius: 20,
                  fontSize: 12,
                  background: job.status === 'completed' ? 'rgba(34,197,94,0.2)' :
                              job.status === 'processing' ? 'rgba(59,130,246,0.2)' : 'rgba(239,68,68,0.2)',
                  color: job.status === 'completed' ? '#22c55e' :
                         job.status === 'processing' ? '#3b82f6' : '#ef4444',
                }}>
                  {job.status}
                </span>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* User Info */}
      {user && (
        <div style={{ marginTop: 24, color: '#8899a6', fontSize: 14 }}>
          Logged in as {user.email} ·{' '}
          {user.is_pro ? <span style={{ color: '#533afd' }}>Pro Plan</span> : 'Free Plan'}
        </div>
      )}
    </div>
  );
}
