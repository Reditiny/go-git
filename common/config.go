package common

import "sync"

type Config struct {
	User  string
	Email string
}

var config Config
var once sync.Once

func ReadConfig() Config {
	once.Do(func() {
		// todo read config from file
		config = Config{
			User:  "default",
			Email: "default@red.com",
		}
	})
	return config
}
