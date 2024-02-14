package config

type app struct {
	Desc    string `yaml:"desc"`
	Address string `yaml:"address"`
	Version string `yaml:"version"`
	Env     string `yaml:"env"`
}

type database struct {
	Kind     string `yaml:"kind"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type redis struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type ServerConfig struct {
	App      app      `yaml:"app"`
	Database database `yaml:"database"`
	Redis    redis    `yaml:"redis"`
}
