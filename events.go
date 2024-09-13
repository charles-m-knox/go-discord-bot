package main

import (
	"encoding/json"
	"log"

	"github.com/bwmarrin/discordgo"
)

func writeScheduledEvents(basePath, fileName string, events []*discordgo.GuildScheduledEvent) {
	b, err := json.Marshal(events)
	if err != nil {
		log.Printf("failed to marshal user add event: %v", err.Error())
		return
	}

	writeToFile(b, basePath, fileName)
}
