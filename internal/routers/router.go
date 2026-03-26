package routers

import (
	"notes-service/internal/api"
	"notes-service/internal/database"
	"notes-service/internal/service"
	"notes-service/internal/model"
	"errors"
	"fmt"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
		return func(c *gin.Context) {
				c.Next()
				if len(c.Errors) > 0 {
						err := c.Errors.Last().Err
						
						if apiErr, ok := errors.AsType[*model.ApiError](err); ok {
							status := http.StatusInternalServerError
							switch apiErr.Code {
								case "not-found":
									status = 404
								case "validation error":
									status = 422
							}

							c.JSON(status, apiErr)
							return
						} else {
							c.JSON(http.StatusInternalServerError, map[string]any{
								"error": "Internal Server Error",
							})
						}
				}
		}
}

func Router(db *database.GormDatabase) *gin.Engine {
	r := gin.New()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
		)
	}))

	r.Use(ErrorHandler())

	noteSrv := service.NewNoteService(db) 
	noteApi := api.NewUserAPI(noteSrv)

	notesGroup := r.Group("/notes")
	{
			notesGroup.GET("", noteApi.GetAllNotes)
			notesGroup.POST("", noteApi.CreateNote)
			notesGroup.GET("/:id", noteApi.GetNote)
	}
	
	healthySrv := service.NewHealthyService(db)
	healthyApi := api.NewHealthyAPI(healthySrv)

	healthyGroup := r.Group("/healthy")
	{
		healthyGroup.GET("/alive", healthyApi.Alive)
		healthyGroup.GET("/ready", healthyApi.Ready)
	}
	
	return r
}
