package server

import (
	"context"
	"log"
	"merch_service/new_version/configs"
	"merch_service/new_version/internal/handlers"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

// MerchServer - структура сервера, имплементация Server
// Содержит хендлеры
//   - UserHandler
//   - MerchHandler
//   - TransactionHandler
//
// для обработки соответстующих API запросов
type MerchServer struct {
	http   *http.Server
	config *configs.ServerConfig

	uHandler *handlers.UserHandler
	mHandler *handlers.MerchHandler
	tHandler *handlers.TransactionHandler
}

func (serv *MerchServer) loadConfig(configPath string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("не удалось получить путь до рабочей директорий:", err.Error())
	}

	// Поднимаемся в корень директорий и оттуда смотрим на configPath
	// TODO выглядит криво, что-то придумать!
	configPath = filepath.Join(filepath.Dir(filepath.Dir(wd)), configPath)
	config, err := os.ReadFile(configPath)

	if err != nil {
		log.Fatalln("ошибка при чтений конфигурационного файла:", err.Error())
	}

	sc := configs.ServerConfig{}

	yaml.Unmarshal(config, &sc)

	serv.config = &sc
}

func NewMerchServer() *MerchServer {
	router := gin.Default()

	newServ := MerchServer{
		http: &http.Server{
			Handler: router,
		},
	}

	newServ.loadConfig("configs/server_config.yml")

	return &newServ
}

func (serv *MerchServer) SetupRoutes() {
	router := serv.http.Handler.(*gin.Engine)
	// --- Публичные пути START --- //
	router.POST("/auth/register", serv.uHandler.RegHandler())
	router.POST("/auth/login", serv.uHandler.LoginHandler(serv.config))
	// --- Публичные пути END --- //

	// --- Приватные пути START --- //
	// AuthRequired применяется только к путям в данной группе
	authorized := router.Group("/")
	authorized.Use(handlers.AuthRequired(serv.config))
	{
		authorized.GET("/merch", serv.mHandler.MerchListHandler)
		authorized.POST("/merch/buy", serv.mHandler.BuyMerchHandler)
		authorized.GET("/history/coins", serv.uHandler.CoinsHistoryHandler)
		authorized.GET("/history/purchase", serv.uHandler.PurchaseHistoryHandler)
		authorized.POST("/coins/transfer", serv.tHandler.TransferHandler)
	}

	// --- Приватные пути END --- //
}

// Start - настраивает пути API и запускает сервер по адресу в serverConfig
func (serv *MerchServer) Start() {
	serv.SetupRoutes()

	// graceful stop (см. https://gin-gonic.com/en/docs/examples/graceful-restart-or-stop/)
	serv.http.Addr = serv.config.Host + ":" + strconv.Itoa(serv.config.Port)

	if err := serv.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

// Stop - остановка сервера, только для тестирования
// Для правильного использования нужно выполнить Start в отдельной горутине
func (serv *MerchServer) Stop() {
	log.Println("Выключение сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := serv.http.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Println("Таймаут в 5с.")
	log.Println("Сервер выключен")
}
