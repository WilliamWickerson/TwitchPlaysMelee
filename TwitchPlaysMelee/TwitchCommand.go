package main

import (
	"irc"
	"strings"
)

type Command interface {
	execute(ircClient irc.IRCClient);
}

type helpCommand struct {
}

type joinCommand struct {
	message irc.TwitchMessage;
}

type teamCommand struct {
	message irc.TwitchMessage;
}

func NewCommand(tm irc.TwitchMessage) Command {
	//Instantiate matching object or return nil if unknown command
	if len(tm.Body) < 5 {
		return nil;
	}
	switch(strings.ToLower(tm.Body)[0:5]) {
		case "!help":
			return helpCommand{};
		case "!join":
			return joinCommand{tm};
		case "!team":
			return teamCommand{tm};
		default:
			return nil;
	}
}

func (hc helpCommand) execute(ircClient irc.IRCClient) {
	go ircClient.PRIVMSG("Use !join to join a team, the teams are: red, blue");
}

func (jc joinCommand) execute(ircClient irc.IRCClient) {
	//Join the team, if the user finds himself on the team print a joined message
	if len(jc.message.Body) >= 7 {
		newTeam := strings.ToLower(jc.message.Body[6:len(jc.message.Body)]);
		JoinTeam(jc.message.Sender, newTeam);
		if team := GetTeam(jc.message.Sender); team == newTeam {
			go ircClient.PRIVMSG(jc.message.Sender + " joined team " + team);
		}
	}
}

func (tc teamCommand) execute(ircClient irc.IRCClient) {
	team := GetTeam(tc.message.Sender);
	if team == "" {
		go ircClient.PRIVMSG(tc.message.Sender + " does not have a team");
	} else {
		go ircClient.PRIVMSG(tc.message.Sender + " is on team " + team);
	}
}
