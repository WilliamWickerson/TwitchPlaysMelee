package main

import (
	"irc"
	"fmt"
	"time"
	"controller"
	"controller/vJoy"
)

func smain() {
	//Get the login from login.json
	login := GetLogin();
	//Create an IRC client and join desired chat
	ircClient := irc.NewIRCClient("irc.chat.twitch.tv", 6667, login.Nick, login.Pass);
	ircClient.JOIN("twotwelvedegrees");
	//Create a message channel for the IRC client to output messages
	messagech := make(chan irc.TwitchMessage, 20);
	//Start up main receive loop and console input loop
	go ircClient.MainLoop(messagech);
	go ircClient.LoopInput();
	//Loop infinitely over the messages printing them to the console and executing any commands
	for m := range messagech {
		fmt.Println(m.Sender + ": " + m.Body);
		if command := NewCommand(m); command != nil {
			command.execute(ircClient);
		}
	}

}

func main() {
	gc1 := controller.NewGamecubeController();
	gc2 := controller.NewGamecubeController();
	gc3 := controller.NewGamecubeController();
	gc4 := controller.NewGamecubeController();
	gc2.RANA = .33;
	gc3.RANA = .66;
	gc4.RANA = 1;
	vJoy := vJoy.NewVJoyController();
	for {
		time.Sleep(500*time.Millisecond);
		controller.SetvJoy(vJoy, gc1);
		time.Sleep(500*time.Millisecond);
		controller.SetvJoy(vJoy, gc2);
		time.Sleep(500*time.Millisecond);
		controller.SetvJoy(vJoy, gc3);
		time.Sleep(500*time.Millisecond);
		controller.SetvJoy(vJoy, gc4);
	}
}