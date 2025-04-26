package api

import (
	"db/api/handler"
	"db/ent"
	"db/service"

	"github.com/gin-gonic/gin"
)

// SetupRouter настраивает маршрутизацию Gin и возвращает экземпляр роутера
func SetupRouter(client *ent.Client) *gin.Engine {
	r := gin.Default()

	// Инициализация сервисов и обработчиков
	userService := service.NewUserService(client)
	userHandler := handler.NewUserHandler(userService)

	// Группа API endpoints
	api := r.Group("/api")
	{
		// Маршруты для пользователей
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
		}
	}

	return r
}
