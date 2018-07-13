package main

import (
	"irc"
	"fmt"
)

func main() {
	login := GetLogin();
	ircClient := irc.NewIRCClient("irc.chat.twitch.tv", 6667, login.Nick, login.Pass);
	ircClient.JOIN("twotwelvedegrees");
	//ircClient.PRIVMSG("hello there!");

	messagech := make(chan irc.TwitchMessage, 20);
	go ircClient.MainLoop(messagech);
	go ircClient.LoopInput();

	for m := range messagech {
		fmt.Println(m.Sender + ": " + m.Body);
		if command := NewCommand(m); command != nil {
			command.execute(ircClient);
		}
	}

}