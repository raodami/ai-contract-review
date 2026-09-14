package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ai-contract-review/internal/auth"
	"ai-contract-review/internal/nlp"
	"ai-contract-review/internal/store"
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
			allowed, remaining, _ := s.CheckQuota("", "")
			c.JSON(http.StatusOK, gin.H{
				"allowed":       allowed,
				"remaining_min": remaining,
				"free_quota":    store.FreeQuota,
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

			// Mark as pro
			if err := s.SetUserPro(userID, true); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Subscribed to " + req.Plan})
		})
	}

	// Contract analysis routes
	contract := r.Group("/api/contract")
	{
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

			// Check quota
			allowed, _, _ := s.CheckQuota(userID, "")
			if !allowed {
				c.JSON(http.StatusPaymentRequired, gin.H{"error": "Quota exceeded", "checkout_url": "/api/subscribe"})
				return
			}

			// Analyze contract
			clauses := nlp.AnalyzeKeywords(req.Text)
			terms := nlp.ExtractKeyTerms(req.Text)

			// Calculate risk score
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

			// Summarize
			summary := buildSummary(clauses, terms)

			// Increment usage
			s.IncrementUsage(userID, 0)

			c.JSON(http.StatusOK, gin.H{
				"summary":           summary,
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
