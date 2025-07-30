package core

import (
	"fmt"
	"runtime"
	"strings"
)

// Enhanced error handling with descriptive messages and context

// ErrorType represents different types of errors
type ErrorType int

const (
	ErrorTypeInvalidParameter ErrorType = iota
	ErrorTypeFileNotFound
	ErrorTypeInvalidFormat
	ErrorTypeOutOfBounds
	ErrorTypeInvalidState
	ErrorTypeMemoryError
	ErrorTypeRenderError
	ErrorTypeUnsupportedOperation
)

// AdvanceError represents an enhanced error with context
type AdvanceError struct {
	Type        ErrorType
	Message     string
	Context     map[string]interface{}
	StackTrace  []string
	Suggestions []string
	Code        string
}

// Error implements the error interface
func (e *AdvanceError) Error() string {
	var builder strings.Builder
	
	// Error type and message
	builder.WriteString(fmt.Sprintf("[%s] %s", e.getTypeString(), e.Message))
	
	// Add context if available
	if len(e.Context) > 0 {
		builder.WriteString("\n\nContext:")
		for key, value := range e.Context {
			builder.WriteString(fmt.Sprintf("\n  %s: %v", key, value))
		}
	}
	
	// Add suggestions if available
	if len(e.Suggestions) > 0 {
		builder.WriteString("\n\nSuggestions:")
		for _, suggestion := range e.Suggestions {
			builder.WriteString(fmt.Sprintf("\n  • %s", suggestion))
		}
	}
	
	// Add error code if available
	if e.Code != "" {
		builder.WriteString(fmt.Sprintf("\n\nError Code: %s", e.Code))
	}
	
	return builder.String()
}

// getTypeString returns a human-readable string for the error type
func (e *AdvanceError) getTypeString() string {
	switch e.Type {
	case ErrorTypeInvalidParameter:
		return "INVALID_PARAMETER"
	case ErrorTypeFileNotFound:
		return "FILE_NOT_FOUND"
	case ErrorTypeInvalidFormat:
		return "INVALID_FORMAT"
	case ErrorTypeOutOfBounds:
		return "OUT_OF_BOUNDS"
	case ErrorTypeInvalidState:
		return "INVALID_STATE"
	case ErrorTypeMemoryError:
		return "MEMORY_ERROR"
	case ErrorTypeRenderError:
		return "RENDER_ERROR"
	case ErrorTypeUnsupportedOperation:
		return "UNSUPPORTED_OPERATION"
	default:
		return "UNKNOWN_ERROR"
	}
}

// WithContext adds context to the error
func (e *AdvanceError) WithContext(key string, value interface{}) *AdvanceError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithSuggestion adds a suggestion to the error
func (e *AdvanceError) WithSuggestion(suggestion string) *AdvanceError {
	e.Suggestions = append(e.Suggestions, suggestion)
	return e
}

// WithCode adds an error code
func (e *AdvanceError) WithCode(code string) *AdvanceError {
	e.Code = code
	return e
}

// NewError creates a new AdvanceError with stack trace
func NewError(errorType ErrorType, message string) *AdvanceError {
	return &AdvanceError{
		Type:       errorType,
		Message:    message,
		Context:    make(map[string]interface{}),
		StackTrace: captureStackTrace(),
	}
}

// captureStackTrace captures the current stack trace
func captureStackTrace() []string {
	var trace []string
	
	// Skip the first few frames (this function and error creation)
	for i := 2; i < 10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		
		// Only include frames from our package
		if strings.Contains(file, "advancegg") {
			trace = append(trace, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
		}
	}
	
	return trace
}

// Specific error creation functions

// NewInvalidParameterError creates an error for invalid parameters
func NewInvalidParameterError(paramName string, value interface{}, expected string) *AdvanceError {
	return NewError(ErrorTypeInvalidParameter, 
		fmt.Sprintf("Invalid parameter '%s'", paramName)).
		WithContext("parameter", paramName).
		WithContext("value", value).
		WithContext("expected", expected).
		WithSuggestion(fmt.Sprintf("Ensure '%s' is %s", paramName, expected)).
		WithCode("E001")
}

// NewFileNotFoundError creates an error for missing files
func NewFileNotFoundError(filepath string) *AdvanceError {
	return NewError(ErrorTypeFileNotFound,
		fmt.Sprintf("File not found: %s", filepath)).
		WithContext("filepath", filepath).
		WithSuggestion("Check if the file path is correct").
		WithSuggestion("Ensure the file exists and is readable").
		WithCode("E002")
}

// NewInvalidFormatError creates an error for invalid file formats
func NewInvalidFormatError(filepath string, expectedFormats []string) *AdvanceError {
	return NewError(ErrorTypeInvalidFormat,
		fmt.Sprintf("Invalid file format: %s", filepath)).
		WithContext("filepath", filepath).
		WithContext("expected_formats", expectedFormats).
		WithSuggestion(fmt.Sprintf("Supported formats: %s", strings.Join(expectedFormats, ", "))).
		WithCode("E003")
}

