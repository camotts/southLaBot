package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	commandPrefix string
	botID         string
	apiKey        string
)

func init() {
	flag.StringVar(&apiKey, "api-key", "", "API key for the discord bot")
}

func main() {
	flag.Parse()
	discord, err := discordgo.New(fmt.Sprintf("Bot %s", apiKey))
	errCheck("error creating discord session", err)
	user, err := discord.User("@me")
	errCheck("error retrieving account", err)

	botID = user.ID

	rp := RoleParser{}

	parser := NewParserChain().AddParser(func(discord *discordgo.Session, message *discordgo.MessageCreate) (bool, error) {
		return IsBot(message.Author), nil
	}).AddParser(rp.ListRole, rp.AddRole, rp.RemoveRole, rp.ConfigureRole, rp.DeleteRole)

	discord.AddHandler(func(discord *discordgo.Session, message *discordgo.MessageCreate) {
		parser.Parse(discord, message)
	})

	err = discord.Open()
	errCheck("Error opening connection to Discord", err)
	defer discord.Close()

	LoadRoles()
	defer SaveRoles()

	commandPrefix = "!"
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		fmt.Println("Got terminate signal")
	}

}

func errCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}
