package controllers

import (
	"time"

	validator "github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/reverie/configs"
	"github.com/reverie/models/mongo"
	"github.com/reverie/types"
	"github.com/reverie/utils"
)

// registerUser handles registration of new users
func registerUser(c *fiber.Ctx, role string) error {
	user := &types.User{}
	if err := c.BodyParser(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if result, err := validator.ValidateStruct(user); !result {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	unique, err := mongo.IsUniqueEmail(user.GetEmail())
	if err != nil {
		return utils.ServerError("User-Controller-1", err)
	}
	if !unique {
		return fiber.NewError(fiber.StatusBadRequest, "Email already registered")
	}

	hashedPass, err := utils.HashPassword(user.GetPassword())
	if err != nil {
		return utils.ServerError("User-Controller-2", err)
	}
	user.SetPassword(hashedPass)
	user.SetRole(role)

	if _, err = mongo.RegisterUser(user); err != nil {
		return utils.ServerError("User-Controller-3", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}

// RegisterClient registers a new client
func RegisterClient(c *fiber.Ctx) error {
	return registerUser(c, types.Client)
}

// RegisterVendor registers a new vendor
func RegisterVendor(c *fiber.Ctx) error {
	return registerUser(c, types.Vendor)
}

// GetUserInfo gets info regarding particular user
func GetUserInfo(c *fiber.Ctx) error {
	user, err := mongo.FetchSingleUserWithoutPassword(c.Params("user"))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.NewError(fiber.StatusNotFound, "No such user exists")
		}
		return utils.ServerError("User-Controller-4", err)
	}
	user.SetSuccess(true)
	return c.Status(fiber.StatusOK).JSON(user)
}

// GetLoggedInUserInfo returns info regarding the current logged in user
func GetLoggedInUserInfo(c *fiber.Ctx) error {
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("User-Controller-5", utils.ErrFailedExtraction)
	}
	user, err := mongo.FetchSingleUserWithoutPassword(claims.GetEmail())
	if err != nil {
		return utils.ServerError("User-Controller-6", err)
	}
	user.SetSuccess(true)
	return c.Status(fiber.StatusOK).JSON(user)
}

// UpdatePassword updates the password of a user
func UpdatePassword(c *fiber.Ctx) error {
	passwordUpdate := &types.PasswordUpdate{}
	if err := c.BodyParser(passwordUpdate); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("User-Controller-7", utils.ErrFailedExtraction)
	}
	user, err := mongo.FetchSingleUser(claims.GetEmail())
	if err != nil {
		return utils.ServerError("User-Controller-8", err)
	}
	if !utils.CompareHashWithPassword(user.GetPassword(), passwordUpdate.GetOldPassword()) {
		return fiber.NewError(fiber.StatusUnauthorized, "Old password is invalid")
	}
	hashedPass, err := utils.HashPassword(passwordUpdate.GetNewPassword())
	if err != nil {
		return utils.ServerError("User-Controller-9", err)
	}
	if err = mongo.UpdatePassword(user.GetEmail(), hashedPass); err != nil {
		return utils.ServerError("User-Controller-10", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}

// DeleteUser deletes the user from database
func DeleteUser(c *fiber.Ctx) error {
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("User-Controller-11", utils.ErrFailedExtraction)
	}
	filter := types.M{
		mongo.EmailKey: claims.GetEmail(),
	}
	updatePayload := types.M{
		"deleted": true,
	}
	err := mongo.UpdateUser(filter, updatePayload)
	if err != nil {
		return utils.ServerError("User-Controller-12", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}

// UpdateInventory updates the inventory for a vendor user
func UpdateInventory(c *fiber.Ctx) error {
	claims := utils.ExtractClaims(c)
	if claims == nil {
		return utils.ServerError("User-Controller-13", utils.ErrFailedExtraction)
	}
	inventory := &types.Inventory{}
	if err := c.BodyParser(inventory); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := mongo.UpdateVendorInventory(claims.GetEmail(), inventory); err != nil {
		return utils.ServerError("User-Controller-14", err)
	}
	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
	})
}

// Login handles the user login process
func Login(c *fiber.Ctx) error {
	auth := &types.Login{}
	if err := c.BodyParser(auth); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	user, err := mongo.FetchSingleUser(auth.GetEmail())
	if err != nil {
		return utils.ServerError("User-Controller-15", err)
	}
	if !utils.CompareHashWithPassword(user.GetPassword(), auth.GetPassword()) {
		return fiber.NewError(fiber.StatusUnauthorized, "Incorrect Email or Password")
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims[types.EmailKey] = user.GetEmail()
	claims[types.UsernameKey] = user.GetName()
	claims[types.RoleKey] = user.GetRole()

	expiry := time.Now().Add(time.Second * configs.JWTConfig.Timeout).Unix()
	claims["exp"] = expiry

	// Generate encoded token and send it as response.
	encryptedToken, err := token.SignedString([]byte(configs.JWTConfig.Secret))
	if err != nil {
		return utils.ServerError("User-Controller-16", err)
	}

	return c.Status(fiber.StatusOK).JSON(types.M{
		types.Success: true,
		"token":       encryptedToken,
		"expiry":      expiry,
	})
}
