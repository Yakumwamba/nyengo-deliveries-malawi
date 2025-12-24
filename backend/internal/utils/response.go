package utils

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"pageSize,omitempty"`
	TotalCount int `json:"totalCount,omitempty"`
	TotalPages int `json:"totalPages,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{Success: true, Data: data})
}

func ErrorResponse(c *fiber.Ctx, code int, errCode, message string) error {
	return c.Status(code).JSON(Response{
		Success: false,
		Error:   &Error{Code: errCode, Message: message},
	})
}

func PaginatedResponse(c *fiber.Ctx, data interface{}, meta *Meta) error {
	return c.JSON(Response{Success: true, Data: data, Meta: meta})
}
