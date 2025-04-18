package server

import (
	"log"
	"merch_service/new_version/configs"
	"merch_service/new_version/internal/handlers"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type MerchServer struct {
	router *gin.Engine
	config *configs.ServerConfig

	uHandler *handlers.UserHandler
	mHandler *handlers.MerchHandler
	tHandler *handlers.TransactionHandler
}

func (serv *MerchServer) LoadConfig(configPath string) {
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
		router: router,
	}

	newServ.LoadConfig("configs/server_config.yml")

	return &newServ
}

func (serv *MerchServer) SetupRoutes() {
	// REMOVE AFTER DEBUG !!!
	// serv.router.POST("/ping", func(c *gin.Context) {
	// 	val := c.GetHeader("Abba")
	// 	log.Println("DEBUG:", val, len(val))
	// 	c.String(200, "OK!")
	// })

	// --- Публичные пути START --- //
	serv.router.POST("/auth/register", serv.uHandler.RegHandler()) // DONE
	// serv.router.POST("/auth/login", LoginHandler(serv.config))     // DONE
	// --- Публичные пути END --- //

	// --- Приватные пути START --- //
	// AuthRequired применяется только к путям в данной группе
	authorized := serv.router.Group("/")
	// authorized.Use(AuthRequired(serv.config))
	{
		authorized.GET("/merch", serv.mHandler.MerchListHandler)
		authorized.POST("/merch/buy", serv.mHandler.BuyMerchHandler)
		authorized.GET("/history/coins", serv.uHandler.CoinsHistoryHandler)
		authorized.GET("/history/purchase")
		authorized.POST("/coins/transfer", serv.tHandler.TransferHandler)
	}

	// --- Приватные пути END --- //
}

func (serv MerchServer) Start() {
	// Здесь требуется инициализация базы данных
	serv.SetupRoutes()
	serv.router.Run(serv.config.Host + ":" + strconv.Itoa(serv.config.Port))
}
