package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type ParserFunc func(*discordgo.Session, *discordgo.MessageCreate) (bool, error)

type ParserChain struct {
	parsers []ParserFunc
}

func NewParserChain() *ParserChain {
	return &ParserChain{
		parsers: make([]ParserFunc, 0),
	}
}

func (p *ParserChain) AddParser(parser ...ParserFunc) *ParserChain {
	p.parsers = append(p.parsers, parser...)
	return p
}

func (p *ParserChain) Parse(discord *discordgo.Session, message *discordgo.MessageCreate) {
	for _, parse := range p.parsers {
		handled, err := parse(discord, message)
		if err != nil {
			fmt.Println(err)
		}
		if handled {
			return
		}
	}
}
