package main

var (
	teamMap = make(map[string]string);
	availableTeams = map[string]int {
		"red" : 1,
		"blue" : 2,
	}
)

func GetTeam(name string) string {
	team, ok := teamMap[name];
	if !ok {
		return "";
	}
	return team;
}

func JoinTeam(name string, team string) {
	if _,ok := availableTeams[team]; ok {
		teamMap[name] = team;
	}
}
