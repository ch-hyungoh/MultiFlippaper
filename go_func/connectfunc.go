package go_func

func Connectfunc(team_color map[int]int, jsonData map[string]interface{}) map[string]interface{} {

	jsonData["team_color"] = team_color

	return jsonData
}

func TeamWherefunc(clientID int, team_color map[int]int, team string) map[int]int {

	if team == "red" {
		if team_color[0] == -1 {
			team_color[0] = clientID
		}
	} else if team == "blue" {
		if team_color[1] == -1 {
			team_color[1] = clientID
		}
	}
	return team_color
}
