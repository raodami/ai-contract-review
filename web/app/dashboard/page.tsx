'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

export default function Dashboard() {
  const router = useRouter();
  const [user, setUser] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);
  const [jobs, setJobs] = useState<any[]>([]);
  const [result, setResult] = useState<any>(null);
  const [selectedJob, setSelectedJob] = useState<any>(null);
  const [analyzing, setAnalyzing] = useState(false);

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

  const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
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
        // Poll for completion
        const jobId = data.job_id;
        if (jobId && token) {
          setTimeout(() => checkJobStatus(jobId, token), 2000);
        }
      } else {
        alert('Error: ' + (data.error || 'Unknown error'));
      }
    } catch (err: any) {
      alert('Upload failed: ' + (err.message || String(err)));
    } finally {
      setUploading(false);
    }
  };

  const checkJobStatus = async (jobId: string, token: string) => {
    try {
      const res = await fetch(`/api/contract/jobs/${jobId}`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      if (res.ok) {
        const job = await res.json();
        if (job.status === 'completed') {
          setSelectedJob(job);
          setResult(JSON.parse(job.result));
          loadJobs();
          return;
        }
        if (job.status === 'processing') {
          setTimeout(() => checkJobStatus(jobId, token), 2000);
        }
      }
    } catch (e) {
      console.error('Failed to check job status', e);
    }
  };

  const fetchJobResult = async (jobId: string) => {
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

  const analyzeText = async () => {
    const token = localStorage.getItem('token');
    if (!token) return;
    
    const text = prompt('Enter contract text to analyze:');
    if (!text) return;
    
    setAnalyzing(true);
    try {
      const res = await fetch('/api/contract/analyze', {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}` 
        },
        body: JSON.stringify({ text }),
      });
      if (res.ok) {
        const data = await res.json();
        setResult(data);
      }
    } catch (e) {
      console.error('Analysis failed', e);
    } finally {
      setAnalyzing(false);
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
        <h3 style={{ marginBottom: 16 }}>Upload Contract (.docx / .pdf)</h3>
        <input
          type="file"
          accept=".docx,.pdf"
          onChange={handleUpload}
          disabled={uploading}
          style={{ marginBottom: 16, display: 'block', color: '#8899a6' }}
        />
        {uploading && <p style={{ color: '#533afd' }}>Processing...</p>}
        <p style={{ color: '#8899a6', fontSize: 14 }}>
          Max 10MB. AI analysis runs automatically.
        </p>
      </div>

      {/* Text Analysis */}
      <div style={{ background: 'rgba(255,255,255,0.05)', borderRadius: 12, padding: 24, marginBottom: 24 }}>
        <h3 style={{ marginBottom: 16 }}>Analyze Text</h3>
        <button
          onClick={analyzeText}
          disabled={analyzing}
          style={{
            padding: '10px 20px',
            background: '#533afd',
            color: '#fff',
            border: 'none',
            borderRadius: 6,
            cursor: analyzing ? 'not-allowed' : 'pointer',
          }}
        >
          {analyzing ? 'Analyzing...' : 'Analyze Sample Text'}
        </button>
      </div>

      {/* Result Section */}
      {result && (
        <div style={{ background: 'rgba(255,255,255,0.05)', borderRadius: 12, padding: 24, marginBottom: 24 }}>
          <h3 style={{ marginBottom: 16 }}>Analysis Result</h3>
          
          {/* Score */}
          <div style={{ marginBottom: 24, textAlign: 'center' }}>
            <span style={{
              fontSize: 64,
              fontWeight: 700,
              color: result.score >= 80 ? '#22c55e' : result.score >= 50 ? '#f59e0b' : '#ef4444',
            }}>
              {result.score}/100
            </span>
            <div style={{ color: '#8899a6', marginTop: 8 }}>Risk Score</div>
          </div>

          {/* Summary */}
          {result.summary && (
            <div style={{ marginBottom: 24, padding: 16, background: 'rgba(0,0,0,0.2)', borderRadius: 8 }}>
              <h4 style={{ marginBottom: 8, color: '#f8fafc' }}>Summary</h4>
              <p style={{ color: '#e2e8f0', whiteSpace: 'pre-wrap', lineHeight: 1.6 }}>{result.summary}</p>
            </div>
          )}

          {/* Risks */}
          {result.risks && result.risks.length > 0 && (
            <div style={{ marginBottom: 24 }}>
              <h4 style={{ marginBottom: 12, color: '#f8fafc' }}>Detected Risks ({result.risks.length})</h4>
              {result.risks.map((risk: any, i: number) => (
                <div key={i} style={{
                  padding: 16,
                  marginBottom: 8,
                  borderRadius: 8,
                  background: risk.risk_level === 'high' ? 'rgba(239,68,68,0.1)' :
                              risk.risk_level === 'medium' ? 'rgba(245,158,11,0.1)' : 'rgba(34,197,94,0.1)',
                  borderLeft: `4px solid ${risk.risk_level === 'high' ? '#ef4444' : risk.risk_level === 'medium' ? '#f59e0b' : '#22c55e'}`,
                }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4 }}>
                    <strong style={{ color: '#e2e8f0' }}>{risk.clause}</strong>
                    <span style={{
                      padding: '2px 8px',
                      borderRadius: 4,
                      fontSize: 11,
                      background: risk.risk_level === 'high' ? 'rgba(239,68,68,0.3)' :
                                  risk.risk_level === 'medium' ? 'rgba(245,158,11,0.3)' : 'rgba(34,197,94,0.3)',
                      color: risk.risk_level === 'high' ? '#ef4444' : risk.risk_level === 'medium' ? '#f59e0b' : '#22c55e',
                    }}>
                      {risk.risk_level.toUpperCase()}
                    </span>
                  </div>
                  <p style={{ margin: 0, fontSize: 13, color: '#8899a6' }}>{risk.explanation}</p>
                </div>
              ))}
            </div>
          )}

          {/* Clauses */}
          {result.clauses && result.clauses.length > 0 && (
            <div style={{ marginBottom: 24 }}>
              <h4 style={{ marginBottom: 12, color: '#f8fafc' }}>Classified Clauses ({result.clauses.length})</h4>
              <div style={{ display: 'grid', gap: 8 }}>
                {result.clauses.map((clause: any, i: number) => (
                  <div key={i} style={{
                    padding: 12,
                    borderRadius: 8,
                    background: 'rgba(255,255,255,0.03)',
                    border: '1px solid rgba(255,255,255,0.05)',
                  }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4 }}>
                      <span style={{ color: '#533afd', fontSize: 12, fontWeight: 600 }}>{clause.type}</span>
                      {clause.risk_level && (
                        <span style={{
                          fontSize: 10,
                          padding: '2px 6px',
                          borderRadius: 4,
                          background: clause.risk_level === 'high' ? 'rgba(239,68,68,0.2)' : 'rgba(245,158,11,0.2)',
                          color: clause.risk_level === 'high' ? '#ef4444' : '#f59e0b',
                        }}>
                          {clause.risk_level.toUpperCase()}
                        </span>
                      )}
                    </div>
                    <div style={{ color: '#e2e8f0', fontSize: 13 }}>{clause.content}</div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Compliance */}
          {result.compliance && result.compliance.length > 0 && (
            <div style={{ marginBottom: 24 }}>
              <h4 style={{ marginBottom: 12, color: '#f8fafc' }}>Compliance Check</h4>
              <div style={{ display: 'grid', gap: 8 }}>
                {result.compliance.map((check: any, i: number) => (
                  <div key={i} style={{
                    padding: 12,
                    borderRadius: 8,
                    background: check.status === 'present' ? 'rgba(34,197,94,0.05)' :
                                check.status === 'missing' ? 'rgba(239,68,68,0.05)' : 'rgba(245,158,11,0.05)',
                    border: '1px solid rgba(255,255,255,0.05)',
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                  }}>
                    <div>
                      <div style={{ color: '#e2e8f0', fontSize: 13 }}>{check.item}</div>
                      {check.issue && <div style={{ color: '#8899a6', fontSize: 11, marginTop: 2 }}>{check.issue}</div>}
                    </div>
                    <span style={{
                      padding: '4px 10px',
                      borderRadius: 12,
                      fontSize: 11,
                      fontWeight: 600,
                      background: check.status === 'present' ? 'rgba(34,197,94,0.2)' :
                                  check.status === 'missing' ? 'rgba(239,68,68,0.2)' : 'rgba(245,158,11,0.2)',
                      color: check.status === 'present' ? '#22c55e' : check.status === 'missing' ? '#ef4444' : '#f59e0b',
                    }}>
                      {check.status.toUpperCase()}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Suggestions */}
          {result.suggestions && result.suggestions.length > 0 && (
            <div>
              <h4 style={{ marginBottom: 12, color: '#f8fafc' }}>Suggestions</h4>
              <ul style={{ color: '#e2e8f0', paddingLeft: 20, lineHeight: 1.8 }}>
                {result.suggestions.map((s: string, i: number) => <li key={i}>{s}</li>)}
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
            {jobs.map((job: any) => (
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
