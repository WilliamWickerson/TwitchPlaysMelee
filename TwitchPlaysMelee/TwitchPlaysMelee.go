package main

import (
	"irc"
	"fmt"
	"scripting"
	"controller"
	"scripting/token"
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
	scanner := scripting.NewScanner("press x 1; unpress y 3-12,15,17-19,42-90; stick left 3-10");
	parser := scripting.NewParser(scanner);
	for _,c := range parser.Parse().Commands {
		if c == nil {
			fmt.Println("Fuck it's nil!");
			continue;
		}
		c.Execute(make([]controller.GamecubeController,60));
	}
	for t := scanner.NextToken(); t.Type() != token.EOF; t = scanner.NextToken() {
		fmt.Printf("%s %d\n", t.Identifier(), t.Type());
	}
	t := token.New(token.FLOATLITERAL, "1");
	val,_ := t.Float();
	fmt.Printf("%f", val);
}