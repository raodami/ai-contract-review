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
  const [analytics, setAnalytics] = useState<any>(null);

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
        loadAnalytics(token);
      })
      .catch(() => router.push('/login'))
      .finally(() => setLoading(false));
  }, []);

  const loadAnalytics = async (token: string) => {
    try {
      const res = await fetch('/api/analytics', {
        headers: { Authorization: `Bearer ${token}` }
      });
      if (res.ok) {
        const data = await res.json();
        setAnalytics(data);
      }
    } catch (e) {
      console.error('Failed to load analytics', e);
    }
  };

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
        // Handle batch results
        const jobs = data.jobs || [];
        if (jobs.length === 1) {
          const jobId = jobs[0].job_id;
          if (jobId && token) {
            setTimeout(() => checkJobStatus(jobId, token), 2000);
          }
        }
        // Show batch summary
        const completed = jobs.filter((j: any) => j.status === 'completed').length;
        const failed = jobs.filter((j: any) => j.error).length;
        if (completed > 0 || failed > 0) {
          alert(`Batch complete: ${completed} succeeded, ${failed} failed. Remaining quota: ${data.remaining_min} min`);
        }
        loadJobs();
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
        const data = await res.json();
        setResult(data.result ? JSON.parse(data.result) : null);
        loadJobs();
      }
    } catch (e) {
      console.error('Failed to get result', e);
    }
  };

  const fetchJobResult = async (jobId: string) => {
    const token = localStorage.getItem('token');
    try {
      const res = await fetch(`/api/contract/jobs/${jobId}`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      if (res.ok) {
        const data = await res.json();
        setSelectedJob(data);
        setResult(data.result ? JSON.parse(data.result) : null);
      }
    } catch (e) {
      console.error('Failed to fetch job', e);
    }
  };

  const downloadExport = async (format: string) => {
    if (!selectedJob) {
      alert('Please select a job from the history first');
      return;
    }
    const token = localStorage.getItem('token');
    const res = await fetch(`/api/contract/export/${selectedJob.id}?format=${format}`, {
      headers: { Authorization: `Bearer ${token}` }
    });
    if (res.ok) {
      const blob = await res.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = res.headers.get('content-disposition')?.split('filename=')[1] || `report.${format}`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    }
  };

  const analyzeText = async () => {
    const sampleText = "This Employment Agreement is entered into between TechCorp Inc. and John Doe. The Employee shall receive a salary of $75,000 per year, paid bi-weekly. Benefits include health insurance and 20 days PTO. Probation period of 90 days applies. Non-compete clause restricts employment with competitors within 50 miles for 2 years. Either party may terminate with 30 days written notice.";
    setAnalyzing(true);
    setAnalyzing(true);
    try {
      const res = await fetch('/api/contract/upload', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ text: sampleText }),
      });
      if (res.ok) {
        const data = await res.json();
        setResult(data.result);
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
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 32 }}>
        <h1 style={{ fontSize: 28, color: '#f8fafc' }}>Contract Review Dashboard</h1>
        {user && (
          <div style={{ textAlign: 'right' }}>
            <div style={{ color: user.is_pro ? '#533afd' : '#e2e8f0', fontSize: 14, fontWeight: 500 }}>
              {user.is_pro ? '★ Pro Plan' : 'Free Plan'}
            </div>
            <div style={{ color: '#8899a6', fontSize: 12 }}>
              {user.remaining_min ?? 30} min remaining
            </div>
          </div>
        )}
      </div>

      {/* Analytics Overview */}
      {analytics && (
        <div style={{ marginBottom: 24 }}>
          <h3 style={{ marginBottom: 16, color: '#f8fafc' }}>Usage Analytics</h3>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: 16 }}>
            {/* Stats Cards */}
            <div style={{ background: 'rgba(255,255,255,0.05)', borderRadius: 12, padding: 20 }}>
              <div style={{ color: '#8899a6', fontSize: 13, marginBottom: 8 }}>Total Jobs</div>
              <div style={{ fontSize: 32, fontWeight: 700, color: '#533afd' }}>{analytics.total_jobs}</div>
            </div>
            <div style={{ background: 'rgba(255,255,255,0.05)', borderRadius: 12, padding: 20 }}>
              <div style={{ color: '#8899a6', fontSize: 13, marginBottom: 8 }}>Completed</div>
              <div style={{ fontSize: 32, fontWeight: 700, color: '#22c55e' }}>{analytics.completed_jobs}</div>
            </div>
            <div style={{ background: 'rgba(255,255,255,0.05)', borderRadius: 12, padding: 20 }}>
              <div style={{ color: '#8899a6', fontSize: 13, marginBottom: 8 }}>Avg Score</div>
              <div style={{ fontSize: 32, fontWeight: 700, color: analytics.avg_score >= 70 ? '#22c55e' : '#f59e0b' }}>
                {Math.round(analytics.avg_score)}
              </div>
            </div>
            <div style={{ background: 'rgba(255,255,255,0.05)', borderRadius: 12, padding: 20 }}>
              <div style={{ color: '#8899a6', fontSize: 13, marginBottom: 8 }}>Failed</div>
              <div style={{ fontSize: 32, fontWeight: 700, color: '#ef4444' }}>{analytics.failed_jobs}</div>
            </div>
          </div>

          {/* Score Distribution Bar */}
          <div style={{ background: 'rgba(255,255,255,0.05)', borderRadius: 12, padding: 20, marginTop: 16 }}>
            <div style={{ color: '#8899a6', fontSize: 13, marginBottom: 12 }}>Score Distribution (Last 30 Days)</div>
            <div style={{ display: 'flex', height: 24, borderRadius: 6, overflow: 'hidden', gap: 2 }}>
              <div style={{ width: `${analytics.score_distribution.low || 0}%`, background: '#22c55e', minWidth: 2 }} />
              <div style={{ width: `${analytics.score_distribution.medium || 0}%`, background: '#f59e0b', minWidth: 2 }} />
              <div style={{ width: `${analytics.score_distribution.high || 0}%`, background: '#ef4444', minWidth: 2 }} />
            </div>
            <div style={{ display: 'flex', gap: 16, marginTop: 8, fontSize: 12, color: '#8899a6' }}>
              <span>✓ Low Risk ({analytics.score_distribution.low || 0})</span>
              <span>⚠ Medium ({analytics.score_distribution.medium || 0})</span>
              <span>✗ High Risk ({analytics.score_distribution.high || 0})</span>
            </div>
          </div>

          {/* Job Trend Chart (Simple) */}
          {analytics.job_trend && analytics.job_trend.length > 0 && (
            <div style={{ background: 'rgba(255,255,255,0.05)', borderRadius: 12, padding: 20, marginTop: 16 }}>
              <div style={{ color: '#8899a6', fontSize: 13, marginBottom: 12 }}>Job Activity Trend</div>
              <div style={{ display: 'flex', alignItems: 'flex-end', gap: 8, height: 80 }}>
                {analytics.job_trend.slice(-14).map((item: any, i: number) => (
                  <div key={i} style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 4 }}>
                    <div style={{
                      width: '100%',
                      background: '#533afd',
                      borderRadius: 4,
                      height: `${Math.min(100, (item.score / 100) * 80)}px`,
                      minWidth: 8
                    }} />
                    <div style={{ fontSize: 10, color: '#8899a6' }}>
                      {item.date.slice(5)}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {/* Upload Section */}
      <div style={{ background: 'rgba(255,255,255,0.05)', borderRadius: 12, padding: 24, marginBottom: 24 }}>
        <h3 style={{ marginBottom: 16 }}>Upload Contract (.docx / .pdf)</h3>
        <input
          type="file"
          accept=".docx,.pdf"
          onChange={handleUpload}
          disabled={uploading}
          multiple
          style={{ marginBottom: 16, display: 'block', color: '#8899a6' }}
        />
        {uploading && <p style={{ color: '#533afd' }}>Processing...</p>}
        <p style={{ color: '#8899a6', fontSize: 14 }}>
          Max 10MB per file. Support multiple files for batch analysis.
        </p>
      </div>

      {/* Text Analysis */}
      <div style={{ background: 'rgba(255,255,255,0.05)', borderRadius: 12, padding: 24, marginBottom: 24 }}>
        <h3 style={{ marginBottom: 16 }}>Analyze Sample Text</h3>
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
          {analyzing ? 'Analyzing...' : 'Analyze Sample Employment Contract'}
        </button>
      </div>

      {/* Result Section */}
      {result && (
        <div style={{ background: 'rgba(255,255,255,0.05)', borderRadius: 12, padding: 24, marginBottom: 24 }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
            <h3 style={{ margin: 0 }}>Analysis Result</h3>
            <div style={{ display: 'flex', gap: 8 }}>
              <button onClick={() => downloadExport('text')} style={{ padding: '8px 16px', background: 'rgba(255,255,255,0.1)', border: '1px solid rgba(255,255,255,0.2)', borderRadius: 6, color: '#e2e8f0', cursor: 'pointer', fontSize: 13 }}>
                📄 Text
              </button>
              <button onClick={() => downloadExport('html')} style={{ padding: '8px 16px', background: 'rgba(255,255,255,0.1)', border: '1px solid rgba(255,255,255,0.2)', borderRadius: 6, color: '#e2e8f0', cursor: 'pointer', fontSize: 13 }}>
                🌐 HTML
              </button>
              <button onClick={() => downloadExport('json')} style={{ padding: '8px 16px', background: 'rgba(255,255,255,0.1)', border: '1px solid rgba(255,255,255,0.2)', borderRadius: 6, color: '#e2e8f0', cursor: 'pointer', fontSize: 13 }}>
                📋 JSON
              </button>
            </div>
          </div>
          
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
          {result.suggestions && result.suggestions.length > 0 && ('suggestions' in result) && (
            <div>
              <h4 style={{ marginBottom: 12, color: '#f8fafc' }}>Suggestions</h4>
              <ul style={{ color: '#e2e8f0', paddingLeft: 20, lineHeight: 1.8 }}>
                {result.suggestions.map((s: string, i: number) => <li key={i}>{s}</li>)}
              </ul>
            </div>
          )}

          {/* Scenario Analysis */}
          {result.scenario && (
            <div style={{ marginTop: 24 }}>
              <h4 style={{ marginBottom: 12, color: '#f8fafc' }}>
                {result.scenario.type_name || 'Contract Type Analysis'}
              </h4>
              {result.scenario.issues && result.scenario.issues.length > 0 && (
                <div style={{ marginBottom: 16 }}>
                  <div style={{ color: '#8899a6', fontSize: 12, marginBottom: 8 }}>Issues Found: {result.scenario.issues.filter((i: any) => !i.found).length}</div>
                  {result.scenario.issues.map((issue: any, i: number) => (
                    <div key={i} style={{
                      padding: 12,
                      marginBottom: 8,
                      borderRadius: 8,
                      background: issue.found ? 'rgba(34,197,94,0.05)' : 'rgba(239,68,68,0.05)',
                      borderLeft: `4px solid ${issue.found ? '#22c55e' : '#ef4444'}`,
                    }}>
                      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4 }}>
                        <span style={{ color: '#e2e8f0', fontSize: 13 }}>{issue.title}</span>
                        <span style={{
                          fontSize: 10,
                          padding: '2px 6px',
                          borderRadius: 4,
                          background: issue.found ? 'rgba(34,197,94,0.2)' : 'rgba(239,68,68,0.2)',
                          color: issue.found ? '#22c55e' : '#ef4444',
                        }}>
                          {issue.found ? 'PRESENT' : 'MISSING'}
                        </span>
                      </div>
                      {!issue.found && issue.recommend && (
                        <div style={{ color: '#8899a6', fontSize: 12 }}>{issue.recommend}</div>
                      )}
                    </div>
                  ))}
                </div>
              )}
              {result.scenario.requirements && result.scenario.requirements.length > 0 && (
                <div>
                  <div style={{ color: '#8899a6', fontSize: 12, marginBottom: 8 }}>
                    Required Clauses: {result.scenario.requirements.filter((r: any) => r.found).length}/{result.scenario.requirements.length}
                  </div>
                  <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
                    {result.scenario.requirements.map((req: any, i: number) => (
                      <span key={i} style={{
                        padding: '4px 10px',
                        borderRadius: 12,
                        fontSize: 11,
                        background: req.found ? 'rgba(34,197,94,0.15)' : 'rgba(239,68,68,0.15)',
                        color: req.found ? '#22c55e' : '#ef4444',
                        border: '1px solid rgba(255,255,255,0.1)',
                      }}>
                        {req.found ? '✓' : '○'} {req.name}
                      </span>
                    ))}
                  </div>
                </div>
              )}
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
