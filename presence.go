package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func writePresenceChange(pc presenceChange) {
	pcb, err := json.Marshal(pc)
	if err != nil {
		log.Printf("failed to marshal user add event: %v", err.Error())
		return
	}

	pcb = append(pcb, '\n')

	writeOrAppendToFile(pcb, presenceChangesBasePath, presenceChangesFileName)
}

func memberJoin(s *discordgo.Session, a *discordgo.GuildMemberAdd) {
	log.Printf("user add: %v %v", a.User.Username, a.User.Discriminator)

	pc := presenceChange{
		Presence: "join",
		Time:     time.Now().Unix(),
		User:     a.User,
		GuildID:  a.GuildID,
	}

	writePresenceChange(pc)
}

func memberLeave(s *discordgo.Session, a *discordgo.GuildMemberRemove) {
	log.Printf("user del: %v %v", a.User.Username, a.User.Discriminator)

	pc := presenceChange{
		Presence: "leave",
		Time:     time.Now().Unix(),
		User:     a.User,
		GuildID:  a.GuildID,
	}

	writePresenceChange(pc)
}
