package main

import (
	"irc"
	"strings"
)

type Command interface {
	execute(ircClient irc.IRCClient);
}

type helpCommand struct {
	Command;
}

func NewCommand(tm irc.TwitchMessage) Command {
	//Instantiate matching object or return nil if unknown command
	switch(strings.ToLower(tm.Body)) {
		case "!help":
			return helpCommand{};
		default:
			return nil;
	}
}

func (hc helpCommand) execute(ircClient irc.IRCClient) {
	ircClient.PRIVMSG("No help for you!");
}
