package store

import (
	"database/sql"
	"time"
)

// AnalyticsData holds statistical data for dashboard
type AnalyticsData struct {
	TotalJobs       int              `json:"total_jobs"`
	CompletedJobs   int              `json:"completed_jobs"`
	FailedJobs      int              `json:"failed_jobs"`
	AvgScore        float64          `json:"avg_score"`
	JobTrend        []timeScore      `json:"job_trend"`
	ScoreDistribution map[string]int `json:"score_distribution"`
}

type timeScore struct {
	Date string  `json:"date"`
	Score float64 `json:"score"`
}

// GetAnalytics returns usage statistics for the given user
func (s *Store) GetAnalytics(userID string, days int) (*AnalyticsData, error) {
	data := &AnalyticsData{}

	// Total jobs
	err := s.db.QueryRow("SELECT COUNT(*) FROM jobs WHERE user_id = ?", userID).Scan(&data.TotalJobs)
	if err != nil {
		return nil, err
	}

	// Completed jobs
	err = s.db.QueryRow("SELECT COUNT(*) FROM jobs WHERE user_id = ? AND status = ?", userID, JobCompleted).Scan(&data.CompletedJobs)
	if err != nil {
		return nil, err
	}

	// Failed jobs
	err = s.db.QueryRow("SELECT COUNT(*) FROM jobs WHERE user_id = ? AND status = ?", userID, JobFailed).Scan(&data.FailedJobs)
	if err != nil {
		return nil, err
	}

	// Total usage (sum of estimated minutes from results)
	var totalUsage sql.NullFloat64
	err = s.db.QueryRow(`
		SELECT AVG(CAST(json_extract(result, '$.score') AS REAL)) 
		FROM jobs WHERE user_id = ? AND status = ?
	`, userID, JobCompleted).Scan(&totalUsage)
	if err != nil {
		return nil, err
	}
	if totalUsage.Valid {
		data.AvgScore = totalUsage.Float64
	}

	// Job trend (last N days)
	startDate := time.Now().AddDate(0, 0, -days)
	rows, err := s.db.Query(`
		SELECT date(created_at / 1000, 'unixepoch') as day,
			   COUNT(*) as cnt,
			   AVG(CAST(json_extract(result, '$.score') AS REAL)) as avg_score
		FROM jobs 
		WHERE user_id = ? AND created_at >= ?
		GROUP BY day
		ORDER BY day DESC
		LIMIT ?
	`, userID, startDate.Unix()*1000, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ts timeScore
		if err := rows.Scan(&ts.Date, &ts.Score); err != nil {
			continue
		}
		data.JobTrend = append([]timeScore{ts}, data.JobTrend...)
	}

	// Score distribution
	var low, medium, high int
	row := s.db.QueryRow(`
		SELECT 
			COUNT(CASE WHEN CAST(json_extract(result, '$.score') AS REAL) >= 80 THEN 1 END) as low,
			COUNT(CASE WHEN CAST(json_extract(result, '$.score') AS REAL) BETWEEN 50 AND 79 THEN 1 END) as medium,
			COUNT(CASE WHEN CAST(json_extract(result, '$.score') AS REAL) < 50 THEN 1 END) as high
		FROM jobs 
		WHERE user_id = ? AND status = ?
	`, userID, JobCompleted)
	if err := row.Scan(&low, &medium, &high); err != nil {
		return nil, err
	}
	data.ScoreDistribution = map[string]int{
		"low":      low,
		"medium":   medium,
		"high":     high,
	}

	return data, nil
}
