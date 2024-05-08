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

type EventController struct {
	eventService service.Event
}

func NewEventRoutes(g *echo.Group, service service.Service) {
	eventController := &EventController{
		eventService: service.Event,
	}

	authMiddleware := &AuthMiddleware{
		authService: service.Auth,
	}

	g.GET("/events", eventController.GetAllEvents)
	g.GET("/events/:id", eventController.GetEventByID)

	withMiddleware := g.Group("")
	withMiddleware.Use(authMiddleware.IsAdmin)

	withMiddleware.POST("/events", eventController.CreateEvent)
	withMiddleware.PUT("/events/:id", eventController.UpdateEvent)
	withMiddleware.DELETE("/events/:id", eventController.DeleteEvent)
}

// CreateEvent godoc
// @Summary Create a new event
// @Description Create a new event with the provided event details
// @Tags events
// @Accept  json
// @Produce  json
// @Param input body models.Event true "Event details"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /events [post]
func (e *EventController) CreateEvent(c echo.Context) error {
	var input models.Event
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input: "+err.Error())
	}

	if err := validation.ValidatePayload(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, err := e.eventService.CreateEvent(&input)
	if err != nil {
		if errors.Is(err, errors.New(postgres.ErrEventNameAlreadyExists)) {
			return echo.NewHTTPError(http.StatusConflict, "event already exists")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]uint{"id": id})
}

// GetAllEvents godoc
// @Summary Get all events
// @Description Get details of all events
// @Tags events
// @Produce  json
// @Success 200 {array} models.Event
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /events [get]
func (e *EventController) GetAllEvents(c echo.Context) error {
	events, err := e.eventService.GetAllEvents()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, events)
}

// GetEventByID godoc
// @Summary Get an event by ID
// @Description Get detailed information about a specific event
// @Tags events
// @Produce  json
// @Param id path int true "Event ID"
// @Success 200 {object} models.Event
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /events/{id} [get]
func (e *EventController) GetEventByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	event, err := e.eventService.GetEventByID(uint(id))
	if err != nil {
		if errors.Is(err, errors.New(postgres.ErrEventNotFound)) {
			return echo.NewHTTPError(http.StatusNotFound, "event not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, event)
}

type EventUpdate struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Date        string  `json:"date" validate:"omitempty,datetime=2006-01-02 15:04"`
	Price       float64 `json:"price" validate:"omitempty,twoDecimalPlaces"`
}

// UpdateEvent godoc
// @Summary Update an event
// @Description Update an event with the provided event details
// @Tags events
// @Accept  json
// @Produce  json
// @Param id path int true "Event ID"
// @Param input body EventUpdate true "Event details"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /events/{id} [put]
func (e *EventController) UpdateEvent(c echo.Context) error {
	var input EventUpdate
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input: "+err.Error())
	}

	if err := validation.ValidatePayload(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	event := models.Event{
		Name:        input.Name,
		Description: input.Description,
		Date:        input.Date,
		Price:       input.Price,
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	event.ID = id
	err = e.eventService.UpdateEvent(&event)
	if err != nil {
		if errors.Is(err, errors.New(postgres.ErrEventNotFound)) {
			return echo.NewHTTPError(http.StatusNotFound, "event not found")
		} else if errors.Is(err, errors.New(postgres.ErrEventNameAlreadyExists)) {
			return echo.NewHTTPError(http.StatusConflict, "event already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// DeleteEvent godoc
// @Summary Delete an event
// @Description Delete an event by its ID
// @Tags events
// @Produce  json
// @Param id path int true "Event ID"
// @Success 204 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /events/{id} [delete]
func (e *EventController) DeleteEvent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	err = e.eventService.DeleteEvent(uint(id))
	if err != nil {
		if errors.Is(err, errors.New(postgres.ErrEventNotFound)) {
			return echo.NewHTTPError(http.StatusNotFound, "event not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
