package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func isURI(str string) bool {
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}

	if u.Scheme == "" {
		return false
	}

	if u.Host == "" {
		return false
	}

	if !strings.Contains(u.Host, ".") {
		return false
	}

	return true
}

// ValidationMiddleware creates a middleware for request validation
func ValidationMiddleware() gin.HandlerFunc {
	validate := validator.New()

	validate.RegisterValidation("uri", func(fl validator.FieldLevel) bool {
		uri := fl.Field().String()
		return isURI(uri)
	})

	return gin.HandlerFunc(func(c *gin.Context) {
		c.Set("validator", validate)
		c.Next()
	})
}

// ValidateStruct validates a struct and returns formatted error messages
func ValidateStruct(validate *validator.Validate, s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())
		tag := err.Tag()

		switch tag {
		case "required":
			errors[field] = field + " is required"
		case "email":
			errors[field] = field + " must be a valid email"
		case "min":
			if field == "weight" || field == "height" {
				errors[field] = field + " must be at least " + err.Param()
			} else {
				errors[field] = field + " must be at least " + err.Param() + " characters"
			}
		case "max":
			if field == "weight" || field == "height" {
				errors[field] = field + " must be at most " + err.Param()
			} else {
				errors[field] = field + " must be at most " + err.Param() + " characters"
			}
		case "oneof":
			errors[field] = field + " must be one of: " + err.Param()
		case "url", "uri":
			errors[field] = field + " must be a valid URL"
		default:
			errors[field] = field + " is invalid"
		}
	}

	return errors
}

// HandleValidationError handles validation errors and sends appropriate response
func HandleValidationError(c *gin.Context, err error) bool {
	if err != nil {
		// Check if it's a JSON binding error (empty body, malformed JSON, etc.)
		if strings.Contains(err.Error(), "EOF") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is required"})
			return true
		}
		if strings.Contains(err.Error(), "invalid character") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return true
		}
		// Generic binding error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return true
	}
	return false
}

// ContentTypeMiddleware validates Content-Type header for POST/PATCH requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == "POST" || method == "PATCH" || method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if contentType == "" || !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be application/json"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func CheckForExplicitNulls(body []byte, fields []string) error {
	var rawJSON map[string]interface{}
	if err := json.Unmarshal(body, &rawJSON); err != nil {
		return nil // Let the normal JSON binding handle this error
	}

	for _, field := range fields {
		if value, exists := rawJSON[field]; exists && value == nil {
			return &ValidationError{Field: field, Message: field + " cannot be null"}
		}
	}

	return nil
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// ValidateJSONForNulls middleware that checks for explicit null values
func ValidateJSONForNulls(fieldsToCheck []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PATCH" || c.Request.Method == "PUT" {
			// Read the body
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
				c.Abort()
				return
			}

			// Restore the body for subsequent handlers
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			// Check for explicit nulls
			if err := CheckForExplicitNulls(body, fieldsToCheck); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
