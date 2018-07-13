package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type Login struct {
	Nick string `json:"nick"`;
	Pass string `json:"pass"`;
}

func GetLogin() Login {
	//Open the login file, if it can't be opened then exit
	loginFile, err := os.Open("login.json");
	if err != nil {
		fmt.Println(err);
		os.Exit(1);
	}
	//Make sure the file gets closed
	defer loginFile.Close();
	//Read in all of the bytes from the file
	bytes, _ := ioutil.ReadAll(loginFile);
	//Create a login object and have go parse the .json file
	var login Login;
	json.Unmarshal(bytes, &login);
	return login;
}
