package main

import (
	"fmt"

	"github.com/pistarlabs/configs"
)

func main() {
	// Initialize config
	cfg, err := configs.Load("/path/toconfig.json")
	if err != nil {
		panic(err)
	}

	// Get development environment configuration
	cfg, err = cfg.Get("development")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Database host is %s", cfg.UString("database.host"))
}
