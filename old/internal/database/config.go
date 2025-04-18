package database

type DBConfig struct {
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	Addr   string `yaml:"addr"`
	Port   int    `yaml:"port"`
	DBName string `yaml:"name"`
}
