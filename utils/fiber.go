package utils

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/reverie/types"
)

// ServerError sends internal server error messages
func ServerError(logContext string, err error) error {
	LogError(logContext, err)
	return err
}

var (
	// ErrFailedExtraction occurs when the request fails to extract claims from JWT
	ErrFailedExtraction = errors.New("Failed to extract JWT claims")
)

// ExtractClaims takes the fiber context and returns the User
// IMPORTANT: claims shall only provide the email, username and the role of a user
func ExtractClaims(c *fiber.Ctx) *types.Claims {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email, ok := claims[types.EmailKey].(string)
	if !ok {
		return nil
	}
	username, ok := claims[types.UsernameKey].(string)
	if !ok {
		return nil
	}
	role, ok := claims[types.RoleKey].(string)
	if !ok {
		return nil
	}
	return &types.Claims{
		Email:    email,
		Username: username,
		Role:     role,
	}
}

func unsafeBytes(s string) (bs []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return
}

// ImmutableString copies a string to make it immutable
// This is required becase fiber has zero allocation policy and its params are used per request
// i.e it cannot be used in goroutines
// For ex:- If you used `bar := c.Params("foo")` and then `go someFunc(bar)`, then bar will be invalid inside the goroutine
// Hence this function is necessary to make a copy of the string ensuring it stays valid even in goroutines
// IMPORTANT :- Use this function only when the extracted params are to be passed to some goroutine inside the handler
// Dont use this unnecessarily to avoid wasteful memory allocations
func ImmutableString(s string) string {
	return string(unsafeBytes(s))
}
