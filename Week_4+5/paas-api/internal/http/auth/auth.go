package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"` //seconds
}

type Config struct {
	AuthUser    string
	AuthPass    string
	JWTSecret   []byte
	JWTIssuer   string
	JWTAudience string
	JWTTL       time.Duration
}

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func RegisterAuthRoutes(e *echo.Echo, cfg Config) {
	e.POST("/auth/login", func(c echo.Context) error {
		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid json")
		}
		if req.Username != cfg.AuthUser || req.Password != cfg.AuthPass {
			// do not leak which part is wrong
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
		}

		now := time.Now()
		exp := now.Add(cfg.JWTTL)

		// MVP role: admin if auth user matches env user
		claims := Claims{
			Role: "admin",
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   req.Username,
				Issuer:    cfg.JWTIssuer,
				Audience:  jwt.ClaimStrings{cfg.JWTAudience},
				IssuedAt:  jwt.NewNumericDate(now),
				ExpiresAt: jwt.NewNumericDate(exp),
			},
		}

		tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := tok.SignedString(cfg.JWTSecret)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to sign token")
		}

		return c.JSON(http.StatusOK, LoginResponse{
			AccessToken: signed,
			TokenType:   "Bearer",
			ExpiresIn:   int64(cfg.JWTTL.Seconds()),
		})
	})

	e.GET("/auth/me", func(c echo.Context) error {
		claims := GetClaims(c)
		return c.JSON(http.StatusOK, map[string]any{
			"sub":  claims.Subject,
			"role": claims.Role,
			"exp":  claims.ExpiresAt.Time.Unix(),
			"iss":  claims.Issuer,
			"aud":  claims.Audience,
		})
	}, RequireJWT(cfg))
}

// --- Middleware helpers ---

const claimsKey = "jwt_claims"

func RequireJWT(cfg Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authz := c.Request().Header.Get("Authorization")
			const prefix = "Bearer "
			if len(authz) <= len(prefix) || authz[:len(prefix)] != prefix {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing bearer token")
			}
			raw := authz[len(prefix):]

			token, err := jwt.ParseWithClaims(raw, &Claims{}, func(t *jwt.Token) (any, error) {
				// enforce HS256
				if t.Method != jwt.SigningMethodHS256 {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid token method")
				}
				return cfg.JWTSecret, nil
			}, jwt.WithIssuer(cfg.JWTIssuer), jwt.WithAudience(cfg.JWTAudience))
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}
			claims, ok := token.Claims.(*Claims)
			if !ok || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			c.Set(claimsKey, claims)
			return next(c)
		}
	}
}

func RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := GetClaims(c)
			if claims == nil || claims.Role != role {
				return echo.NewHTTPError(http.StatusForbidden, "forbidden")
			}
			return next(c)
		}
	}
}

func GetClaims(c echo.Context) *Claims {
	v := c.Get(claimsKey)
	if v == nil {
		return nil
	}
	claims, _ := v.(*Claims)
	return claims
}
