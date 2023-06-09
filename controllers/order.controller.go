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

	userData, err := services.FindUserById(userId.(primitive.ObjectID))
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	if userData.Role == db.RoleEmployee {
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
		// fmt.Println("getId", getId)
		data, err := services.FindProductById(getId)

		if err != nil ||
			requestBody.OrderQuantity[index] <= 0 ||
			requestBody.OrderQuantity[index] > data.Stock {
			response.StatusCode = http.StatusBadRequest
			response.Message = "Unable to place order"
			response.SendResponse(c)
			return
		}

		if data.Category == "Premium" {
			totalPremium++
		}

		totalCost += data.Price * requestBody.OrderQuantity[index]
	}

	err = services.UpdateMultipleProductStock(requestBody.ProductId, requestBody.OrderQuantity)

	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "order data is invalid"
		response.SendResponse(c)
		return
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
		userData.ID,
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

	services.PlaceUserOrder(userData.ID, order.ID)

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{"order": order}
	response.SendResponse(c)
}

func getZeroId() primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex("")

	return id

}

func CancelOrder(c *gin.Context) {
	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	idHex := c.Param("id")
	orderId, _ := primitive.ObjectIDFromHex(idHex)

	userId, exists := c.Get("userId")
	if !exists {
		response.StatusCode = http.StatusBadRequest
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	userData, err := services.FindUserById(userId.(primitive.ObjectID))
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	orderData, err := services.FindOrderById(orderId)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	if userData.Role == db.RoleEmployee {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Only User Or Admin Can Callable"
		response.SendResponse(c)
		return
	}

	if orderData.UserId != userData.ID || orderData.OrderStatus != "Placed" {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Your not authorized to cancel the order or Invalid Order Status"
		response.SendResponse(c)
		return
	}

	for index, getId := range orderData.ProductId {
		services.UpdateProductCancelOrder(getId, orderData.OrderQuantity[index])
	}

	order, err := services.SetOrderStatus(orderId, getZeroId(), "Cancelled")
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	services.CancelUserOrder(orderData.UserId, orderId)

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{"order": order}
	response.SendResponse(c)
}

func DispatchOrder(c *gin.Context) {
	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	idHex := c.Param("id")
	orderId, _ := primitive.ObjectIDFromHex(idHex)

	userId, exists := c.Get("userId")
	if !exists {
		response.StatusCode = http.StatusBadRequest
		response.Message = "cannot get order"
		response.SendResponse(c)
		return
	}

	empData, err := services.FindEmployeeById(userId.(primitive.ObjectID))
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	orderData, err := services.FindOrderById(orderId)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	if orderData.OrderStatus != "Placed" {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Invalid Order Status"
		response.SendResponse(c)
		return
	}

	if empData.Role == db.RoleUser {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Only Employee Or Admin Can Callable"
		response.SendResponse(c)
		return
	}

	order, err := services.SetOrderStatus(orderId, empData.ID, "Dispatched")
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{"order": order}
	response.SendResponse(c)
}

func CompleteOrder(c *gin.Context) {
	response := &utils.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	idHex := c.Param("id")
	orderId, _ := primitive.ObjectIDFromHex(idHex)

	userId, exists := c.Get("userId")
	if !exists {
		response.StatusCode = http.StatusBadRequest
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	userData, err := services.FindUserById(userId.(primitive.ObjectID))
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	if userData.Role == db.RoleEmployee {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Only User Or Admin Can Callable"
		response.SendResponse(c)
		return
	}

	orderData, err := services.FindOrderById(orderId)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	if orderData.OrderStatus != "Dispatched" {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Invalid Order Status"
		response.SendResponse(c)
		return
	}

	order, err := services.SetOrderStatus(orderId, getZeroId(), "Completed")
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	services.CompleteUserOrder(userData.ID, orderId)

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

	_, exists := c.Get("userId")
	if !exists {
		response.Message = "cannot get user"
		response.SendResponse(c)
		return
	}

	pageQuery := c.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageQuery)
	limit := 5

	orders, _ := services.GetOrders(page, limit)
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

	_, exists := c.Get("userId")
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

	order, err := services.GetOrderById(orderId)
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
