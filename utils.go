package main

import "github.com/bwmarrin/discordgo"

func IsBot(user *discordgo.User) bool {
	return user.Bot || user.ID == botID
}
