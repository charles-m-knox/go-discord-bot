package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	ch, ok := channels[m.ChannelID]
	if !ok {
		log.Printf("[%v] %v: %v", m.ChannelID, m.Author.Username, m.Content)
		return
	}

	log.Printf("[%v] %v#%v (%v): %v", ch, m.Author.Username, m.Author.Discriminator, m.Author.GlobalName, m.Content)

	// ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
		if err != nil {
			log.Printf("failed to send pong message: %v", err.Error())
		}
	}

	if m.Content == "pong" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Ping!")
		if err != nil {
			log.Printf("failed to send ping message: %v", err.Error())
		}
	}

	writeMsgLog(msgLog{
		Message: *m.Message,
		Time:    time.Now().Unix(),
	})
}

func writeMsgLog(m msgLog) {
	mb, err := json.Marshal(m)
	if err != nil {
		log.Printf("failed to marshal msg log event: %v", err.Error())
		return
	}

	mb = append(mb, '\n')

	writeOrAppendToFile(mb, msgLogBasePath, msgLogFileName)
}
