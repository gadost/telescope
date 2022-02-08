package alert

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/gadost/telescope/conf"
)

func DiscordSend(s string, m string) {
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
	dg.ChannelMessageSend(fmt.Sprint(conf.MainConfig.Discord.ChannelID), "***"+s+"***"+": \n"+m)

	dg.Close()

}
