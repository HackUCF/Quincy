package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/HackUCF/Quincy/common/log"
	_ "github.com/joho/godotenv/autoload"
)

const num_agents int = 25

func main() {
	r := getRunner()

	for range num_agents {
		id, err := uuid.NewRandom()
		if err != nil {
			log.Panic("uuid failed to generate. this should never happen")
		}

		go r.loop(id.String())
	}

	for {
		time.Sleep(time.Millisecond)
	}
}
