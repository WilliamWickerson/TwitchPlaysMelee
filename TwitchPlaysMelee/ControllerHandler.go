package main

import (
	"irc"
	"controller"
	"scripting"
	"scripting/AST"
	"controller/vJoy"
	"time"
)

var (
	controllerSemaphore = make(chan bool, 1);
	gcArrays = make(map[string][]controller.GamecubeController);
	vJoys = make(map[string]vJoy.VJoyController);
)

func init() {
	for team,_ := range availableTeams {
		gcArrays[team] = make([]controller.GamecubeController, AST.MAXFRAMEWINDOW);
		vJoys[team] = vJoy.NewVJoyController(availableTeams[team]);
	}
}

func ControllerLoop(messagech chan irc.TwitchMessage) {
	go NextWindowLoop();
	for m := range messagech {
		team := GetTeam(m.Sender);
		if team == "" {
			continue;
		}
		commands := scripting.NewParser(scripting.NewScanner(m.Body)).Parse();
		controllerSemaphore <- true;
		for _,command := range commands {
			command.Execute(gcArrays[team]);
		}
		go controller.SetvJoy(vJoys[team], gcArrays[team][0]);
		<-controllerSemaphore;
	}
}

func NextWindowLoop() {
	for ; true; time.Sleep(17 * time.Millisecond) {
		controllerSemaphore <- true;
		for team,_ := range availableTeams {
			go controller.SetvJoy(vJoys[team], gcArrays[team][0]);
			gcArrays[team] = append(gcArrays[team], controller.GamecubeController{})[1:AST.MAXFRAMEWINDOW+1];
		}
		<-controllerSemaphore;
	}
}
