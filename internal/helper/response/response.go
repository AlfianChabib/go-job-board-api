package response

import "github.com/gofiber/fiber/v3"

type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
}

type ErrorResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  T      `json:"errors,omitempty"`
}

// SimpleResponse untuk endpoint yang hanya mengembalikan pesan (Logout, Delete, dll.)
type SimpleResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Success untuk response yang membawa payload data
func Success[T any](c fiber.Ctx, status int, message string, data T) error {
	return c.Status(status).JSON(Response[T]{
		Success: true,
		Message: message,
		Data:    &data,
	})
}

// OK shortcut cepat untuk status 200 OK dengan data
func OK[T any](c fiber.Ctx, message string, data T) error {
	return Success(c, fiber.StatusOK, message, data)
}

// Message untuk response tanpa data (misal: Logout, Delete)
func Message(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(SimpleResponse{
		Success: true,
		Message: message,
	})
}

// Error helper untuk error umum (400, 401, 404, 500)
func Error(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(ErrorResponse[any]{
		Success: false,
		Message: message,
	})
}

// ValidationError helper khusus error validasi (biasanya status 422 atau 400)
func ValidationError[T any](c fiber.Ctx, message string, details T) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(ErrorResponse[T]{
		Success: false,
		Message: message,
		Errors:  details,
	})
}
