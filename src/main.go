package main

import (
	"github.com/HackUCF/Quincy/cmd"

	// automatically load .env files
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cmd.Execute()
}
