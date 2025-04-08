package server

import (
	"log"
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
	config := loadConfig("configs/config.yml")

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
	serv.router.POST("/auth/register", RegHandler(serv.config)) // TODO handler
	serv.router.POST("/auth/login", LoginHandler(serv.config))  // TODO handler
	// --- Публичные пути END --- //

	// --- Приватные пути START --- //
	// AuthRequired применяется только к путям в данной группе
	authorized := serv.router.Group("/")
	authorized.Use(AuthRequired(serv.config))
	{
		authorized.GET("/merch", MerchList)
		authorized.GET("/history", WalletHistory)
		authorized.POST("/merch/buy")      // TODO handler
		authorized.POST("/coins/transfer") // TODO handler
	}
	// --- Приватные пути END --- //
}

func (serv Server) Run() {
	serv.SetupRoutes()
	serv.router.Run(serv.config.Host + ":" + strconv.Itoa(serv.config.Port))
}
