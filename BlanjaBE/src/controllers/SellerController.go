package controllers

import (
	"gofiber-marketplace/src/helpers"
	"gofiber-marketplace/src/middlewares"
	"gofiber-marketplace/src/models"
	"gofiber-marketplace/src/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetSellers(c *fiber.Ctx) error {
	keyword := c.Query("search")
	sort := helpers.GetSortParams(c.Query("sorting"), c.Query("orderBy"))

	sellers := models.SelectAllSellers(keyword, sort)
	if len(sellers) == 0 {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"status":     "no content",
			"statusCode": 202,
			"message":    "Seller is empty",
		})
	}

	resultSellers := make([]map[string]interface{}, len(sellers))
	for i, seller := range sellers {
		products := make([]map[string]interface{}, len(seller.Products))
		for j, product := range seller.Products {
			var categoryName string
			if product.CategoryID != 0 {
				category := models.SelectCategoryById(int(product.CategoryID))
				if category.ID != 0 {
					categoryName = category.Name
				}
			}

			products[j] = map[string]interface{}{
				"id":            product.ID,
				"created_at":    product.CreatedAt,
				"updated_at":    product.UpdatedAt,
				"name":          product.Name,
				"price":         product.Price,
				"rating":        product.Rating,
				"category_name": categoryName,
			}
		}

		resultSellers[i] = map[string]interface{}{
			"id":         seller.ID,
			"created_at": seller.CreatedAt,
			"updated_at": seller.UpdatedAt,
			"name":       seller.Name,
			"user_id":    seller.User.ID,
			"email":      seller.User.Email,
			"image":      seller.Image,
			"phone":      seller.Phone,
			"desc":       seller.Description,
			"role":       seller.User.Role,
			"products":   products,
		}
	}

	// return c.JSON(categories)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"data":       resultSellers,
	})
}

func GetDetailSeller(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid ID format",
		})
	}

	seller := models.SelectSellerById(id)
	if seller.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Seller not found",
		})
	}

	if seller.User.Role != "seller" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Role of this user is customer or not seller",
		})
	}

	products := make([]map[string]interface{}, len(seller.Products))
	for i, product := range seller.Products {
		var categoryName string
		if product.CategoryID != 0 {
			category := models.SelectCategoryById(int(product.CategoryID))
			if category.ID != 0 {
				categoryName = category.Name
			}
		}

		products[i] = map[string]interface{}{
			"id":            product.ID,
			"created_at":    product.CreatedAt,
			"updated_at":    product.UpdatedAt,
			"name":          product.Name,
			"price":         product.Price,
			"rating":        product.Rating,
			"category_name": categoryName,
		}
	}

	resultSeller := map[string]interface{}{
		"id":         seller.ID,
		"created_at": seller.CreatedAt,
		"updated_at": seller.UpdatedAt,
		"name":       seller.Name,
		"user_id":    seller.User.ID,
		"email":      seller.User.Email,
		"image":      seller.Image,
		"phone":      seller.Phone,
		"desc":       seller.Description,
		"role":       seller.User.Role,
		"products":   products,
	}

	// return c.JSON(product)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"data":       resultSeller,
	})
}

