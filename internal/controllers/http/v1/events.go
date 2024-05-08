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

func (e *EventController) GetAllEvents(c echo.Context) error {
	events, err := e.eventService.GetAllEvents()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, events)
}

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

func (e *EventController) UpdateEvent(c echo.Context) error {
	type EventUpdate struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Date        string  `json:"date" validate:"omitempty,datetime=2006-01-02 15:04"`
		Price       float64 `json:"price" validate:"omitempty,twoDecimalPlaces"`
	}

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