// NewOutOfBoundsError creates an error for out-of-bounds operations
func NewOutOfBoundsError(operation string, value, min, max interface{}) *AdvanceError {
	return NewError(ErrorTypeOutOfBounds,
		fmt.Sprintf("Value out of bounds in %s", operation)).
		WithContext("operation", operation).
		WithContext("value", value).
		WithContext("min", min).
		WithContext("max", max).
		WithSuggestion(fmt.Sprintf("Ensure value is between %v and %v", min, max)).
		WithCode("E004")
}

// NewInvalidStateError creates an error for invalid context state
func NewInvalidStateError(operation string, currentState string, requiredState string) *AdvanceError {
	return NewError(ErrorTypeInvalidState,
		fmt.Sprintf("Invalid state for operation '%s'", operation)).
		WithContext("operation", operation).
		WithContext("current_state", currentState).
		WithContext("required_state", requiredState).
		WithSuggestion(fmt.Sprintf("Ensure context is in '%s' state before calling %s", requiredState, operation)).
		WithCode("E005")
}

// NewMemoryError creates an error for memory-related issues
func NewMemoryError(operation string, requestedSize int64) *AdvanceError {
	return NewError(ErrorTypeMemoryError,
		fmt.Sprintf("Memory allocation failed for %s", operation)).
		WithContext("operation", operation).
		WithContext("requested_size", requestedSize).
		WithSuggestion("Try reducing image size or clearing caches").
		WithSuggestion("Check available system memory").
		WithCode("E006")
}

// NewRenderError creates an error for rendering issues
func NewRenderError(operation string, details string) *AdvanceError {
	return NewError(ErrorTypeRenderError,
		fmt.Sprintf("Rendering failed for %s: %s", operation, details)).
		WithContext("operation", operation).
		WithContext("details", details).
		WithSuggestion("Check if all required resources are available").
		WithSuggestion("Verify operation parameters are valid").
		WithCode("E007")
}

// NewUnsupportedOperationError creates an error for unsupported operations
func NewUnsupportedOperationError(operation string, reason string) *AdvanceError {
	return NewError(ErrorTypeUnsupportedOperation,
		fmt.Sprintf("Unsupported operation: %s", operation)).
		WithContext("operation", operation).
		WithContext("reason", reason).
		WithSuggestion("Check documentation for supported operations").
		WithCode("E008")
}

// Validation helpers

// ValidatePositive validates that a number is positive
func ValidatePositive(name string, value float64) error {
	if value <= 0 {
		return NewInvalidParameterError(name, value, "a positive number")
	}
	return nil
}

// ValidateRange validates that a value is within a range
func ValidateRange(name string, value, min, max float64) error {
	if value < min || value > max {
		return NewOutOfBoundsError(name, value, min, max)
	}
	return nil
}

// ValidateNotNil validates that a pointer is not nil
func ValidateNotNil(name string, value interface{}) error {
	if value == nil {
		return NewInvalidParameterError(name, value, "not nil")
	}
	return nil
}

// ValidateImageBounds validates coordinates are within image bounds
func ValidateImageBounds(x, y, width, height int) error {
	if x < 0 || y < 0 {
		return NewOutOfBoundsError("coordinate", fmt.Sprintf("(%d, %d)", x, y), "(0, 0)", fmt.Sprintf("(%d, %d)", width-1, height-1))
	}
	if x >= width || y >= height {
		return NewOutOfBoundsError("coordinate", fmt.Sprintf("(%d, %d)", x, y), "(0, 0)", fmt.Sprintf("(%d, %d)", width-1, height-1))
	}
	return nil
}

// Context validation methods

// ValidateContext validates that a context is in a valid state
func (dc *Context) ValidateContext() error {
	if dc == nil {
		return NewInvalidStateError("context_operation", "nil", "initialized")
	}
	if dc.im == nil {
		return NewInvalidStateError("drawing_operation", "no_image", "image_initialized")
	}
	return nil
}

// ValidateFontLoaded validates that a font is loaded
func (dc *Context) ValidateFontLoaded() error {
	if dc.fontFace == nil {
		return NewInvalidStateError("text_operation", "no_font", "font_loaded").
			WithSuggestion("Load a font using LoadFontFace() before drawing text")
	}
	return nil
}

// Error recovery helpers

// RecoverFromPanic recovers from panics and converts them to AdvanceErrors
func RecoverFromPanic() error {
	if r := recover(); r != nil {
		return NewError(ErrorTypeRenderError, fmt.Sprintf("Panic recovered: %v", r)).
			WithSuggestion("This is likely a bug - please report it").
			WithCode("E999")
	}
	return nil
}

// SafeExecute executes a function with panic recovery
func SafeExecute(operation string, fn func() error) error {
	defer func() {
		if r := recover(); r != nil {
			// Log the panic for debugging
			fmt.Printf("Panic in %s: %v\n", operation, r)
		}
	}()
	
	return fn()
}

// ErrorHandler is a function type for handling errors
type ErrorHandler func(error)

// DefaultErrorHandler is the default error handler
var DefaultErrorHandler ErrorHandler = func(err error) {
	if err != nil {
		fmt.Printf("AdvanceGG Error: %v\n", err)
	}
}

// SetErrorHandler sets the global error handler
func SetErrorHandler(handler ErrorHandler) {
	DefaultErrorHandler = handler
}

// HandleError handles an error using the global error handler
func HandleError(err error) {
	if err != nil && DefaultErrorHandler != nil {
		DefaultErrorHandler(err)
	}
}
