package config

import "time"

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
	Host     string        `yaml:"host"`
	Port     string        `yaml:"port"`
	Password string        `yaml:"password"`
	Username string        `yaml:"user"`
	TTL      time.Duration `yaml:"ttl"`
}

type schedule struct {
	FirstInterval time.Duration `yaml:"firstInterval"`
	Interval      time.Duration `yaml:"interval"`
	Period        time.Duration `yaml:"period"`
}

type localCache struct {
	MaxSize  int           `yaml:"maxSize"`
	Interval time.Duration `yaml:"interval"`
	Period   time.Duration `yaml:"period"`
	TTL      time.Duration `yaml:"ttl"`
}

type hotdata struct {
	AgeStart int `yaml:"ageStart"`
	AgeEnd   int `yaml:"ageEnd"`
}

type ServerConfig struct {
	App        app        `yaml:"app"`
	Database   database   `yaml:"database"`
	Redis      redis      `yaml:"redis"`
	Schedule   schedule   `yaml:"schedule"`
	LocalCache localCache `yaml:"local"`
	HotData    hotdata    `yaml:"hotdata"`
}
