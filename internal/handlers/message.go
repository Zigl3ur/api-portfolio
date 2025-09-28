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
		return errors.New("Name must be at least 2 char long")
	} else if nameLength > 15 {
		return errors.New("Name must be at most 15 char long")
	}

	if _, err := mail.ParseAddress(data.Email); err != nil {
		return errors.New("Email is invalid")
	}

	if subjectLength < 5 {
		return errors.New("Subject must be at least 5 char long")
	} else if subjectLength > 100 {
		return errors.New("Name must be at most 100 char long")
	}
	if messageLength < 10 {
		return errors.New("Message must be at least 10 char long")
	} else if messageLength > 1500 {
		return errors.New("Message must be at most 1500 char long")
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
		c.Status(fiber.StatusBadRequest)
		c.JSON(fiber.Map{
			"error": err.Error(),
		})
		return nil
	}

	webHookData := fiber.Map{
		"username": "Portfolio",
		"content":  fmt.Sprintf("**Name:** %s\n**Email:** %s\n**Subject:** %s\n**Message:** %s", d.Name, d.Email, d.Subject, d.Message),
	}

	payload, err := c.App().Config().JSONEncoder(webHookData)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		c.JSON(fiber.Map{
			"error": "Failed to encode webhook data",
		})
		return nil
	}

	_, err = http.Post(webhookUrl, "application/json", bytes.NewReader(payload))
	if err != nil {
		c.Status(fiber.StatusBadGateway)
		c.JSON(fiber.Map{
			"error": "Failed to send data to webhook",
		})
		return nil
	}

	c.SendStatus(fiber.StatusOK)
	return nil
}
