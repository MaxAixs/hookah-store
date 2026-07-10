package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/anomalyco/hookah-store/user-service/internal/errs"
	jwtpkg "github.com/anomalyco/hookah-store/user-service/pkg/jwt"
	"github.com/gin-gonic/gin"
)

const adminRole = "admin"

func RequireAuth(jwtCfg *jwtpkg.JwtConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := extractToken(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())

			return
		}

		_, err = jwtCfg.Validate(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errs.ErrInvalidToken.Error())

			return
		}

		ctx.Next()
	}
}

func RequireRole(jwtCfg *jwtpkg.JwtConfig, roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := extractToken(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())

			return
		}

		claims, err := jwtCfg.Validate(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errs.ErrInvalidToken.Error())

			return
		}

		for _, role := range roles {
			if claims.Role == role {
				ctx.Next()

				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, errs.ErrAccessDenied.Error())
	}
}

func RequireAdminRole(jwtCfg *jwtpkg.JwtConfig) gin.HandlerFunc {
	return RequireRole(jwtCfg, adminRole)
}

func extractToken(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", fmt.Errorf("invalid authorization header format")
	}

	if parts[1] == "" {
		return "", fmt.Errorf("empty token")
	}

	return parts[1], nil
}
