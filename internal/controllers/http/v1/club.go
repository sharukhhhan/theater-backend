package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"theater/internal/models"
	"theater/internal/repository/postgres"
	"theater/internal/service"
	"theater/pkg/validation"
)

type ClubController struct {
	clubService service.Club
}

func NewClubRoutes(g *echo.Group, service service.Service) {
	clubController := &ClubController{
		clubService: service.Club,
	}

	authMiddleware := &AuthMiddleware{
		authService: service.Auth,
	}

	g.GET("/clubs", clubController.GetAllClubs)
	g.GET("/clubs/:id", clubController.GetClubByID)

	withMiddleware := g.Group("")
	withMiddleware.Use(authMiddleware.IsAdmin)

	withMiddleware.POST("/clubs", clubController.CreateClub)
	withMiddleware.PUT("/clubs/:id", clubController.UpdateClub)
	withMiddleware.DELETE("/clubs/:id", clubController.DeleteClub)
}

// CreateClub godoc
// @Summary Create a new club
// @Description Create a new club with the provided club details
// @Tags clubs
// @Accept  json
// @Produce  json
// @Param input body models.Club true "Club details"
// @Success 200 {object} map[string]number
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /clubs [post]
func (cl *ClubController) CreateClub(c echo.Context) error {
	var input models.Club
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input: "+err.Error())
	}

	if err := validation.ValidatePayload(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, err := cl.clubService.CreateClub(&input)
	if err != nil {
		if errors.Is(err, errors.New(postgres.ErrClubNameAlreadyExists)) {
			return echo.NewHTTPError(http.StatusConflict, "club already exists")

		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]int{"id": id})
}

// GetAllClubs godoc
// @Summary Get all clubs
// @Description Get details of all clubs
// @Tags clubs
// @Produce  json
// @Success 200 {array} models.Club
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /clubs [get]
func (cl *ClubController) GetAllClubs(c echo.Context) error {
	clubs, err := cl.clubService.GetAllClubs()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, clubs)
}

// GetClubByID godoc
// @Summary Get a club by ID
// @Description Get detailed information about a specific club
// @Tags clubs
// @Produce  json
// @Param id path int true "Club ID"
// @Success 200 {object} models.Club
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /clubs/{id} [get]
func (cl *ClubController) GetClubByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	club, err := cl.clubService.GetClubByID(id)
	if err != nil {
		if errors.Is(err, errors.New(postgres.ErrClubNotFound)) {
			return echo.NewHTTPError(http.StatusNotFound, "club not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, club)
}

type ClubUpdate struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Contact     string  `json:"contact"`
	Price       float64 `json:"price" validate:"omitempty,twoDecimalPlaces"`
	SpotsNumber int     `json:"spots_number" validate:"omitempty,number,min=1"`
	IsActive    string  `json:"is_active"`
}

// UpdateClub godoc
// @Summary Update a club
// @Description Update a club with the provided club details
// @Tags clubs
// @Accept  json
// @Produce  json
// @Param id path int true "Club ID"
// @Param input body ClubUpdate true "Club details"
// @Success 200 {object} map[string]number
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /clubs/{id} [put]
func (cl *ClubController) UpdateClub(c echo.Context) error {
	var input ClubUpdate
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input: "+err.Error())
	}

	if err := validation.ValidatePayload(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	club := models.Club{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		Contact:     input.Contact,
		Price:       input.Price,
		SpotsNumber: input.SpotsNumber,
		IsActive:    input.IsActive,
	}

	err = cl.clubService.UpdateClub(&club)
	if err != nil {
		if errors.Is(err, errors.New(postgres.ErrClubNotFound)) {
			return echo.NewHTTPError(http.StatusNotFound, "club not found")
		} else if errors.Is(err, errors.New(postgres.ErrClubNameAlreadyExists)) {
			return echo.NewHTTPError(http.StatusConflict, "club already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// DeleteClub godoc
// @Summary Delete a club
// @Description Delete a club by its ID
// @Tags clubs
// @Produce  json
// @Param id path int true "Club ID"
// @Success 204 {object} map[string]number
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /clubs/{id} [delete]
func (cl *ClubController) DeleteClub(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	err = cl.clubService.DeleteClub(id)
	if err != nil {
		if errors.Is(err, errors.New(postgres.ErrClubNotFound)) {
			return echo.NewHTTPError(http.StatusNotFound, "club not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
