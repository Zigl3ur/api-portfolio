package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
)

type contactData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func validateData(data *contactData) error {
	nameLength := utf8.RuneCountInString(data.Name)
	subjectLength := utf8.RuneCountInString(data.Subject)
	messageLength := utf8.RuneCountInString(data.Message)

	if nameLength < 2 {
		return errors.New("name must be at least 2 char long")
	} else if nameLength > 15 {
		return errors.New("name must be at most 15 char long")
	}

	if _, err := mail.ParseAddress(data.Email); err != nil {
		return errors.New("email is invalid")
	}

	if subjectLength < 5 {
		return errors.New("subject must be at least 5 char long")
	} else if subjectLength > 100 {
		return errors.New("subject must be at most 100 char long")
	}
	if messageLength < 10 {
		return errors.New("message must be at least 10 char long")
	} else if messageLength > 1500 {
		return errors.New("message must be at most 1500 char long")
	}

	return nil
}

func MessageHandler(c *fiber.Ctx, webhookUrl string) error {
	c.Accepts("application/json")

	d := &contactData{}

	if err := c.BodyParser(d); err != nil {
		return err
	}

	err := validateData(d)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	webHookData := fiber.Map{
		"username": "Portfolio",
		"content":  fmt.Sprintf("**IP:** %s\n**Name:** %s\n**Email:** %s\n**Subject:** %s\n**Message:** %s", c.IP(), d.Name, d.Email, d.Subject, d.Message),
	}

	payload, err := c.App().Config().JSONEncoder(webHookData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to encode webhook data",
		})
	}

	_, err = http.Post(webhookUrl, "application/json", bytes.NewReader(payload))
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "Failed to send data to webhook",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
