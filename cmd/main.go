package main

import (
	"fmt"
	"golang/configs"
	"golang/internal/auth"
	"golang/internal/link"
	"golang/internal/stat"
	"golang/internal/user"
	"golang/packages/db"
	"golang/packages/event"
	"golang/packages/middleware"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDB(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	// Repositories
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepository := stat.NewStatRepository(db)

	// Services
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepository,
	})

	// handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		Config:         conf,
		EventBus:       eventBus,
	})
	stat.NewStatHandler(router, &stat.StatHandlerDeps{
		StatRepository: statRepository,
		Config:         conf,
	})

	// stack - это массив middleware, которые запускаются поочередно.
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging)

	// Здесь прописываются настройки сервера, который будет запущен в работу. Собственно,
	// Addr - это порт. Handler - функция, которая будет запускаться (запуск идет через router).
	// Каждый middlerware работают по принципу ПЕРЕХВАТИЛ ЗАПРОС/ОТВЕТ ОТ ВНЕШНЕГО СЕРВЕРА -> ДОПОЛНИЛ СВОЕЙ ИНФОЙ ->
	// -> ОТПРАВИЛ ЕГО ДАЛЬШЕ.
	server := http.Server{
		Addr:    "",
		Handler: stack(router),
	}

	go statService.AddClick()

	// ListenAndServe по факту представляет из себя горутину, которая нон-стопом слушает запросы на себя.
	fmt.Println("Server is working")
	server.ListenAndServe()
}
