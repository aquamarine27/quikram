package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"quikram-service/internal/config"
	"quikram-service/internal/service"
)

const (
	cookieName   = "refresh_token"
	cookiePath   = "/api/v1/auth"
)

type AuthHandler struct {
	authService *service.AuthService
	cfg         *config.Config
}

func NewAuthHandler(authService *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{authService: authService, cfg: cfg}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	result, err := h.authService.Register(req.Email, req.Password, req.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	h.setRefreshCookie(c, result.RefreshToken)
	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	result, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	h.setRefreshCookie(c, result.RefreshToken)
	return c.JSON(result)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies(cookieName)
	if refreshToken == "" {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.BodyParser(&req); err != nil || req.RefreshToken == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "refresh token required"})
		}
		refreshToken = req.RefreshToken
	}

	result, err := h.authService.Refresh(refreshToken)
	if err != nil {
		h.clearRefreshCookie(c)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	h.setRefreshCookie(c, result.RefreshToken)
	return c.JSON(result)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies(cookieName)
	if refreshToken == "" {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.BodyParser(&req); err == nil && req.RefreshToken != "" {
			refreshToken = req.RefreshToken
		}
	}

	if refreshToken != "" {
		h.authService.Logout(refreshToken)
	}

	h.clearRefreshCookie(c)
	return c.JSON(fiber.Map{"message": "logged out successfully"})
}

func (h *AuthHandler) setRefreshCookie(c *fiber.Ctx, token string) {
	secure := h.cfg.Env == "production"
	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     cookiePath,
		Expires:  time.Now().Add(h.cfg.JWTRefreshExpire),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Strict",
	})
}

func (h *AuthHandler) clearRefreshCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     cookiePath,
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		SameSite: "Strict",
	})
}
