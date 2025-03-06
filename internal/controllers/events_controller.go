package controllers

import (
	"strconv"

	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/dtos"
	"github.com/OvniCore-SA/api_go_whatsapp_chatbot/internal/services"
	"github.com/gofiber/fiber/v2"
)

// EventsController maneja las rutas relacionadas con los eventos
type EventsController struct {
	eventsService services.EventsService
}

// NewEventsController crea una nueva instancia del controlador
func NewEventsController(eventsService services.EventsService) *EventsController {
	return &EventsController{eventsService: eventsService}
}

// Crear un evento
func (ec *EventsController) CreateEvent(c *fiber.Ctx) error {
	var eventDTO dtos.EventsDto
	if err := c.BodyParser(&eventDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := ec.eventsService.Create(eventDTO); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Event created successfully"})
}

// Obtener un evento por ID
func (ec *EventsController) GetEventByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
	}

	event, err := ec.eventsService.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
	}

	return c.JSON(event)
}

// Obtener todos los eventos
func (ec *EventsController) GetAllEvents(c *fiber.Ctx) error {
	events, err := ec.eventsService.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving events"})
	}

	return c.JSON(events)
}

// Actualizar un evento
func (ec *EventsController) UpdateEvent(c *fiber.Ctx) error {
	var eventDTO dtos.EventsDto
	if err := c.BodyParser(&eventDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := ec.eventsService.Update(eventDTO); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Event updated successfully"})
}

// Eliminar un evento por ID
func (ec *EventsController) DeleteEvent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
	}

	if err := ec.eventsService.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Event deleted successfully"})
}

// Cancelar un evento por código
func (ec *EventsController) CancelEvent(c *fiber.Ctx) error {
	codeEvent := c.Params("codeEvent")

	if err := ec.eventsService.Cancel(codeEvent); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Event canceled successfully"})
}

// Obtener eventos por contacto y fecha
func (ec *EventsController) GetEventsByContactAndDate(c *fiber.Ctx) error {
	contactID, err := strconv.ParseInt(c.Params("contactID"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid contact ID"})
	}

	date := c.Params("date")
	currentTime := c.Params("currentTime")

	events, err := ec.eventsService.GetEventByContactAndDate(contactID, date, currentTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(events)
}
