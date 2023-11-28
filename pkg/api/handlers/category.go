package handlers

import (
	services "Zhooze/pkg/usecase/interface"
	"Zhooze/pkg/utils/models"
	"Zhooze/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryUseCase services.CategoryUseCase
}

func NewCategoryHandler(useCase services.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		categoryUseCase: useCase,
	}
}

// @Summary		Get Category
// @Description	Retrieve All Category
// @Tags			Admin Category Management
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Success		200		{object}	response.Response{}
// @Failure		500		{object}	response.Response{}
// @Router			/admin/category   [GET]
func (cy *CategoryHandler) GetCategory(c *gin.Context) {
	category, err := cy.categoryUseCase.GetCategory()
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Couldn't displayed categories", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
	}
	success := response.ClientResponse(http.StatusBadRequest, "Display All category", category, nil)
	c.JSON(http.StatusBadRequest, success)
}

// admin
// @Summary		Add Category
// @Description	Admin can add new categories for products
// @Tags			Admin Category Management
// @Accept			json
// @Produce		    json
// @Param			category	body	models.Category	true	"category"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/admin/category [POST]
func (cy *CategoryHandler) AddCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	cate, err := cy.categoryUseCase.AddCategory(category)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Could not add the Category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully added Category", cate, nil)
	c.JSON(http.StatusOK, success)

}

// @Summary		Delete Category
// @Description	Admin can delete a category
// @Tags			Admin Category Management
// @Accept			json
// @Produce		    json
// @Param			id	query	string	true	"id"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/admin/category     [DELETE]
func (cy *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err = cy.categoryUseCase.DeleteCategory(id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not delete the specified category", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully deleted the category", nil, nil)
	c.JSON(http.StatusOK, success)

}

// @Summary		Update Category
// @Description	Admin can update name of a category into new name
// @Tags			Admin Category Management
// @Accept			json
// @Produce		    json
// @Param			set_new_name	body	models.SetNewName	true	"set new name"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/admin/category     [PUT]
func (cy *CategoryHandler) UpdateCategory(c *gin.Context) {
	var categoryUpdate models.SetNewName
	if err := c.ShouldBindJSON(&categoryUpdate); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	ok, err := cy.categoryUseCase.UpdateCategory(categoryUpdate.Current, categoryUpdate.New)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "could not update the product", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully updated Category", ok, nil)
	c.JSON(http.StatusOK, success)

}
