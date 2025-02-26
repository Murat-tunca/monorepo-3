package command

import (
	"context"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/senchabot-opensource/monorepo/apps/discord-bot/internal/service"
	"github.com/senchabot-opensource/monorepo/apps/discord-bot/internal/service/streamer"
	"github.com/senchabot-opensource/monorepo/config"
	"github.com/senchabot-opensource/monorepo/packages/gosenchabot"
)

func (c *commands) DelTwitchCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, service service.Service) {
	options := i.ApplicationCommandData().Options

	switch options[0].Name {
	// del-twitch streamer
	case "streamer":
		options = options[0].Options
		twitchUsername := options[0].StringValue()
		twitchUsername = gosenchabot.ParseTwitchUsernameURLParam(twitchUsername)

		response0, uInfo := streamer.GetTwitchUserInfo(twitchUsername, c.twitchAccessToken)
		if response0 != "" {
			ephemeralRespond(s, i, response0)
			return
		}

		ok, err := service.DeleteDiscordTwitchLiveAnno(ctx, uInfo.ID, i.GuildID)
		if err != nil {
			ephemeralRespond(s, i, config.ErrorMessage+"#XXXX")
			return
		}

		if !ok {
			ephemeralRespond(s, i, "`"+twitchUsername+"` kullanıcı adlı Twitch yayıncısı veritabanında bulunamadı.")
			return
		}

		streamers := streamer.GetStreamersData(i.GuildID)
		delete(streamers, uInfo.Login)
		ephemeralRespond(s, i, "`"+uInfo.Login+"` kullanıcı adlı Twitch streamer veritabanından silindi.")

		// del-twitch event-channel
	case "event-channel":
		options = options[0].Options
		channelId := options[0].ChannelValue(s).ID
		channelName := options[0].ChannelValue(s).Name

		ok, err := service.DeleteAnnouncementChannel(ctx, channelId)
		if err != nil {
			ephemeralRespond(s, i, config.ErrorMessage+"#XXYX")
			log.Println("Error while deleting announcement channel:", err)
			return
		}
		if !ok {
			ephemeralRespond(s, i, "`"+channelName+"` isimli yazı kanalı yayın etkinlik yazı kanalları listesinde bulunamadı.")
			return
		}
		ephemeralRespond(s, i, "`"+channelName+"` isimli yazı kanalı yayın etkinlik yazı kanalları listesinden kaldırıldı.")

		// del-twitch announcement
	case "announcement":
		options = options[0].Options
		switch options[0].Name {
		// del-twitch announcement default-channel
		case "default-channel":
			ok, err := service.DeleteDiscordBotConfig(ctx, i.GuildID, "stream_anno_default_channel")
			if err != nil {
				log.Printf("Error while deleting Discord bot config: %v", err)
				ephemeralRespond(s, i, config.ErrorMessage+"#0001")
				return
			}

			if !ok {
				ephemeralRespond(s, i, config.ErrorMessage+"#0002")
				return
			}
			ephemeralRespond(s, i, "Varsayılan Twitch canlı yayın duyuru kanalı ayarı kaldırıldı.")

			// del-twitch announcement default-content
		case "default-content":
			_, err := service.DeleteDiscordBotConfig(ctx, i.GuildID, "stream_anno_default_content")
			if err != nil {
				log.Printf("Error while setting Discord bot config: %v", err)
				ephemeralRespond(s, i, config.ErrorMessage+"#0001")
				return
			}

			ephemeralRespond(s, i, "Yayın duyuru mesajı içeriği varsayılan olarak ayarlandı: `{stream.user}, {stream.category} yayınına başladı! {stream.url}`")

			// del-twitch announcement custom-content
		case "custom-content":
			options = options[0].Options
			twitchUsername := options[0].StringValue()
			twitchUsername = gosenchabot.ParseTwitchUsernameURLParam(twitchUsername)

			response0, uInfo := streamer.GetTwitchUserInfo(twitchUsername, c.twitchAccessToken)
			if response0 != "" {
				ephemeralRespond(s, i, response0)
				return
			}

			ok, err := service.UpdateTwitchStreamerAnnoContent(ctx, uInfo.ID, i.GuildID, nil)
			if err != nil {
				log.Printf("Error while deleting streamer announcement custom content: %v", err)
				ephemeralRespond(s, i, config.ErrorMessage+"del-twitch:custom-content#0001")
				return
			}

			if !ok {
				ephemeralRespond(s, i, config.ErrorMessage+"del-twitch:custom-content#0002")
				return
			}

			ephemeralRespond(s, i, twitchUsername+" kullanıcı adlı Twitch yayıncısına özgü yayın duyuru mesajı silindi.")
		}
	}
}

