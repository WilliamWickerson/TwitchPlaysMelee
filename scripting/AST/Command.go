package AST

import (
	"scripting/token"
	"controller"
	"fmt"
)

type Command interface {
	Execute(gcArray []controller.GamecubeController);
}

type buttonCommand struct {
	command token.Type;
	button token.Type;
	duration Duration;
}

func NewButtonCommand(command token.Type, button token.Type, duration Duration) Command {
	return buttonCommand{command, button, duration};
}

func (bc buttonCommand) Execute(gcArray []controller.GamecubeController) {
	fmt.Println("Button Command!");
	for _,frame := range bc.duration.Frames() {
		setting := true;
		if (bc.command == token.KW_UNPRESS) {
			setting = false;
		}
		switch (bc.button) {
			case token.KW_A:
				gcArray[frame].A = setting;
				break;
			case token.KW_B:
				gcArray[frame].B = setting;
				break;
			case token.KW_X:
				gcArray[frame].X = setting;
				break;
			case token.KW_Y:
				gcArray[frame].Y = setting;
				break;
			case token.KW_Z:
				gcArray[frame].Z = setting;
				break;
			case token.KW_START:
				gcArray[frame].START = setting;
				break;
			case token.KW_DLEFT:
				gcArray[frame].DLEFT = setting;
				break;
			case token.KW_DRIGHT:
				gcArray[frame].DRIGHT = setting;
				break;
			case token.KW_DUP:
				gcArray[frame].DUP = setting;
				break;
			case token.KW_DDOWN:
				gcArray[frame].DDOWN = setting;
				break;
		}
	}
}

type sliderCommand struct {
	command token.Type;
	button token.Type;
	value float64;
	duration Duration;
}

func NewSliderCommand(command token.Type, button token.Type, value float64, duration Duration) Command {
	return sliderCommand{command, button, value, duration};
}

func (sc sliderCommand) Execute(gcArray []controller.GamecubeController) {
	fmt.Println("Slider Command!");
	/*
	switch (sc.button) {
		case token.KW_L:
			gcArray[frame].B = setting;
		case token.KW_R:
			gcArray[frame].B = setting;
	}
	*/
}

type stickCommand struct {
	stick token.Type;
	direction Direction;
	duration Duration;
}

func NewStickCommand(stick token.Type, direction Direction, duration Duration) Command {
	return stickCommand{stick, direction, duration};
}

func (sc stickCommand) Execute(gcArray []controller.GamecubeController) {
	fmt.Println("Stick Command!");
}

type macroCommand struct {
	Command;
}
