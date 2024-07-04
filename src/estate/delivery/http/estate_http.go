package http

import (
	"encoding/json"
	"net/http"

	"github.com/davidyunus/sawitpro-estate/src/domain"
	"github.com/davidyunus/sawitpro-estate/src/helper"
	"github.com/labstack/echo/v4"
)

type estateHandler struct {
	estateUsecase domain.EstateUsecase
}

func NewEstateHandler(e *echo.Echo, estateUsecase domain.EstateUsecase) {
	handler := estateHandler{
		estateUsecase: estateUsecase,
	}

	e.POST("/estate", handler.CreateEstate)
	e.POST(`/estate/:id/tree`, handler.PlantPalmTree)
	e.GET("/estate/:id/stats", handler.GetTreeStats)
	e.GET("/estate/:id/drone-plan", handler.GetDroneFlyingDistance)
}

// @Summary Create Estate
// @Description Create Estate
// @Tags estates
// @Accept  json
// @Produce  json
// @Param   estate body domain.Estate true "Estate Payload"
// @Success 201 {object} domain.Estate
// @Failure 400 {object} helper.HttpResponse
// @Router /estate [post]
func (e *estateHandler) CreateEstate(c echo.Context) error {
	ctx := c.Request().Context()

	payload := &domain.Estate{}
	err := json.NewDecoder(c.Request().Body).Decode(&payload)
	if err != nil {
		err = domain.ErrInvalidInput
		code := helper.GetStatusCode(err)
		return c.JSON(code, helper.Response(code, err.Error(), nil, err.Error()))
	}
	c.Echo().Validator = helper.NewValidator()
	err = c.Echo().Validator.Validate(payload)
	if err != nil {
		err = domain.ErrInvalidInput
		code := helper.GetStatusCode(err)
		return c.JSON(code, helper.Response(code, err.Error(), nil, err.Error()))
	}

	resp, err := e.estateUsecase.CreateEstate(ctx, payload)
	if err != nil {
		response := helper.Response(http.StatusBadRequest, err.Error(), nil, err.Error())
		return c.JSON(response.Code, response)
	}

	response := helper.Response(http.StatusCreated, "Success create estate", resp, nil)
	return c.JSON(http.StatusCreated, response)
}

// @Summary Plant Palm Tree
// @Description Plant a palm tree in an estate
// @Tags estates
// @Accept  json
// @Produce  json
// @Param   id    path  string           true "Estate ID"
// @Param   tree  body  domain.PalmTree  true "Palm Tree Payload"
// @Success 201 {object} nil
// @Failure 400 {object} helper.HttpResponse
// @Router /estate/{id}/tree [post]
func (e *estateHandler) PlantPalmTree(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	payload := &domain.PalmTree{}
	err := json.NewDecoder(c.Request().Body).Decode(&payload)
	if err != nil {
		err = domain.ErrInvalidInput
		code := helper.GetStatusCode(err)
		return c.JSON(code, helper.Response(code, err.Error(), nil, err.Error()))
	}
	c.Echo().Validator = helper.NewValidator()
	err = c.Echo().Validator.Validate(payload)
	if err != nil {
		err = domain.ErrInvalidInput
		code := helper.GetStatusCode(err)
		return c.JSON(code, helper.Response(code, err.Error(), nil, err.Error()))
	}
	resp, err := e.estateUsecase.PlantPalmTree(ctx, id, payload)
	if err != nil {
		code := helper.GetStatusCode(err)
		response := helper.Response(code, err.Error(), nil, err.Error())
		return c.JSON(response.Code, response)
	}

	response := helper.Response(http.StatusCreated, "Success plant palm tree", resp, nil)
	return c.JSON(http.StatusCreated, response)
}

// @Summary Get Tree Stats
// @Description Get statistics of trees in an estate
// @Tags estates
// @Accept  json
// @Produce  json
// @Param   id    path  string  true "Estate ID"
// @Success 200 {object} domain.GetTreeStatsResponse
// @Failure 400 {object} helper.HttpResponse
// @Router /estate/{id}/stats [get]
func (e *estateHandler) GetTreeStats(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	treeStats, err := e.estateUsecase.GetTreeStats(ctx, id)
	if err != nil {
		code := helper.GetStatusCode(err)
		response := helper.Response(code, err.Error(), nil, err.Error())
		return c.JSON(response.Code, response)
	}

	response := helper.Response(http.StatusOK, "Success get tree stats", treeStats, nil)
	return c.JSON(http.StatusCreated, response)
}

// @Summary Get Drone Flying Distance
// @Description Get the flying distance plan for a drone in an estate
// @Tags estates
// @Accept  json
// @Produce  json
// @Param   id    path  string  true "Estate ID"
// @Success 200 {object} domain.GetDroneFlyingDistanceResponse
// @Failure 400 {object} helper.HttpResponse
// @Router /estate/{id}/drone-plan [get]
func (e *estateHandler) GetDroneFlyingDistance(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	treeStats, err := e.estateUsecase.GetDroneFlyingDistance(ctx, id)
	if err != nil {
		code := helper.GetStatusCode(err)
		response := helper.Response(code, err.Error(), nil, err.Error())
		return c.JSON(response.Code, response)
	}

	response := helper.Response(http.StatusOK, "Success get drone flying distance", treeStats, nil)
	return c.JSON(http.StatusCreated, response)
}
