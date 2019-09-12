package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Parser interface {
	Parse(*discordgo.Session, *discordgo.MessageCreate) bool
}

type ParserChain struct {
	parsers []Parser
}

func parseCommand(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == botID || user.Bot {
		return
	}

	content := message.Content

	if content == "!test" {
		discord.ChannelMessageSend(message.ChannelID, "Testing..")
	}

	if strings.HasPrefix(content, "!addRole") {
		fmt.Println("Adding role")
		roles := strings.Split(strings.ToLower(content), " ")[1:]
		currentRoles, _ := discord.GuildRoles(message.GuildID)
		serverRoles := make(map[string]string)
		for _, cr := range currentRoles {
			serverRoles[strings.ToLower(cr.Name)] = cr.ID
		}
		for _, role := range roles {
			fmt.Println(role)
			if id, exists := serverRoles[role]; exists {
				if ok := Roles[id]; ok {
					err := discord.GuildMemberRoleAdd(message.GuildID, message.Author.ID, id)
					if err != nil {
						fmt.Println(err)
					} else {
						discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Successfully added %s to %s", user.Mention(), role))
					}
				}
			}
		}
	}

	if strings.HasPrefix(content, "!listRole") {
		fmt.Println("Listing role")
		currentRoles, _ := discord.GuildRoles(message.GuildID)
		serverRoles := make(map[string]string)
		for _, cr := range currentRoles {
			serverRoles[strings.ToLower(cr.Name)] = cr.ID
		}
		rolesToSend := make([]string, 0)
		for role, id := range serverRoles {
			fmt.Println(role)
			if ok := Roles[id]; ok {
				rolesToSend = append(rolesToSend, role)
			}
		}
		discord.ChannelMessageSend(message.ChannelID, strings.Join(rolesToSend, ", "))
	}

	if strings.HasPrefix(content, "!removeRole") {
		fmt.Println("Adding role")
		roles := strings.Split(strings.ToLower(content), " ")[1:]
		currentRoles, _ := discord.GuildRoles(message.GuildID)
		serverRoles := make(map[string]string)
		for _, cr := range currentRoles {
			serverRoles[strings.ToLower(cr.Name)] = cr.ID
		}
		for _, role := range roles {
			fmt.Println(role)
			if id, exists := serverRoles[role]; exists {
				if ok := Roles[id]; ok {
					err := discord.GuildMemberRoleRemove(message.GuildID, message.Author.ID, id)
					if err != nil {
						fmt.Println(err)
					} else {
						discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Successfully removed %s from %s", user.Mention(), role))
					}
				}
			}
		}
	}

	if strings.HasPrefix(content, "!configureRole") {
		mem, _ := discord.GuildMember(message.GuildID, message.Author.ID)
		for _, roleId := range mem.Roles {
			role, err := discord.State.Role(message.GuildID, roleId)
			if err != nil {
				return
			}
			if role.Permissions&discordgo.PermissionAdministrator != 0 {
				roles := strings.Split(strings.ToLower(content), " ")[1:]
				currentRoles, _ := discord.GuildRoles(message.GuildID)
				serverRoles := make(map[string]string)
				fmt.Println("Adding: ", roles)
				for _, cr := range currentRoles {
					serverRoles[strings.ToLower(cr.Name)] = cr.ID
				}
				for _, role := range roles {
					if id, exists := serverRoles[role]; exists {
						Roles[id] = true
					}
				}
			}
		}
	}

	if strings.HasPrefix(content, "!deleteRole") {
		mem, _ := discord.GuildMember(message.GuildID, message.Author.ID)
		for _, roleId := range mem.Roles {
			role, err := discord.State.Role(message.GuildID, roleId)
			if err != nil {
				return
			}
			if role.Permissions&discordgo.PermissionAdministrator != 0 {
				roles := strings.Split(strings.ToLower(content), " ")[1:]
				currentRoles, _ := discord.GuildRoles(message.GuildID)
				serverRoles := make(map[string]string)
				fmt.Println("Deleting: ", roles)
				for _, cr := range currentRoles {
					serverRoles[strings.ToLower(cr.Name)] = cr.ID
				}
				for _, role := range roles {
					if id, exists := serverRoles[role]; exists {
						Roles[id] = false
					}
				}
			}
		}
	}
}
