package server

type Server interface {
	LoadConfig(configPath string)
	Start()
}
