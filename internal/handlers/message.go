package handlers

import (
	"fmt"
	"time"

	z "github.com/Oudwins/zog"
	"github.com/Zigl3ur/api-portfolio/internal/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/client"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

type contactData struct {
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Subject string `json:"subject,omitempty"`
	Message string `json:"message,omitempty"`
}

var MessageLimiter = limiter.New(limiter.Config{
	Max:        2,
	Expiration: time.Hour * 24 * 7, // 1 week
	LimitReached: func(c fiber.Ctx) error {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Rate limit reached. Try again later.",
		})
	},
})

type MessageHandler struct {
	cfg *config.Config
}

func NewMessageHandler(cfg *config.Config) *MessageHandler {
	return &MessageHandler{cfg: cfg}
}

var messageSchema = z.Struct(z.Shape{
	"name":    z.String().Min(2, z.Message("name must be at least 2 characters long")).Max(15, z.Message("name must be at most 15 characters long")).Required(z.Message("name is required")),
	"email":   z.String().Email(z.Message("email is invalid")).Required(z.Message("email is required")),
	"subject": z.String().Min(5, z.Message("subject must be at least 5 characters long")).Max(100, z.Message("subject must be at most 100 characters long")).Required(z.Message("subject is required")),
	"message": z.String().Min(10, z.Message("message must be at least 10 characters long")).Max(1500, z.Message("message must be at most 1500 characters long")).Required(z.Message("message is required")),
})

func (h *MessageHandler) Handler(c fiber.Ctx) error {
	d := new(contactData)

	if err := c.Bind().Body(d); err != nil {
		return err
	}

	if parsed := messageSchema.Validate(d); len(parsed) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": z.Issues.Flatten(parsed),
		})
	}

	webHookData := fiber.Map{
		"username": "Portfolio",
		"content":  fmt.Sprintf("**IP:** %s\n**Name:** %s\n**Email:** %s\n**Subject:** %s\n**Message:** %s", c.IP(), d.Name, d.Email, d.Subject, d.Message),
	}

	cc := client.New()
	if _, err := cc.Post(h.cfg.DiscordWebhook, client.Config{
		Body: webHookData,
	}); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "Failed to send data to webhook",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
