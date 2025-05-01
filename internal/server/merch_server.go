package server

import (
	"context"
	"fmt"
	"log"
	"merch_service/configs"
	"merch_service/internal/handlers"
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
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		log.Fatalln("не удалось определить абсолютный путь:", err.Error())
	}

	fmt.Println(absPath)
	config, err := os.ReadFile(absPath)
	if err != nil {
		log.Fatalln("ошибка при чтении конфигурационного файла:", err.Error())
	}

	sc := configs.ServerConfig{}
	if err := yaml.Unmarshal(config, &sc); err != nil {
		log.Fatalln("ошибка при разборе yaml:", err.Error())
	}

	serv.config = &sc
}

func NewMerchServer(u *handlers.UserHandler, t *handlers.TransactionHandler, m *handlers.MerchHandler, configPath string) *MerchServer {
	router := gin.Default()

	newServ := MerchServer{
		http: &http.Server{
			Handler: router,
		},
		uHandler: u,
		tHandler: t,
		mHandler: m,
	}

	// Хардоженые пути, сорян =(
	if configPath == "" {
		configPath = "configs/server_config.yml"
	}

	newServ.loadConfig(configPath)

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
