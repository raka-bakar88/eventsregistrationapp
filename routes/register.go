package routes

import (
	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("user_id")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event id"})
	}

	event, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not find event"})
		return
	}

	err = event.Register(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could register for an event"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Registered!"})
}

func cancelRegistration(context *gin.Context) {
	userId := context.GetInt64("user_id")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event id"})
	}

	event, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not find event"})
		return
	}

	err = event.Unregister(userId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could unregister for an event"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Unregistered successfully!"})
}
