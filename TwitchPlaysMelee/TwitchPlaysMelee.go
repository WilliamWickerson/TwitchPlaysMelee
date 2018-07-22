package main

import (
	"irc"
	"fmt"
	"scripting"
	"controller"
	"controller/vJoy"
	"scripting/token"
	"scripting/AST"
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
	scanner := scripting.NewScanner("wavedash 3 left");
	parser := scripting.NewParser(scanner);
	gcSlice := make([]controller.GamecubeController, 60);
	for _,c := range parser.Parse() {
		if c == nil {
			fmt.Println("Fuck it's nil!");
			continue;
		}
		c.Execute(gcSlice);
		controller.PrintGC(gcSlice[0]);
		controller.PrintGC(gcSlice[3]);
	}
	for t := scanner.NextToken(); t.Type() != token.EOF; t = scanner.NextToken() {
		fmt.Printf("%s %d\n", t.Text(), t.Type());
	}
	vJoyC := vJoy.NewVJoyController();
	controller.SetvJoy(vJoyC, gcSlice[0]);
	fmt.Println(AST.GetMacro("wavedash", 2));
}