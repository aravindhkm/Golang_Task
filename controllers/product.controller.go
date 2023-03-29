package controllers

import (
	db "Hdfc_Assignment/models"
	"Hdfc_Assignment/services"
	"Hdfc_Assignment/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateNewProduct godoc
// @Summary      Create Product
// @Description  creates a new product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        req  body      models.ProductRequest true "Product Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /products [post]
// @Security     ApiKeyAuth
func CreateNewProduct(c *gin.Context) {
	var requestBody utils.ProductRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	userId, exists := c.Get("userId")
	if !exists {
		response.StatusCode = http.StatusBadRequest
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	res, err := services.FindAdminById(userId.(primitive.ObjectID))
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	if res.Role != db.RoleAdmin {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Only Admin Can Callable"
		response.SendResponse(c)
		return
	}

	product, err := services.CreateProduct(
		requestBody.Title,
		requestBody.Description,
		requestBody.Price,
		requestBody.Rating,
		requestBody.Stock,
		requestBody.Brand,
		requestBody.Type,
		requestBody.Category,
		requestBody.Thumbnail,
		requestBody.Images,
	)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{"product": product}
	response.SendResponse(c)
}

// GetProducts godoc
// @Summary      Get Products
// @Description  gets user products with pagination
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page  query    string  false  "Switch page by 'page'"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /products [get]
// @Security     ApiKeyAuth
func GetProducts(c *gin.Context) {
	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	userId, exists := c.Get("userId")
	if !exists {
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	pageQuery := c.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageQuery)
	limit := 5

	products, _ := services.GetProducts(userId.(primitive.ObjectID), page, limit)
	hasPrev := page > 0
	hasNext := len(products) > limit

	if hasNext {
		products = products[:limit]
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"products": products, "prev": hasPrev, "next": hasNext}
	response.SendResponse(c)
}

func GetAllProducts(c *gin.Context) {
	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	_, exists := c.Get("userId")
	if !exists {
		response.StatusCode = http.StatusBadRequest
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	pageQuery := c.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageQuery)
	limit := 5

	// isAdmin, err := services.FindAdminById(userId.(primitive.ObjectID))
	// if err != nil || isAdmin.Role != db.RoleAdmin {
	// 	response.StatusCode = http.StatusBadRequest
	// 	response.Message = "admin only accessible"
	// 	response.SendResponse(c)
	// 	return
	// }

	products, _ := services.GetAllProducts(page, limit)
	hasPrev := page > 0
	hasNext := len(products) > limit

	if hasNext {
		products = products[:limit]
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"products": products, "prev": hasPrev, "next": hasNext}
	response.SendResponse(c)
}

// GetOneProduct godoc
// @Summary      Get a product
// @Description  get product by id
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /products/{id} [get]
// @Security     ApiKeyAuth
func GetOneProduct(c *gin.Context) {
	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	idHex := c.Param("id")
	productId, _ := primitive.ObjectIDFromHex(idHex)

	_, exists := c.Get("userId")

	if !exists {
		response.StatusCode = http.StatusBadRequest
		response.Message = "cannot get product"
		response.SendResponse(c)
		return
	}

	//product, err := services.GetProductFromCache(userId.(primitive.ObjectID), productId)
	// if err == nil {
	// 	utils.SendResponseData(c, gin.H{"product": product, "cache": true})
	// 	return
	// }

	product, err := services.GetProductById(productId)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	//	go services.CacheOneProduct(userId.(primitive.ObjectID), product)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"product": product}
	response.SendResponse(c)
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  updates a product by id
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id     path    string  true  "Product ID"
// @Param        req    body    models.ProductRequest true "Product Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /products/{id} [put]
// @Security     ApiKeyAuth
func UpdateProduct(c *gin.Context) {
	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	idHex := c.Param("id")
	productId, _ := primitive.ObjectIDFromHex(idHex)

	userId, exists := c.Get("userId")
	if !exists {
		response.StatusCode = http.StatusBadRequest
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	adminData, err := services.FindAdminById(userId.(primitive.ObjectID))
	if err != nil || adminData.Role != db.RoleAdmin {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	var productRequest utils.ProductRequest
	_ = c.ShouldBindBodyWith(&productRequest, binding.JSON)

	err = services.UpdateProduct(productId, &productRequest)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.SendResponse(c)
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  deletes product by id
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /products/{id} [delete]
// @Security     ApiKeyAuth
func DeleteProduct(c *gin.Context) {
	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	idHex := c.Param("id")
	productId, _ := primitive.ObjectIDFromHex(idHex)

	userId, exists := c.Get("userId")
	if !exists {
		response.StatusCode = http.StatusBadRequest
		response.Message = "cannot get product"
		response.SendResponse(c)
		return
	}

	adminData, err := services.FindAdminById(userId.(primitive.ObjectID))
	if err != nil || adminData.Role != db.RoleAdmin {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	err = services.DeleteProduct(productId)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.SendResponse(c)
}
