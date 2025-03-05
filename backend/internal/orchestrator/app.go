package orchestrator

import (
	"calc-website/config"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func Run(cfg *config.Config) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	service := NewAPIService(cfg)
	apiHandler := NewAPIHandler(service)
	router := apiHandler.Router()
	// Настраиваем CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Разрешить запросы с любого источника
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// Запускаем сервер с CORS
	handler := c.Handler(router)
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		return err
	}
	return nil
}
