package main

import (
	"fmt"
	"sync"
)

type Config struct {
	AppName string
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{
			AppName: "My App",
		}
	})
	return instance
}

func main() {
	c1 := GetConfig()
	c2 := GetConfig()

	fmt.Println(c1 == c2) // true
}
