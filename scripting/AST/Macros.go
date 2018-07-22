package AST

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"fmt"
)

type Macro struct {
	Name string `json:"name"`;
	Text string `json:"text"`;
}

var (
	macroMap = make(map[string]string);
)

func GetMacro(name string, inputs int) string {
	//Get the macro corresponding to the name and number of inputs
	mapKey := name + "[" + strconv.Itoa(inputs) + "]";
	text, ok := macroMap[mapKey];
	//If it doesn't exist then just return an empty string
	if !ok {
		return "";
	}
	return text;
}

func init() {
	//Open the file, if it can't be opened then exit
	loginFile, err := os.Open("../scripting/AST/macros.json");
	if err != nil {
		fmt.Println(err.Error());
		os.Exit(1);
	}
	//Make sure the file gets closed
	defer loginFile.Close();
	//Read in all of the bytes from the file
	bytes, _ := ioutil.ReadAll(loginFile);
	//Create a login object and have go parse the .json file
	macros := make([]Macro,0);
	if json.Unmarshal(bytes, &macros) != nil {
		fmt.Println(err.Error());
	}
	//Read through the macros and stick them in the map
	for _,macro := range macros {
		macroMap[macro.Name] = macro.Text;
	}
}