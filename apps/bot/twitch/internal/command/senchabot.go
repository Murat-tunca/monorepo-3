package command

import (
	"github.com/gempir/go-twitch-irc/v3"
	"github.com/senchabot-dev/monorepo/apps/bot/twitch/client"
	"github.com/senchabot-dev/monorepo/apps/bot/twitch/server"
)

func SenchabotCommand(client *client.Clients, server *server.SenchabotAPIServer, message twitch.PrivateMessage, commandName string, params []string) {
	client.Twitch.Say(message.Channel, "https://github.com/senchabot-dev/monorepo")
}
