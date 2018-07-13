package main

import (
	"irc"
	"fmt"
)

func main() {
	nick := "twotwelvedegrees";
	pass := "oauth:xlqq3hr1t516rdnk3kmeo837fmw6dk";
	ircClient := irc.NewIRCClient("irc.chat.twitch.tv", 6667, nick, pass);
	ircClient.JOIN("twotwelvedegrees");
	//ircClient.PRIVMSG("hello there!");

	messagech := make(chan irc.TwitchMessage, 20);
	go ircClient.MainLoop(messagech);
	go ircClient.LoopInput();

	fmt.Println("Hey! I sent them all!");

	for m := range messagech {
		fmt.Println(m.Sender + ": " + m.Body);
		if command := NewCommand(m); command != nil {
			command.execute(ircClient);
		}
	}

}