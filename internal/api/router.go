package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"encoding/json"
	"ai-contract-review/internal/auth"
	"ai-contract-review/internal/llm"
	"ai-contract-review/internal/nlp"
	"ai-contract-review/internal/parser"
	"ai-contract-review/internal/store"
	"ai-contract-review/internal/payment"
	"ai-contract-review/internal/export"
)

// RegisterRequest represents a register request
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserInfo represents current user info
type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	IsPro     bool   `json:"is_pro"`
	UsageMin  int    `json:"usage_minutes"`
	FreeQuota int    `json:"free_quota"`
}

// SetupRoutes sets up auth routes
func SetupRoutes(r *gin.Engine, s *store.Store) {
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", func(c *gin.Context) {
			var req RegisterRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Check if user exists
			existing, _ := s.GetUserByEmail(req.Email)
			if existing != nil {
				c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
				return
			}

			// Hash password
			hashedPassword, err := auth.HashPassword(req.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
				return
			}

			// Create user
			userID := uuid.New().String()
			if err := s.CreateUser(userID, req.Email, hashedPassword); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}

			// Generate JWT
			token, err := auth.GenerateToken(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}

			c.JSON(http.StatusCreated, gin.H{
				"token": token,
				"user": UserInfo{
					ID:        userID,
					Email:     req.Email,
					FreeQuota: store.FreeQuota,
				},
			})
		})

		authGroup.POST("/login", func(c *gin.Context) {
			var req LoginRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			user, err := s.GetUserByEmail(req.Email)
			if err != nil || user == nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}

			if !auth.CheckPassword(req.Password, user.Password) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}

			token, err := auth.GenerateToken(user.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"token": token,
				"user": UserInfo{
					ID:        user.ID,
					Email:     user.Email,
					IsPro:     user.IsPro,
					UsageMin:  user.UsageMin,
					FreeQuota: store.FreeQuota,
				},
			})
		})

		authGroup.GET("/me", func(c *gin.Context) {
			tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

			userID, err := auth.ParseToken(tokenStr)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			user, err := s.GetUserByID(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
				return
			}

			allowed, remaining, _ := s.CheckQuota(userID, "")
			c.JSON(http.StatusOK, gin.H{
				"id":            user.ID,
				"email":         user.Email,
				"is_pro":        user.IsPro,
				"usage_minutes": user.UsageMin,
				"free_quota":    store.FreeQuota,
				"allowed":       allowed,
				"remaining_min": remaining,
			})
		})
	}

	userGroup := r.Group("/api/user")
	{
		userGroup.GET("/usage", func(c *gin.Context) {
			tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
			userID, err := auth.ParseToken(tokenStr)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}
			user, _ := s.GetUserByID(userID)
			if user == nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				return
			}
			allowed, remaining, limit := payment.CheckQuota(user.IsPro, user.UsageMin)
			c.JSON(http.StatusOK, gin.H{
				"id":            user.ID,
				"is_pro":        user.IsPro,
				"allowed":       allowed,
				"remaining_min": remaining,
				"usage_minutes": user.UsageMin,
				"quota_limit":   limit,
				"free_quota":    payment.FreeQuotaMinutes,
			})
		})

		userGroup.PUT("/subscribe", func(c *gin.Context) {
			tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
			userID, err := auth.ParseToken(tokenStr)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			var req struct {
				Plan string `json:"plan" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if err := s.SetUserPro(userID, true); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Subscribed to " + req.Plan})
		})
	}

	// Payment routes
	paymentGroup := r.Group("/api/payment")
	{
		// GET /api/payment/plans - list available plans
		paymentGroup.GET("/plans", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"plans": payment.GetPlans()})
		})

		// POST /api/payment/checkout - create checkout session
		paymentGroup.POST("/checkout", func(c *gin.Context) {
			tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
			userID, err := auth.ParseToken(tokenStr)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			var req struct {
				PlanID string `json:"plan_id" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			user, _ := s.GetUserByID(userID)
			email := ""
			if user != nil {
				email = user.Email
			}

			sessionURL, err := payment.CheckoutSession(payment.CheckoutSessionInput{
				UserID: userID,
				Email:  email,
				PlanID: req.PlanID,
				SuccessURL: "https://your-domain.com/dashboard?success=true",
				CancelURL:  "https://your-domain.com/pricing?canceled=true",
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"checkout_url": sessionURL,
				"plan_id":      req.PlanID,
			})
		})

		// POST /api/payment/webhook - Stripe webhook
		paymentGroup.POST("/webhook", func(c *gin.Context) {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
				return
			}

			event, err := payment.WebhookHandler(body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": event})
		})
	}

	// Contract analysis routes
	contract := r.Group("/api/contract")
	{
		// POST /api/contract/upload — upload file, return job_id
		contract.POST("/upload", func(c *gin.Context) {
			tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
			userID, err := auth.ParseToken(tokenStr)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			allowed, remaining, limit := payment.CheckQuota(false, 0)
			if !allowed {
				c.JSON(http.StatusPaymentRequired, gin.H{"error": "Quota exceeded", "checkout_url": "/api/subscribe"})
				return
			}

			form, err := c.MultipartForm()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			files := form.File["file"]
			if len(files) == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
				return
			}

			var results []map[string]any
			for _, file := range files {
				f, err := file.Open()
				if err != nil {
					results = append(results, map[string]any{"file": file.Filename, "error": "Failed to open"})
					continue
				}

				data, err := io.ReadAll(f)
				f.Close()
				if err != nil {
					results = append(results, map[string]any{"file": file.Filename, "error": "Failed to read"})
					continue
				}
				if len(data) > 10*1024*1024 {
					results = append(results, map[string]any{"file": file.Filename, "error": "File too large (max 10MB)"})
					continue
				}

				jobID := uuid.New().String()
				if err := s.CreateJob(jobID, userID, file.Filename, int64(len(data))); err != nil {
					results = append(results, map[string]any{"file": file.Filename, "job_id": jobID, "error": "Failed to create job"})
					continue
				}

				text, err := parser.ParseDOCX(data)
				if err != nil {
					if parser.IsPDF(data) {
						text, err = parser.ParsePDF(data)
					}
				}
				if err != nil {
					s.UpdateJobStatus(jobID, store.JobFailed, fmt.Sprintf("parse error: %v", err))
					results = append(results, map[string]any{"file": file.Filename, "job_id": jobID, "error": "Parse failed"})
					continue
				}
				if err := s.UpdateJobContent(jobID, text); err != nil {
					s.UpdateJobStatus(jobID, store.JobFailed, "save error")
					results = append(results, map[string]any{"file": file.Filename, "job_id": jobID, "error": "Save failed"})
					continue
				}
				s.UpdateJobStatus(jobID, store.JobProcessing, "")

				client := llm.NewDeepSeekClient()
				analysis, err := client.AnalyzeContract(text)
				if err != nil {
					s.UpdateJobStatus(jobID, store.JobFailed, fmt.Sprintf("ai error: %v", err))
					results = append(results, map[string]any{"file": file.Filename, "job_id": jobID, "error": "AI failed"})
					continue
				}
				resultJSON, _ := json.Marshal(analysis)
				s.UpdateJobStatus(jobID, store.JobCompleted, string(resultJSON))
				s.IncrementUserUsage(userID, 0)

				results = append(results, map[string]any{"file": file.Filename, "job_id": jobID, "status": "completed"})
			}

			newRemaining := remaining - len(files)
			if newRemaining < 0 {
				newRemaining = 0
			}

			c.JSON(http.StatusOK, gin.H{
				"jobs":         results,
				"total":        len(results),
				"remaining_min": newRemaining,
				"quota_limit":  limit,
			})
		})

		// GET /api/contract/jobs — list recent jobs
		contract.GET("/jobs", func(c *gin.Context) {
			tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
			userID, err := auth.ParseToken(tokenStr)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}
			jobs, err := s.ListJobs(userID, 20)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list jobs"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"jobs": jobs})
		})

		// GET /api/contract/jobs/:id — get job result
		contract.GET("/jobs/:id", func(c *gin.Context) {
			tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
			userID, _ := auth.ParseToken(tokenStr)
			jobID := c.Param("id")
			job, err := s.GetJob(jobID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
				return
			}
			if job.UserID != userID {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
				return
			}
			c.JSON(http.StatusOK, job)
		})

		// GET /api/contract/export/:id?format=text — export analysis report
		contract.GET("/export/:id", func(c *gin.Context) {
			tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
			userID, _ := auth.ParseToken(tokenStr)
			jobID := c.Param("id")
			format := c.Query("format")
			if format == "" {
				format = "text"
			}

			job, err := s.GetJob(jobID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
				return
			}
			if job.UserID != userID {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
				return
			}

			var result map[string]any
			if job.Result != "" {
				json.Unmarshal([]byte(job.Result), &result)
			}

			reportData := export.NewReportData(job.FileName, job.FileSize, result)
			opts := export.ExportOptions{Format: format}
			filename, content := export.GenerateExport(reportData, opts)

			c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

			switch format {
			case "html":
				c.Data(http.StatusOK, "text/html; charset=utf-8", content)
			case "json":
				c.Data(http.StatusOK, "application/json; charset=utf-8", content)
			case "csv":
				c.Data(http.StatusOK, "text/csv; charset=utf-8", content)
			case "docx":
				c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", content)
			default:
				c.Data(http.StatusOK, "text/plain; charset=utf-8", content)
			}
		})

	// GET /api/analytics — get usage statistics
	userGroup.GET("/analytics", func(c *gin.Context) {
		tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		userID, err := auth.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		analytics, err := s.GetAnalytics(userID, 30)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get analytics"})
			return
		}

		c.JSON(http.StatusOK, analytics)
	})

	// POST /api/contract/analyze — analyze text (for testing)
		contract.POST("/analyze", func(c *gin.Context) {
			tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
			userID, err := auth.ParseToken(tokenStr)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			var req struct {
				Text string `json:"text" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			allowed, _, _ := s.CheckQuota(userID, "")
			if !allowed {
				c.JSON(http.StatusPaymentRequired, gin.H{"error": "Quota exceeded", "checkout_url": "/api/subscribe"})
				return
			}

			client := llm.NewDeepSeekClient()
			if client.IsConfigured() {
				analysis, err := client.AnalyzeContract(req.Text)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "AI analysis failed: " + err.Error()})
					return
				}
				c.JSON(http.StatusOK, analysis)
				return
			}

			// Fallback to keyword analysis
			clauses := nlp.AnalyzeKeywords(req.Text)
			terms := nlp.ExtractKeyTerms(req.Text)
			score := 100
			for _, clause := range clauses {
				switch clause.Risk {
				case nlp.RiskCritical:
					score -= 25
				case nlp.RiskHigh:
					score -= 15
				case nlp.RiskMedium:
					score -= 5
				}
			}
			if score < 0 {
				score = 0
			}
			s.IncrementUserUsage(userID, 0)
			c.JSON(http.StatusOK, gin.H{
				"summary":           buildSummary(clauses, terms),
				"dangerous_clauses": clauses,
				"key_terms":         terms,
				"score":             score,
			})
		})
	}
}

func buildSummary(clauses []nlp.RiskClause, terms map[string]string) string {
	var sb strings.Builder
	sb.WriteString("Contract Analysis Summary:\n\n")

	if len(clauses) > 0 {
		sb.WriteString("## Risk Clauses Detected:\n")
		for _, c := range clauses {
			sb.WriteString(fmt.Sprintf("- [%s] %s\n", strings.ToUpper(c.Risk.String()), c.Term))
			if c.Suggestion != "" {
				sb.WriteString(fmt.Sprintf("  → %s\n", c.Suggestion))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Key Terms:\n")
	for k, v := range terms {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", k, v))
	}

	return sb.String()
}