func DelTwitchCommandMetadata() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "del-twitch",
		Description: "Delete configuration setting.",
		DescriptionLocalizations: &map[discordgo.Locale]string{
			discordgo.Turkish: "Yapılandırma ayarlarını kaldır.",
		},
		DefaultMemberPermissions: &setdeletePermissions,
		Options: []*discordgo.ApplicationCommandOption{
			// del-twitch streamer
			{
				Name:        "streamer",
				Description: "Delete the stream from live stream announcements.",
				DescriptionLocalizations: map[discordgo.Locale]string{
					discordgo.Turkish: "Yayın duyuru mesajı atılan yayıncıyı sil.",
				},
				Type: discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "twitch-username-or-url",
						Description: "Twitch profile url or username",
						DescriptionLocalizations: map[discordgo.Locale]string{
							discordgo.Turkish: "Twitch kullanıcı profil linki veya kullanıcı adı",
						},
						Required: true,
					},
				},
			},
			// del-twitch announcement
			{
				Name:        "announcement",
				Description: "Annoucement group",
				Options: []*discordgo.ApplicationCommandOption{
					// del-twitch announcement default-channel
					{
						Name:        "default-channel",
						Description: "Delete the default channel configuration for live stream announcements.",
						DescriptionLocalizations: map[discordgo.Locale]string{
							discordgo.Turkish: "Yayın duyuru mesajlarının atılacağı varsayılan kanal ayarını kaldır.",
						},
						Type: discordgo.ApplicationCommandOptionSubCommand,
					},
					// del-twitch announcement default-content
					{
						Name:        "default-content",
						Description: "Delete the default announcement message content configuration.",
						DescriptionLocalizations: map[discordgo.Locale]string{
							discordgo.Turkish: "Varsayılan yayın duyuru mesajını sil.",
						},
						Type: discordgo.ApplicationCommandOptionSubCommand,
					},
					// del-twitch announcement custom-content
					{
						Name:        "custom-content",
						Description: "Delete the streamer specific custom live stream announcement message content.",
						DescriptionLocalizations: map[discordgo.Locale]string{
							discordgo.Turkish: "Yayıncıya özgü yayın duyuru mesajını sil.",
						},
						Type: discordgo.ApplicationCommandOptionSubCommand,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "twitch-username-or-url",
								Description: "Twitch profile url or username",
								DescriptionLocalizations: map[discordgo.Locale]string{
									discordgo.Turkish: "Twitch kullanıcı profil linki veya kullanıcı adı",
								},
								Required: true,
							},
						},
					},
				},
				Type: discordgo.ApplicationCommandOptionSubCommandGroup,
			},
			// del-twitch event-channel
			{
				Name:        "event-channel",
				Description: "Delete the live stream announcements channel setting to create Discord events for live streams.",
				DescriptionLocalizations: map[discordgo.Locale]string{
					discordgo.Turkish: "Canlı yayınların Discord etkinliklerini oluşturmak için canlı yayın duyuruları kanalını seç.",
				},
				Type: discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionChannel,
						Name:        "channel",
						Description: "The text channel where Twitch live stream announcements will be unfollowed",
						DescriptionLocalizations: map[discordgo.Locale]string{
							discordgo.Turkish: "Twitch yayın duyurularının takipten çıkarılacağı yazı kanalı",
						},
						ChannelTypes: []discordgo.ChannelType{
							discordgo.ChannelTypeGuildNews,
							discordgo.ChannelTypeGuildText,
						},
						Required: true,
					},
				},
			},
		},
	}
}
