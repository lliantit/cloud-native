package main

import (
	"fmt"
	"go-service/internal/controller"
	"go-service/internal/middleware"
	"go-service/internal/repository"
	"go-service/internal/service"
	"net/http"
)

func main() {
	// 1. Ініціалізація шарів
	repo := repository.NewInMemoryUserRepository()
	svc := service.NewUserService(repo)
	ctrl := controller.NewUserController(svc)

	// 2. Налаштування маршрутизатора
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/api/users", ctrl.HandleUsers)
	mux.HandleFunc("/api/users/", ctrl.HandleUserByID)

	// 3. Збірка мідлварів (Validation -> Logging -> Router)
	var handler http.Handler = mux
	handler = middleware.ContentTypeValidation(handler)
	handler = middleware.Logging(handler)

	// 4. Запуск сервера
	fmt.Println("Go server is running on port 8088")
	if err := http.ListenAndServe(":8088", handler); err != nil {
		fmt.Printf("Server error: %s\n", err)
	}
}
