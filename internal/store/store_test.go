package store

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func newTestDB(t *testing.T) *Store {
	t.Helper()
	store, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test DB: %v", err)
	}
	return store
}

func TestCreateUser(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	id := uuid.New().String()
	err := s.CreateUser(id, "test@example.com", "hashed_password")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	user, err := s.GetUserByEmail("test@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected 'test@example.com', got '%s'", user.Email)
	}
	if user.Password != "hashed_password" {
		t.Errorf("Expected 'hashed_password', got '%s'", user.Password)
	}
}

func TestIncrementUserUsage(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	id := uuid.New().String()
	s.CreateUser(id, "test@example.com", "pass")
	s.IncrementUserUsage(id, 5)
	s.IncrementUserUsage(id, 10)

	user, _ := s.GetUserByID(id)
	if user.UsageMin != 15 {
		t.Errorf("Expected 15, got %d", user.UsageMin)
	}
}

func TestCheckQuota_FreeUser(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	id := uuid.New().String()
	s.CreateUser(id, "test@example.com", "pass")

	allowed, remaining, _ := s.CheckQuota(id, "")
	if !allowed {
		t.Error("Expected allowed=true for free user under quota")
	}
	if remaining != FreeQuota {
		t.Errorf("Expected remaining %d, got %d", FreeQuota, remaining)
	}

	s.IncrementUserUsage(id, FreeQuota)
	allowed, remaining, _ = s.CheckQuota(id, "")
	if allowed {
		t.Error("Expected allowed=false after exceeding quota")
	}
	if remaining != 0 {
		t.Errorf("Expected remaining 0, got %d", remaining)
	}
}

func TestCheckQuota_ProUser(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	id := uuid.New().String()
	s.CreateUser(id, "pro@example.com", "pass")
	s.SetUserPro(id, true)

	allowed, remaining, _ := s.CheckQuota(id, "")
	if !allowed {
		t.Error("Expected allowed=true for pro user")
	}
	if remaining != -1 {
		t.Errorf("Expected remaining -1 (unlimited), got %d", remaining)
	}
}

func TestAnonymousSession(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	cookieHash := "abc123"
	sess, err := s.GetOrCreateSession(cookieHash, "1.2.3.4")
	if err != nil {
		t.Fatalf("GetOrCreateSession failed: %v", err)
	}
	if sess.CookieHash != cookieHash {
		t.Errorf("Expected cookie hash '%s', got '%s'", cookieHash, sess.CookieHash)
	}

	// Same cookie should return same session
	sess2, err := s.GetOrCreateSession(cookieHash, "5.6.7.8")
	if err != nil {
		t.Fatalf("GetOrCreateSession failed: %v", err)
	}
	if sess2.ID != sess.ID {
		t.Error("Expected same session ID for same cookie hash")
	}

	// Increment usage
	s.IncrementSessionUsage(cookieHash, 10)
	sess3, _ := s.GetOrCreateSession(cookieHash, "")
	if sess3.UsageMin != 10 {
		t.Errorf("Expected 10, got %d", sess3.UsageMin)
	}
}

func TestFreeQuota_Constant(t *testing.T) {
	if FreeQuota != 30 {
		t.Errorf("Expected FreeQuota to be 30, got %d", FreeQuota)
	}
}

func TestCreateJob(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	id := uuid.New().String()
	err := s.CreateJob(id, "user123", "test.docx", 5000000)
	if err != nil {
		t.Fatalf("CreateJob failed: %v", err)
	}

	job, err := s.GetJob(id)
	if err != nil {
		t.Fatalf("GetJob failed: %v", err)
	}
	if job.Status != JobPending {
		t.Errorf("Expected status 'pending', got %v", job.Status)
	}
	if job.FileName != "test.docx" {
		t.Errorf("Expected file_name 'test.docx', got %v", job.FileName)
	}
}

func TestUpdateJobStatus(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	id := uuid.New().String()
	s.CreateJob(id, "user123", "test.docx", 1000000)

	err := s.UpdateJobStatus(id, JobCompleted, `{"summary":"ok"}`)
	if err != nil {
		t.Fatalf("UpdateJobStatus failed: %v", err)
	}

	job, _ := s.GetJob(id)
	if job.Status != JobCompleted {
		t.Errorf("Expected status 'completed', got %v", job.Status)
	}
	if job.Result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestListJobs(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	for i := 0; i < 3; i++ {
		id := uuid.New().String()
		s.CreateJob(id, "user123", fmt.Sprintf("test%d.docx", i), int64(1000*i))
	}

	jobs, err := s.ListJobs("user123", 10)
	if err != nil {
		t.Fatalf("ListJobs failed: %v", err)
	}
	if len(jobs) != 3 {
		t.Errorf("Expected 3 jobs, got %d", len(jobs))
	}
}

func TestGetAnalytics(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	userID := uuid.New().String()
	s.CreateUser(userID, "test@example.com", "pass")

	// Create some jobs with results
	for i := 0; i < 5; i++ {
		jobID := uuid.New().String()
		s.CreateJob(jobID, userID, "test.docx", 1000)
		score := 60 + i*10
		s.UpdateJobStatus(jobID, JobCompleted, `{"score":` + fmt.Sprintf("%d", score) + `,"summary":"test"}`)
	}

	analytics, err := s.GetAnalytics(userID, 30)
	if err != nil {
		t.Fatalf("GetAnalytics failed: %v", err)
	}
	if analytics.TotalJobs != 5 {
		t.Errorf("Expected 5 jobs, got %d", analytics.TotalJobs)
	}
	if analytics.CompletedJobs != 5 {
		t.Errorf("Expected 5 completed jobs, got %d", analytics.CompletedJobs)
	}
	if analytics.AvgScore < 60 || analytics.AvgScore > 100 {
		t.Errorf("Expected avg score around 80, got %f", analytics.AvgScore)
	}
}
