package alert

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/gadost/telescope/conf"
)

func (a *Alert) DiscordSend() {
	dg, err := discordgo.New("Bot " + conf.MainConfig.Discord.Token)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return
	}

	err = dg.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}
	dg.ChannelMessageSend(
		fmt.Sprint(conf.MainConfig.Discord.ChannelID),
		"***"+a.Importance+"***"+": \n"+a.Message,
	)

	dg.Close()
}

func DiscordSendTest(t string, c string) error {
	dg, err := discordgo.New("Bot " + t)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return err
	}

	err = dg.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return err
	}
	dg.ChannelMessageSend(fmt.Sprint(c), "***Pong***")

	dg.Close()

	return err

}
