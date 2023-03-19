package controllers

import (
	db "Hdfc_Assignment/models"
	"Hdfc_Assignment/services"
	"Hdfc_Assignment/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateNewOrder godoc
// @Summary      Create Order
// @Description  creates a new order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        req  body      models.OrderRequest true "Order Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /orders [post]
// @Security     ApiKeyAuth
func CreateNewOrder(c *gin.Context) {
	var requestBody utils.OrderRequest
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

	res, err := services.FindUserById(userId.(primitive.ObjectID))
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	if res.Role == db.RoleEmployee {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Only User Or Admin Can Callable"
		response.SendResponse(c)
		return
	}

	if len(requestBody.ProductId) != len(requestBody.OrderQuantity) {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Product And Quantity Invalid"
		response.SendResponse(c)
		return
	}

	// currentDiscount := 10

	var totalPremium int
	var totalCost int
	var grandTotal int
	var discountPercentage int
	var discountAmount int

	//var GetProductId []primitive.ObjectID
	for index, getId := range requestBody.ProductId {
		fmt.Println("getId", getId)
		data, err := services.FindProductById(getId)

		if err != nil ||
			requestBody.OrderQuantity[index] == 0 ||
			requestBody.OrderQuantity[index] > data.Stock {
			response.StatusCode = http.StatusBadRequest
			response.Message = "Unable to place order"
			response.SendResponse(c)
			return
		}

		services.UpdateProductStock(requestBody.ProductId[index], requestBody.OrderQuantity[index])

		if data.Category == "Premium" {
			totalPremium++
		}

		totalCost += data.Price * requestBody.OrderQuantity[index]
	}

	if totalPremium >= 3 {
		discountPercentage = 10
		discountAmount = totalCost / 10
		grandTotal = totalCost - discountAmount
	} else {
		grandTotal = totalCost
		discountPercentage = 0
		discountAmount = 0
	}

	// GetUserId, _ := primitive.ObjectIDFromHex(requestBody.UserId)
	//, _ := primitive.ObjectIDFromHex(requestBody.ProductId)

	order, err := services.CreateOrder(
		// GetUserId,
		// GetProductId,

		requestBody.UserId,
		requestBody.ProductId,
		requestBody.OrderQuantity,
		totalCost,
		grandTotal,
		discountPercentage,
		discountAmount,
		requestBody.Address,
	)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	services.PlaceUserOrder(requestBody.UserId, order.ID)

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{"order": order}
	response.SendResponse(c)
}

// GetOrders godoc
// @Summary      Get Orders
// @Description  gets user orders with pagination
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        page  query    string  false  "Switch page by 'page'"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /orders [get]
// @Security     ApiKeyAuth
func GetOrders(c *gin.Context) {
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

	orders, _ := services.GetOrders(userId.(primitive.ObjectID), page, limit)
	hasPrev := page > 0
	hasNext := len(orders) > limit

	if hasNext {
		orders = orders[:limit]
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"orders": orders, "prev": hasPrev, "next": hasNext}
	response.SendResponse(c)
}

// GetOneOrder godoc
// @Summary      Get a order
// @Description  get order by id
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /orders/{id} [get]
// @Security     ApiKeyAuth
func GetOneOrder(c *gin.Context) {
	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	idHex := c.Param("id")
	orderId, _ := primitive.ObjectIDFromHex(idHex)

	userId, exists := c.Get("userId")
	if !exists {
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	//order, err := services.GetOrderFromCache(userId.(primitive.ObjectID), orderId)
	// if err == nil {
	// 	utils.SendResponseData(c, gin.H{"order": order, "cache": true})
	// 	return
	// }

	order, err := services.GetOrderById(userId.(primitive.ObjectID), orderId)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	//	go services.CacheOneOrder(userId.(primitive.ObjectID), order)

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = gin.H{"order": order}
	response.SendResponse(c)
}

// UpdateOrder godoc
// @Summary      Update a order
// @Description  updates a order by id
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id     path    string  true  "Order ID"
// @Param        req    body    models.OrderRequest true "Order Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /orders/{id} [put]
// @Security     ApiKeyAuth
func UpdateOrder(c *gin.Context) {
	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	idHex := c.Param("id")
	orderId, _ := primitive.ObjectIDFromHex(idHex)

	userId, exists := c.Get("userId")
	if !exists {
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	var orderRequest utils.OrderRequest
	_ = c.ShouldBindBodyWith(&orderRequest, binding.JSON)

	err := services.UpdateOrder(userId.(primitive.ObjectID), orderId, &orderRequest)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.SendResponse(c)
}

// DeleteOrder godoc
// @Summary      Delete a order
// @Description  deletes order by id
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /orders/{id} [delete]
// @Security     ApiKeyAuth
func DeleteOrder(c *gin.Context) {
	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	idHex := c.Param("id")
	orderId, _ := primitive.ObjectIDFromHex(idHex)

	userId, exists := c.Get("userId")
	if !exists {
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	err := services.DeleteOrder(userId.(primitive.ObjectID), orderId)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.SendResponse(c)
}
