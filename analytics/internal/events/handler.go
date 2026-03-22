package events

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/26thavenue/backend_testproj/analytics/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Track(c fiber.Ctx) error {
	var req domain.TrackRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "invalid JSON body")
	}
	req.EventName = strings.TrimSpace(req.EventName)

	if req.EventName == "" {
		return writeError(c, fiber.StatusBadRequest, "event_name is required")
	}

	if err := h.service.Track(c.Context(), req); err != nil {
		return writeError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(map[string]any{
		"status": "ok",
	})
}

func (h *Handler) List(c fiber.Ctx) error {
	filter := domain.Filter{
		EventName: strings.TrimSpace(c.Query("event")),
		UserID:    strings.TrimSpace(c.Query("user_id")),
	}

	from, err := parseTime(c.Query("from"))
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "invalid from timestamp")
	}
	to, err := parseTime(c.Query("to"))
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "invalid to timestamp")
	}
	filter.From = from
	filter.To = to

	events, err := h.service.GetRawEvents(c.Context(), filter)
	if err != nil {
		return writeError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(events)
}

func (h *Handler) Aggregate(c fiber.Ctx) error {
	eventName := strings.TrimSpace(c.Query("event"))
	if eventName == "" {
		return writeError(c, fiber.StatusBadRequest, "event is required")
	}

	from, err := parseTime(c.Query("from"))
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "invalid from timestamp")
	}
	to, err := parseTime(c.Query("to"))
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "invalid to timestamp")
	}

	result, err := h.service.GetAnalytics(c.Context(), domain.AggregateQuery{
		EventName: eventName,
		From:      from,
		To:        to,
	})
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(result)
}

func (h *Handler) CreateEventType(c fiber.Ctx) error {
	var req domain.EventType
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "invalid JSON body")
	}
	req.Name = strings.TrimSpace(req.Name)

	if err := h.service.CreateEventType(c.Context(), req); err != nil {
		return writeError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(map[string]any{
		"status": "ok",
	})
}

func parseTime(value string) (time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, value)
}

func writeError(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(map[string]any{
		"error": message,
	})
}
