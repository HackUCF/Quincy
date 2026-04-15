package main

import (
	"github.com/HackUCF/quincy/cmd"

	// automatically load .env files
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cmd.Execute()
}
