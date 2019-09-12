package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func RoleParser(discord *discordgo.Session, message *discordgo.MessageCreate) {

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
