package server

import (
	"log"
	"merch_service/internal/database"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type Server struct {
	router *gin.Engine
	config *ServerConfig
}

// требуемый API
// Метод 	Путь 	Описание
// POST 	/auth/register 	Регистрация нового сотрудника
// POST 	/auth/login 	Авторизация (получение JWT)
// GET 	    /merch 	Список товаров
// POST 	/merch/buy 	Покупка товара
// POST 	/coins/transfer 	Перевод монет другому сотруднику
// GET 	/history 	История операций пользователя

func loadConfig(configPath string) *ServerConfig {
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

	sc := ServerConfig{}

	yaml.Unmarshal(config, &sc)
	return &sc
}

func NewServer() Server {
	router := gin.Default()
	config := loadConfig("configs/server_config.yml")

	// TODO загрузить секретный ключ
	newServ := Server{
		router: router,
		config: config,
	}

	return newServ
}

func (serv Server) SetupRoutes() {
	// REMOVE AFTER DEBUG !!!
	// serv.router.POST("/ping", func(c *gin.Context) {
	// 	val := c.GetHeader("Abba")
	// 	log.Println("DEBUG:", val, len(val))
	// 	c.String(200, "OK!")
	// })

	// --- Публичные пути START --- //
	serv.router.POST("/auth/register", RegHandler(serv.config)) // DONE
	serv.router.POST("/auth/login", LoginHandler(serv.config))  // DONE
	// --- Публичные пути END --- //

	// --- Приватные пути START --- //
	// AuthRequired применяется только к путям в данной группе
	authorized := serv.router.Group("/")
	authorized.Use(AuthRequired(serv.config))
	{
		authorized.GET("/merch", MerchList)                 // DONE
		authorized.GET("/history", WalletHistory)           // NOT YET
		authorized.POST("/merch/buy")                       // NOT YET
		authorized.POST("/coins/transfer", TransferHandler) // HALF DONE
	}

	// --- Приватные пути END --- //
}

func (serv Server) Run() {
	database.InitDB()
	serv.SetupRoutes()
	serv.router.Run(serv.config.Host + ":" + strconv.Itoa(serv.config.Port))
}
