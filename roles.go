package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type RoleParser struct {
}

var Roles map[string]bool

func LoadRoles() {
	Roles = make(map[string]bool)
	raw, _ := ioutil.ReadFile("roles.txt")
	for _, role := range strings.Split(string(raw), "\n") {
		Roles[role] = true
	}
	fmt.Println(Roles)
}

func SaveRoles() {
	data := ""
	arr := make([]string, 0)
	for role := range Roles {
		arr = append(arr, role)
	}
	data = strings.Join(arr, "\n")
	ioutil.WriteFile("roles.txt", []byte(data), 0777)
}

func (rp *RoleParser) ListRole(discord *discordgo.Session, message *discordgo.MessageCreate) (handled bool, err error) {
	if strings.HasPrefix(message.Content, "!listRole") {
		handled = true
		currentRoles, err := discord.GuildRoles(message.GuildID)
		if err != nil {
			return handled, err
		}
		serverRoles := make(map[string]string)
		for _, cr := range currentRoles {
			serverRoles[strings.ToLower(cr.Name)] = cr.ID
		}
		rolesToSend := make([]string, 0)
		for role, id := range serverRoles {
			if ok := Roles[id]; ok {
				rolesToSend = append(rolesToSend, role)
			}
		}
		discord.ChannelMessageSend(message.ChannelID, strings.Join(rolesToSend, ", "))
	}
	return
}

func (rp *RoleParser) AddRole(discord *discordgo.Session, message *discordgo.MessageCreate) (handled bool, err error) {
	if strings.HasPrefix(message.Content, "!addRole") {
		handled = true
		roles := strings.Split(strings.ToLower(message.Content), " ")[1:]
		currentRoles, _ := discord.GuildRoles(message.GuildID)
		serverRoles := make(map[string]string)
		for _, cr := range currentRoles {
			serverRoles[strings.ToLower(cr.Name)] = cr.ID
		}
		for _, role := range roles {
			if id, exists := serverRoles[role]; exists {
				if ok := Roles[id]; ok {
					err = discord.GuildMemberRoleAdd(message.GuildID, message.Author.ID, id)
					if err == nil {
						discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Successfully added %s to %s", message.Author.Mention(), role))
					}
				}
			}
		}
	}
	return
}

func (rp *RoleParser) RemoveRole(discord *discordgo.Session, message *discordgo.MessageCreate) (handled bool, err error) {
	if strings.HasPrefix(message.Content, "!removeRole") {
		handled = true
		roles := strings.Split(strings.ToLower(message.Content), " ")[1:]
		currentRoles, err := discord.GuildRoles(message.GuildID)
		if err != nil {
			return handled, err
		}
		serverRoles := make(map[string]string)
		for _, cr := range currentRoles {
			serverRoles[strings.ToLower(cr.Name)] = cr.ID
		}
		for _, role := range roles {
			fmt.Println(role)
			if id, exists := serverRoles[role]; exists {
				if ok := Roles[id]; ok {
					err = discord.GuildMemberRoleRemove(message.GuildID, message.Author.ID, id)
					if err == nil {
						discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Successfully removed %s from %s", message.Author.Mention(), role))
					}
				}
			}
		}
	}
	return
}

func (rp *RoleParser) ConfigureRole(discord *discordgo.Session, message *discordgo.MessageCreate) (handled bool, err error) {
	if strings.HasPrefix(message.Content, "!configureRole") {
		handled = true
		mem, err := discord.GuildMember(message.GuildID, message.Author.ID)
		if err != nil {
			return handled, err
		}
		for _, roleId := range mem.Roles {
			role, err := discord.State.Role(message.GuildID, roleId)
			if err != nil {
				return handled, err
			}
			if role.Permissions&discordgo.PermissionAdministrator != 0 {
				roles := strings.Split(strings.ToLower(message.Content), " ")[1:]
				currentRoles, err := discord.GuildRoles(message.GuildID)
				if err != nil {
					return handled, err
				}
				serverRoles := make(map[string]string)
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
	return
}

func (rp *RoleParser) DeleteRole(discord *discordgo.Session, message *discordgo.MessageCreate) (handled bool, err error) {
	if strings.HasPrefix(message.Content, "!deleteRole") {
		handled = true
		mem, err := discord.GuildMember(message.GuildID, message.Author.ID)
		if err != nil {
			return handled, err
		}
		for _, roleId := range mem.Roles {
			role, err := discord.State.Role(message.GuildID, roleId)
			if err != nil {
				return handled, err
			}
			if role.Permissions&discordgo.PermissionAdministrator != 0 {
				roles := strings.Split(strings.ToLower(message.Content), " ")[1:]
				currentRoles, err := discord.GuildRoles(message.GuildID)
				if err != nil {
					return handled, err
				}
				serverRoles := make(map[string]string)
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
	return
}
