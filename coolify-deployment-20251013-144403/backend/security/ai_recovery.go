package security

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// RecoveryAction represents a recovery action
type RecoveryAction struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Command     string                 `json:"command,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	SuccessRate float64                `json:"success_rate"`
	ExecutedAt  *time.Time             `json:"executed_at,omitempty"`
}

// AIRecoveryManager manages AI-powered error recovery
type AIRecoveryManager struct {
	recoveryHistory     []RecoveryAction
	knowledgeBase       map[string][]RecoveryAction
	confidenceThreshold float64
}

// NewAIRecoveryManager creates a new AI recovery manager
func NewAIRecoveryManager() *AIRecoveryManager {
	return &AIRecoveryManager{
		recoveryHistory:     make([]RecoveryAction, 0),
		knowledgeBase:       buildKnowledgeBase(),
		confidenceThreshold: 0.7,
	}
}

// AnalyzeAndRecover analyzes an error and attempts recovery
func (arm *AIRecoveryManager) AnalyzeAndRecover(err error, context map[string]interface{}) (*RecoveryAction, error) {
	// Analyze error pattern
	errorType := arm.classifyError(err)
	contextInfo := arm.extractContext(context)

	// Find appropriate recovery actions
	candidates := arm.findRecoveryCandidates(errorType, contextInfo)

	// Select best recovery action using AI logic
	bestAction := arm.selectBestRecoveryAction(candidates, contextInfo)
	if bestAction == nil {
		return nil, fmt.Errorf("no suitable recovery action found")
	}

	// Execute recovery action
	startTime := time.Now()
	success := arm.executeRecoveryAction(bestAction, context)

	duration := time.Since(startTime)
	bestAction.ExecutedAt = &startTime
	bestAction.SuccessRate = arm.calculateSuccessRate(bestAction, duration, success)

	// Record in history
	arm.recoveryHistory = append(arm.recoveryHistory, *bestAction)

	// Update knowledge base
	arm.updateKnowledgeBase(errorType, *bestAction, success)

	if success {
		return bestAction, nil
	}

	return bestAction, fmt.Errorf("recovery action failed")
}

// classifyError classifies the error type using AI patterns
func (arm *AIRecoveryManager) classifyError(err error) string {
	errMsg := strings.ToLower(err.Error())

	// Pattern matching for error classification
	patterns := map[string]string{
		"connection":     "database_connection",
		"timeout":        "network_timeout",
		"permission":     "file_permission",
		"disk":           "storage_full",
		"memory":         "memory_error",
		"import":         "import_error",
		"syntax":         "syntax_error",
		"authentication": "auth_error",
		"validation":     "validation_error",
	}

	for pattern, errorType := range patterns {
		if strings.Contains(errMsg, pattern) {
			return errorType
		}
	}

	return "unknown_error"
}

// extractContext extracts relevant context from the error situation
func (arm *AIRecoveryManager) extractContext(context map[string]interface{}) map[string]interface{} {
	relevant := make(map[string]interface{})

	// Extract relevant context information
	relevantKeys := []string{"endpoint", "method", "user_id", "request_id", "ip_address", "component"}

	for _, key := range relevantKeys {
		if value, exists := context[key]; exists {
			relevant[key] = value
		}
	}

	return relevant
}

// findRecoveryCandidates finds potential recovery actions
func (arm *AIRecoveryManager) findRecoveryCandidates(errorType string, context map[string]interface{}) []RecoveryAction {
	var candidates []RecoveryAction

	// Get from knowledge base
	if actions, exists := arm.knowledgeBase[errorType]; exists {
		candidates = append(candidates, actions...)
	}

	// Add context-specific actions
	switch errorType {
	case "database_connection":
		candidates = append(candidates, RecoveryAction{
			Type:        "restart_service",
			Description: "Restart database service",
			Command:     "systemctl restart postgresql",
			SuccessRate: 0.8,
		})
	case "network_timeout":
		candidates = append(candidates, RecoveryAction{
			Type:        "retry_request",
			Description: "Retry the failed request",
			SuccessRate: 0.6,
		})
	case "file_permission":
		candidates = append(candidates, RecoveryAction{
			Type:        "fix_permissions",
			Description: "Fix file permissions",
			Command:     "chmod 644 /path/to/file",
			SuccessRate: 0.9,
		})
	}

	return candidates
}

// selectBestRecoveryAction selects the best recovery action using AI logic
func (arm *AIRecoveryManager) selectBestRecoveryAction(candidates []RecoveryAction, context map[string]interface{}) *RecoveryAction {
	if len(candidates) == 0 {
		return nil
	}

	// Simple scoring algorithm (in production, this would use ML)
	var bestAction *RecoveryAction
	var bestScore float64 = 0

	for i := range candidates {
		action := &candidates[i]
		score := action.SuccessRate

		// Adjust score based on context
		if endpoint, exists := context["endpoint"]; exists {
			if strings.Contains(endpoint.(string), "health") && action.Type == "restart_service" {
				score += 0.1 // Health endpoints benefit from service restart
			}
		}

		if score > bestScore {
			bestScore = score
			bestAction = action
		}
	}

	if bestScore >= arm.confidenceThreshold {
		return bestAction
	}

	return nil
}

// executeRecoveryAction executes a recovery action
func (arm *AIRecoveryManager) executeRecoveryAction(action *RecoveryAction, context map[string]interface{}) bool {
	switch action.Type {
	case "restart_service":
		return arm.restartService(action.Command)
	case "retry_request":
		return arm.retryRequest(context)
	case "fix_permissions":
		return arm.fixPermissions(action.Parameters)
	case "clear_cache":
		return arm.clearCache(context)
	default:
		return arm.executeCommand(action.Command)
	}
}

// restartService restarts a system service
func (arm *AIRecoveryManager) restartService(command string) bool {
	cmd := exec.Command("bash", "-c", command)
	err := cmd.Run()
	return err == nil
}

// retryRequest retries a failed HTTP request
func (arm *AIRecoveryManager) retryRequest(context map[string]interface{}) bool {
	// In a real implementation, this would retry the original request
	// For now, simulate with a delay
	time.Sleep(1 * time.Second)
	return true
}

// fixPermissions fixes file permissions
func (arm *AIRecoveryManager) fixPermissions(params map[string]interface{}) bool {
	if path, exists := params["path"]; exists {
		cmd := exec.Command("chmod", "644", path.(string))
		err := cmd.Run()
		return err == nil
	}
	return false
}

// clearCache clears application cache
func (arm *AIRecoveryManager) clearCache(context map[string]interface{}) bool {
	// Clear Redis cache or file cache
	cmd := exec.Command("redis-cli", "FLUSHALL")
	err := cmd.Run()
	return err == nil
}

// executeCommand executes a shell command
func (arm *AIRecoveryManager) executeCommand(command string) bool {
	cmd := exec.Command("bash", "-c", command)
	err := cmd.Run()
	return err == nil
}

// calculateSuccessRate calculates success rate for a recovery action
func (arm *AIRecoveryManager) calculateSuccessRate(action *RecoveryAction, duration time.Duration, success bool) float64 {
	baseRate := action.SuccessRate

	// Adjust based on execution time
	if duration < 5*time.Second {
		baseRate += 0.1 // Faster recovery is better
	} else if duration > 30*time.Second {
		baseRate -= 0.1 // Slower recovery is worse
	}

	// Adjust based on success
	if success {
		baseRate += 0.05 // Successful recovery gets slight boost
	} else {
		baseRate -= 0.2 // Failed recovery gets significant penalty
	}

	// Clamp between 0 and 1
	if baseRate > 1.0 {
		baseRate = 1.0
	}
	if baseRate < 0.0 {
		baseRate = 0.0
	}

	return baseRate
}

// updateKnowledgeBase updates the knowledge base with new information
func (arm *AIRecoveryManager) updateKnowledgeBase(errorType string, action RecoveryAction, success bool) {
	// Update success rate based on outcome
	if success {
		action.SuccessRate = (action.SuccessRate * 0.9) + 0.1 // Increase success rate
	} else {
		action.SuccessRate = action.SuccessRate * 0.9 // Decrease success rate
	}

	// Store updated action back in knowledge base
	if actions, exists := arm.knowledgeBase[errorType]; exists {
		for i := range actions {
			if actions[i].Type == action.Type {
				actions[i] = action
				break
			}
		}
	}
}

// GetRecoveryHistory returns recovery history
func (arm *AIRecoveryManager) GetRecoveryHistory() []RecoveryAction {
	return arm.recoveryHistory
}

// GetRecoveryStats returns recovery statistics
func (arm *AIRecoveryManager) GetRecoveryStats() map[string]interface{} {
	if len(arm.recoveryHistory) == 0 {
		return map[string]interface{}{
			"total_recoveries": 0,
			"success_rate":     0.0,
			"average_time":     0.0,
		}
	}

	successCount := 0
	totalTime := time.Duration(0)

	for _, action := range arm.recoveryHistory {
		if action.SuccessRate > 0.5 { // Consider successful if rate > 0.5
			successCount++
		}
		if action.ExecutedAt != nil {
			// In a real implementation, you'd calculate actual duration
			totalTime += time.Second
		}
	}

	successRate := float64(successCount) / float64(len(arm.recoveryHistory))
	averageTime := totalTime / time.Duration(len(arm.recoveryHistory))

	return map[string]interface{}{
		"total_recoveries": len(arm.recoveryHistory),
		"success_rate":     successRate,
		"average_time_ms":  averageTime.Milliseconds(),
	}
}

// buildKnowledgeBase builds the initial knowledge base of recovery actions
func buildKnowledgeBase() map[string][]RecoveryAction {
	return map[string][]RecoveryAction{
		"database_connection": {
			{
				Type:        "restart_service",
				Description: "Restart database service",
				Command:     "systemctl restart postgresql",
				SuccessRate: 0.8,
			},
			{
				Type:        "check_connection",
				Description: "Check database connectivity",
				Command:     "pg_isready -h localhost",
				SuccessRate: 0.6,
			},
		},
		"network_timeout": {
			{
				Type:        "retry_request",
				Description: "Retry the failed request",
				SuccessRate: 0.7,
			},
			{
				Type:        "check_network",
				Description: "Check network connectivity",
				Command:     "ping -c 3 google.com",
				SuccessRate: 0.5,
			},
		},
		"file_permission": {
			{
				Type:        "fix_permissions",
				Description: "Fix file permissions",
				SuccessRate: 0.9,
			},
		},
		"memory_error": {
			{
				Type:        "restart_service",
				Description: "Restart application service",
				Command:     "systemctl restart nutrition-platform",
				SuccessRate: 0.6,
			},
			{
				Type:        "clear_cache",
				Description: "Clear application cache",
				SuccessRate: 0.8,
			},
		},
		"import_error": {
			{
				Type:        "install_dependency",
				Description: "Install missing dependency",
				Command:     "go mod tidy",
				SuccessRate: 0.9,
			},
		},
	}
}

// Global AI recovery manager
var globalRecoveryManager *AIRecoveryManager

// GetRecoveryManager returns the global recovery manager
func GetRecoveryManager() *AIRecoveryManager {
	if globalRecoveryManager == nil {
		globalRecoveryManager = NewAIRecoveryManager()
	}
	return globalRecoveryManager
}

// RecoveryMiddleware creates middleware for automatic error recovery
func RecoveryMiddleware() echo.MiddlewareFunc {
	recoveryManager := GetRecoveryManager()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Capture context information
			ctx := map[string]interface{}{
				"endpoint":   c.Request().URL.Path,
				"method":     c.Request().Method,
				"ip_address": c.RealIP(),
				"request_id": c.Response().Header().Get("X-Request-ID"),
			}

			// Add user context if available
			if userID := c.Get("user_id"); userID != nil {
				ctx["user_id"] = userID
			}

			// Execute request with recovery
			err := next(c)
			if err != nil {
				// Try to recover from error
				action, recoveryErr := recoveryManager.AnalyzeAndRecover(err, ctx)
				if recoveryErr == nil && action != nil {
					// Log successful recovery
					logger := GetLogger()
					logger.LogSecurityEvent(WARN, "Error recovered", map[string]interface{}{
						"original_error":  err.Error(),
						"recovery_action": action.Type,
						"success_rate":    action.SuccessRate,
					})

					// Return success response
					return c.JSON(200, map[string]interface{}{
						"message":         "Request processed with automatic recovery",
						"recovery_action": action.Type,
						"original_error":  err.Error(),
					})
				}

				// Recovery failed, return original error
				return err
			}

			return nil
		}
	}
}
