package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RBACMiddleware handles role-based access control for customer operations
type RBACMiddleware struct{}

// NewRBACMiddleware creates a new RBAC middleware
func NewRBACMiddleware() *RBACMiddleware {
	return &RBACMiddleware{}
}

// RequireRole middleware checks if user has one of the required roles
func (m *RBACMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Unauthorized: No user role found",
			})
			c.Abort()
			return
		}

		userRoleStr, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Invalid user role format",
			})
			c.Abort()
			return
		}

		// Check if user has one of the required roles
		for _, role := range roles {
			if strings.EqualFold(userRoleStr, role) {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Forbidden: Insufficient role permissions",
		})
		c.Abort()
	}
}

// RequirePermission middleware checks if user has a specific permission
func (m *RBACMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userPermissions, exists := c.Get("user_permissions")
		if !exists {
			// Fall back to role-based check
			userRole, roleExists := c.Get("user_role")
			if !roleExists {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Unauthorized: No permissions found",
				})
				c.Abort()
				return
			}

			// Super admin and manager bypass
			if role, ok := userRole.(string); ok {
				if strings.EqualFold(role, "SUPER_ADMIN") || strings.EqualFold(role, "admin") || strings.EqualFold(role, "MANAGER") {
					c.Next()
					return
				}
			}

			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Forbidden: Missing required permission: " + permission,
			})
			c.Abort()
			return
		}

		hasPermission := false

		switch perms := userPermissions.(type) {
		case []string:
			for _, p := range perms {
				if p == permission {
					hasPermission = true
					break
				}
			}
		case string:
			permissions := strings.Split(perms, ",")
			for _, p := range permissions {
				if strings.TrimSpace(p) == permission {
					hasPermission = true
					break
				}
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Forbidden: Missing required permission: " + permission,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission middleware checks if user has any of the specified permissions
func (m *RBACMiddleware) RequireAnyPermission(permissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userPermissions, exists := c.Get("user_permissions")
		if !exists {
			userRole, roleExists := c.Get("user_role")
			if !roleExists {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Unauthorized: No permissions found",
				})
				c.Abort()
				return
			}

			// Super admin and manager bypass
			if role, ok := userRole.(string); ok {
				if strings.EqualFold(role, "SUPER_ADMIN") || strings.EqualFold(role, "admin") || strings.EqualFold(role, "MANAGER") {
					c.Next()
					return
				}
			}

			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Forbidden: Missing one of required permissions",
			})
			c.Abort()
			return
		}

		hasAnyPermission := false

		switch perms := userPermissions.(type) {
		case []string:
			for _, userPerm := range perms {
				for _, requiredPerm := range permissions {
					if userPerm == requiredPerm {
						hasAnyPermission = true
						break
					}
				}
				if hasAnyPermission {
					break
				}
			}
		case string:
			userPerms := strings.Split(perms, ",")
			for _, userPerm := range userPerms {
				for _, requiredPerm := range permissions {
					if strings.TrimSpace(userPerm) == requiredPerm {
						hasAnyPermission = true
						break
					}
				}
				if hasAnyPermission {
					break
				}
			}
		}

		if !hasAnyPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Forbidden: Missing one of required permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CustomerAdminMiddleware checks if user has customer admin permissions
func CustomerAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Unauthorized: No user role found",
			})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Invalid user role format",
			})
			c.Abort()
			return
		}

		// Allow customer management roles
		allowedRoles := []string{"admin", "superadmin", "SUPER_ADMIN", "MANAGER", "STAFF_ORDERS", "SALES_AGENT"}
		isAllowed := false
		for _, allowed := range allowedRoles {
			if strings.EqualFold(role, allowed) {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Forbidden: Customer admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserIDFromContext retrieves the user ID from the Gin context
func GetUserIDFromContext(c *gin.Context) uuid.UUID {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil
	}

	switch v := userID.(type) {
	case uuid.UUID:
		return v
	case string:
		id, err := uuid.Parse(v)
		if err != nil {
			return uuid.Nil
		}
		return id
	default:
		return uuid.Nil
	}
}

// GetUserRoleFromContext retrieves the user role from the Gin context
func GetUserRoleFromContext(c *gin.Context) string {
	userRole, exists := c.Get("user_role")
	if !exists {
		return ""
	}

	role, ok := userRole.(string)
	if !ok {
		return ""
	}

	return role
}