func GetSellerProfile(c *fiber.Ctx) error {
	id, err := middlewares.JWTAuthorize(c, "seller")
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"status":     fiberErr.Message,
				"statusCode": fiberErr.Code,
				"message":    fiberErr.Message,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "Internal Server Error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    err.Error(),
		})
	}

	seller := models.SelectSellerByUserId(int(id))
	if seller.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Seller not found",
		})
	}

	products := make([]map[string]interface{}, len(seller.Products))
	for i, product := range seller.Products {
		var categoryName string
		if product.CategoryID != 0 {
			category := models.SelectCategoryById(int(product.CategoryID))
			if category.ID != 0 {
				categoryName = category.Name
			}
		}

		products[i] = map[string]interface{}{
			"id":            product.ID,
			"created_at":    product.CreatedAt,
			"updated_at":    product.UpdatedAt,
			"name":          product.Name,
			"price":         product.Price,
			"rating":        product.Rating,
			"stock":         product.Stock,
			"category_name": categoryName,
		}
	}

	resultAddresses := models.SelectAddressesbyUserID(int(id))
	addresses := make([]map[string]interface{}, len(resultAddresses))
	for i, address := range resultAddresses {
		addresses[i] = map[string]interface{}{
			"id":             address.ID,
			"created_at":     address.CreatedAt,
			"updated_at":     address.UpdatedAt,
			"main_address":   address.MainAddress,
			"detail_address": address.DetailAddress,
			"postal_code":    address.PostalCode,
			"name":           address.Name,
			"phone":          address.Phone,
			"primary":        address.Primary,
			"city":           address.City,
		}
	}

	resultSeller := map[string]interface{}{
		"id":         seller.ID,
		"created_at": seller.CreatedAt,
		"updated_at": seller.UpdatedAt,
		"name":       seller.Name,
		"user_id":    seller.User.ID,
		"email":      seller.User.Email,
		"image":      seller.Image,
		"phone":      seller.Phone,
		"desc":       seller.Description,
		"role":       seller.User.Role,
		"products":   products,
		"addresses":  addresses,
	}

	// return c.JSON(product)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"data":       resultSeller,
	})
}

func UpdateSellerProfile(c *fiber.Ctx) error {
	var profileData models.SellerProfile

	id, err := middlewares.JWTAuthorize(c, "seller")
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"status":     fiberErr.Message,
				"statusCode": fiberErr.Code,
				"message":    fiberErr.Message,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "Internal Server Error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    err.Error(),
		})
	}

	seller := models.SelectSellerByUserId(int(id))
	if seller.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Seller not found",
		})
	}

	if err := c.BodyParser(&profileData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	user := middlewares.XSSMiddleware(&profileData).(*models.SellerProfile)
	if errors := helpers.StructValidation(user); len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":     "unprocessable entity",
			"statusCode": 422,
			"message":    "Validation failed",
			"errors":     errors,
		})
	}

	if seller.User.Email != user.Email {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Email already exists",
		})
	}

	updatedUser := models.User{
		Email: user.Email,
	}

	updatedSeller := models.Seller{
		Name:        user.Name,
		Phone:       user.Phone,
		Description: user.Description,
	}

	if err := models.UpdateUser(int(id), &updatedUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to update user",
		})
	}

	if err := models.UpdateSeller(int(seller.ID), &updatedSeller); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to update seller",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    "Profile updated successfully",
	})

}

func UpdateSellerProfilePhoto(c *fiber.Ctx) error {
	id, err := middlewares.JWTAuthorize(c, "seller")
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"status":     fiberErr.Message,
				"statusCode": fiberErr.Code,
				"message":    fiberErr.Message,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "Internal Server Error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    err.Error(),
		})
	}

	seller := models.SelectSellerByUserId(int(id))
	if seller.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Seller not found",
		})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Failed to upload file",
		})
	}

	if err := helpers.ImageValidation(file); len(err) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":     "unprocessable entity",
			"statusCode": 422,
			"message":    "Validation failed",
			"errors":     err,
		})
	}

	uploadResult, err := services.UploadCloudinary(c, file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to save file",
		})
	}

	updatedSeller := models.Seller{
		Image: uploadResult.URL,
	}

	if err := models.UpdateSeller(int(seller.ID), &updatedSeller); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to update profile photo",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    "Profile photo updated successfully",
	})

}

func DeleteSeller(c *fiber.Ctx) error {
	id, err := middlewares.JWTAuthorize(c, "seller")
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"status":     fiberErr.Message,
				"statusCode": fiberErr.Code,
				"message":    fiberErr.Message,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "Internal Server Error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    err.Error(),
		})
	}

	seller := models.SelectSellerByUserId(int(id))
	if seller.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Seller not found",
		})
	}

	if err := models.DeleteSeller(int(seller.ID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to delete seller",
		})
	}

	if err := models.DeleteUser(int(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "server error",
			"statusCode": 500,
			"message":    "Failed to delete user",
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":     "success",
			"statusCode": 200,
			"message":    "User deleted successfully",
		})
	}
}
