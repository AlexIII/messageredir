package userstrings

const (
	BotMsgHelp             = "Hello! I am a message redirect bot.\nCommands:\n/start - Generate a new token\n/config - Set preferences (timezone, SIM names)\n/end - Leave the service\n\nThis is an open-source project: https://github.com/AlexIII/messageredir"
	MsgRedirFmt            = "%s\n%s\n%s\n\n%s"
	UserRemoved            = "You have been removed from the system. Goodbye!"
	UserAdded              = "You are all set!\nYour token: %s\n\nUse the following URL to submit messages."
	UserPreferencesUpdated = "Preferences updated!\n\nCurrent preferences: %+v"
	UserPreferencesHint    = "Example command to change preferences:\n/config UTC: +4, sim1: AT&T, sim2: Germany SIM card\n\nCurrent preferences: %+v"
)
