package server

type ServerConfig struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	Secret        string `yaml:"secret"` // Секретный ключ. Для простоты храним секретный ключ в структуре;
	RefreshSecret string `yaml:"refresh"`
	ExpTimeout    int64  `yaml:"exptimeout"` // Время жизни токена. Задается в секундах
}
