const { SlashCommandBuilder } = require("@discordjs/builders");
const { dynamicVoice } = require("../utils/dynamicVoice");

const CHANNEL_NAME_OPTION = "channel-name";

module.exports = {
  data: new SlashCommandBuilder()
    .setName("voice")
    .setDescription("Create dynamic voice channels.")
    .addStringOption((option) =>
      option
        .setName(CHANNEL_NAME_OPTION)
        .setDescription("A proper channel name")
        .setRequired(true)
    ),
  execute(interaction) {
    const userId = interaction.user.id;

    if (dynamicVoice.userHasChannels(userId)) {
      interaction.reply({
        content: "You have already created a channel.",
        ephemeral: true,
      });
      return;
    }

    const channelName = interaction.options.getString(CHANNEL_NAME_OPTION);

    if (dynamicVoice.channelNameInUse(channelName)) {
      interaction.reply({
        content: "This channel name is aldready in use.",
        ephemeral: true,
      });
      return;
    }

    const guild = interaction.guild;

    guild.channels
      .create(channelName, {
        type: "GUILD_VOICE",
      })
      .then((channel) => {
        dynamicVoice.addChannel({
          channelId: channel.id,
          channelName: channelName,
          ownerId: userId,
        });

        interaction.reply({
          content: `Channel name: ${channelName}`,
          ephemeral: true,
        });
      })
      .catch((e) => {
        interaction.reply({
          content: `There was an error while creating voice channel "${channelName}". Error: ${e}`,
          ephemeral: true,
        });
      });
  },
};
