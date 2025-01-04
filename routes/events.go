package routes

import (
	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func getEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event id"})
	}

	event, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "could not find event"})
	}

	context.JSON(http.StatusOK, gin.H{"event": event})
}

func createEvent(context *gin.Context) {

	var event models.Event
	err := context.ShouldBindBodyWithJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse body"})
		return
	}
	userId := context.GetInt64("userId")
	event.UserID = userId
	err = event.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save event"})
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Event created", "event": event})
}

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get events", "error": err})
		return
	}
	context.JSON(http.StatusOK, events)
}

func updateEvent(context *gin.Context) {
	// extract id from url path
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event id"})
	}
	// check if event is exist in DB
	userId := context.GetInt64("userId")
	event, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not find event"})
		return
	}
	// check if user id is matched with user id of the event
	// only the user can update its own event
	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to update event"})
		return
	}

	// extract request body type event
	var updatedEvent models.Event
	err = context.ShouldBindBodyWithJSON(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event id"})
		return
	}

	// update event in DB
	updatedEvent.ID = eventId
	err = updatedEvent.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update event", "error": err})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event updated", "event": updatedEvent})
}

func deleteEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event id"})
	}

	userId := context.GetInt64("userId")
	event, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not find event"})
		return
	}

	// check if user id is matched with user id of the event
	// only the user can delete its own event
	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to update event"})
		return
	}

	err = event.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete event", "error": err})
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event deleted"})
}
