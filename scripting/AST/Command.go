package AST

import (
	"scripting/token"
	"controller"
)

type Command interface {
	Execute(gcArray []controller.GamecubeController);
	executeDelay(gcArray []controller.GamecubeController, delay int);
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
	bc.executeDelay(gcArray, 0);
}

func (bc buttonCommand) executeDelay(gcArray []controller.GamecubeController, delay int) {
	for _,i := range bc.duration.Frames() {
		frame := i + delay;
		if frame < 0 || frame >= MAXFRAMEWINDOW {
			continue;
		}
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
	sc.executeDelay(gcArray, 0);
}

func (sc sliderCommand) executeDelay(gcArray []controller.GamecubeController, delay int) {
	for _,i := range sc.duration.Frames() {
		frame := i + delay;
		if frame < 0 || frame >= MAXFRAMEWINDOW {
			continue;
		}
		//If the value to set is greater than or equal to 1, set digital press
		digPress := sc.value >= 1;
		switch (sc.button) {
			case token.KW_L:
				gcArray[frame].L = digPress;
				gcArray[frame].LANA = sc.value;
				break;
			case token.KW_R:
				gcArray[frame].R = digPress;
				gcArray[frame].RANA = sc.value;
				break;
		}
	}
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
	sc.executeDelay(gcArray, 0);
}

func (sc stickCommand) executeDelay(gcArray []controller.GamecubeController, delay int) {
	for _,i := range sc.duration.Frames() {
		frame := i + delay;
		if frame < 0 || frame >= MAXFRAMEWINDOW {
			continue;
		}
		//Get the components of the direction and set the correct stick
		x,y := sc.direction.Components();
		switch (sc.stick) {
			case token.KW_STICK:
				gcArray[frame].STICKX = x;
				gcArray[frame].STICKY = y;
				break;
			case token.KW_CSTICK:
				gcArray[frame].CX = x;
				gcArray[frame].CY = y;
				break;
		}
	}
}

type macroCommand struct {
	commands []Command;
	delay int;
}

func NewMacroCommand(commands []Command, delay int) Command {
	return macroCommand{commands, delay};
}

func (mc macroCommand) Execute(gcArray []controller.GamecubeController) {
	mc.executeDelay(gcArray, 0);
}

func (mc macroCommand) executeDelay(gcArray[]controller.GamecubeController, delay int) {
	for _,command := range mc.commands {
		command.executeDelay(gcArray, mc.delay + delay);
	}
}