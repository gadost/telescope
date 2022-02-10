package status

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/gadost/telescope/conf"
)

func DiscordHandler() error {
	dg, err := discordgo.New("Bot " + conf.MainConfig.Discord.Token)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return err
	}
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return err
	}

	// Cleanly close down the Discord session.
	defer dg.Close()
	return err
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "$status" {
		s.ChannelMessageSend(m.ChannelID, StatusCollection())
	}
}
