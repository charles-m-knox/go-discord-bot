package main

import (
	"github.com/bwmarrin/discordgo"
)

type presenceChange struct {
	// Typically "join" or "leave"
	Presence string `json:"presence"`
	// Time of the presence change event
	Time int64 `json:"time"`
	// The User that performed the presence change event
	User *discordgo.User `json:"user"`
	// The ID of the guild that the presence change event
	GuildID string `json:"guildId"`
}

type msgLog struct {
	Message discordgo.Message `json:"msg"`
	Time    int64             `json:"time"`
}
