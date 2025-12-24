package handlers

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse is a standardized error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// SuccessResponse is a standardized success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(APIResponse{Success: true, Data: data})
}

func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(APIResponse{Success: true, Data: data})
}

func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
		Success: false, Error: &APIError{Code: "BAD_REQUEST", Message: message},
	})
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(APIResponse{
		Success: false, Error: &APIError{Code: "UNAUTHORIZED", Message: message},
	})
}

func NotFound(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(APIResponse{
		Success: false, Error: &APIError{Code: "NOT_FOUND", Message: message},
	})
}

func ServerError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
		Success: false, Error: &APIError{Code: "SERVER_ERROR", Message: message},
	})
}

func CustomErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(APIResponse{
		Success: false, Error: &APIError{Code: "ERROR", Message: err.Error()},
	})
}
