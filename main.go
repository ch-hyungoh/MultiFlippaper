package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"std/github.com/ch-hyungoh/MultiFlippaper/go_func"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]int)
var maxClientID = 0
var boardID = 0

var connectStatus = 0
var teamWhereStatus = 1
var finish_team = 2
var squareStatus = 3
var game_start = 4
var new_player = 5

var team_color = map[int]int{
	0: -1,
	1: -1,
}

var reset_team_color = map[int]int{
	0: -1,
	1: -1,
}

var game_board = make(([][]interface{}), 0)

var player_board = make(([][]interface{}), 0)

var game_timer = make(([][]interface{}), 0)

var game_square = map[int]int{
	1:  0,
	2:  1,
	3:  0,
	4:  1,
	5:  0,
	6:  1,
	7:  0,
	8:  1,
	9:  0,
	10: 1,
	11: 0,
	12: 1,
	13: 0,
	14: 1,
	15: 0,
	16: 1,
	17: 0,
	18: 1,
	19: 0,
	20: 1,
	21: 0,
	22: 1,
	23: 0,
	24: 1,
	25: 0,
}

// 클라이언트가 팀을 선택하면 추가 해주기
func addClient(client *websocket.Conn) int {

	// 새로운 클라이언트 번호 생성
	maxClientID = maxClientID + 1

	// 클라이언트와 클라이언트 번호 매핑
	clients[client] = maxClientID

	return clients[client]
}

// 팀원들이 모이면 새로운 게임판과 인원이 추가 시켜준다.
func addGame_Board(player1 int, player2 int) int {
	// 새 항목을 만들어 game_board에 추가
	boardID = boardID + 1

	newplayer := []interface{}{player1, player2, boardID}
	player_board = append(player_board, newplayer)

	// 게임 판 얕은 복사로 만들어 주기
	copiedMap := make(map[int]int)

	// 원본 map의 모든 키-값 쌍을 복사
	for key, value := range game_square {
		copiedMap[key] = value
	}

	game_board = append(game_board, []interface{}{copiedMap})

	return boardID
}

func removeClient(smallestClients []int, client *websocket.Conn) {
	for client := range clients {
		log.Println(client, smallestClients)
	}
	if _, exists := clients[client]; exists {
		// 클라이언트가 존재하는 경우에만 삭제
		delete(clients, client)
	}
}

// 1초 슬립해서 줄어드는 타이머
func countdownTimer() {
	for i := 30; i >= 0; i-- {
		fmt.Printf("%d초\n", i)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("타이머 종료")
}

// /////////////////////////////////////////////////////
// //////웹 서버 실행시켜 주는 부분///////////////////////
// /////////////////////////////////////////////////////
func startHTTPServer() {
	http.Handle("/", http.FileServer(http.Dir("public")))
}

func startWebSocketServer() {
	http.HandleFunc("/ws", handleConnections)
	log.Println("웹 서버가 열렸음")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("웹 소켓 서버 시작 실패 : %v", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 클라이언트와 연결된 후에 데이터를 전송합니다.
	var initialData = map[string]interface{}{
		"msg": "시작하자말자 갑니다",
	}
	if err := ws.WriteJSON(initialData); err != nil {
		log.Printf("Error writing initial data: %v", err)
		return
	}
	defer ws.Close()

	addClient(ws)

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("WebSocket 연결이 정상적으로 닫혔습니다.")
			} else {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		// msg를 문자열로 변환
		msgStr := string(msg)

		// JSON 파싱
		var jsonData map[string]interface{}
		err1 := json.Unmarshal([]byte(msgStr), &jsonData)
		if err1 != nil {
			if websocket.IsCloseError(err1, websocket.CloseNormalClosure) {
				log.Println("json오류 맛탱이감.")
			}
		}

		if jsonData["status"] == float64(connectStatus) {
			go_func.Connectfunc(team_color, jsonData)
			ws.WriteJSON(jsonData)
		} else if jsonData["status"] == float64(teamWhereStatus) {
			team := jsonData["team"].(string)

			teamselect := true
			for _, colorvalue := range team_color {
				if colorvalue == clients[ws] {
					teamselect = false
				}
			}

			if teamselect {

				team_color = go_func.TeamWherefunc(clients[ws], team_color, team)
				jsonData["team_color"] = team_color
				jsonData["team"] = team

				ws.WriteJSON(jsonData)

				for client := range clients {
					err := client.WriteJSON(jsonData)
					if err != nil {
						log.Printf("Error broadcasting message: %v", err)
					}
				}

				// 팀 선택 인원이 모두 정해지면 게임 시작되거나 다시 구하기
				if team_color[0] != -1 && team_color[1] != -1 {
					log.Println("///////////////////////게임이 시작되었습니다.")
					game_id := addGame_Board(team_color[0], team_color[1])

					count := 0

					for nowws, teamValue := range clients {
						for _, colorValue := range team_color {
							if teamValue == colorValue {
								jsonData["team"] = colorValue
								jsonData["status"] = game_start
								jsonData["game_id"] = game_id
								nowws.WriteJSON(jsonData)
								count += 1
							}
						}
						if count == 2 {
							break
						}
					}

					for key, value := range reset_team_color {
						team_color[key] = value

					}

					jsonData["status"] = new_player
					jsonData["game_square"] = game_square
					jsonData["now_game"] = game_id

					for client := range clients {
						err := client.WriteJSON(jsonData)
						if err != nil {
							log.Printf("Error broadcasting message: %v", err)
						}
					}
					log.Println(game_board)
				}
			}
		} else if jsonData["status"] == float64(squareStatus) {
			squareNumber := int(jsonData["number"].(float64))
			game_status := int(jsonData["game_status"].(float64))
			game_team := int(jsonData["myteam"].(float64))

			a := game_board[game_status-1][0]
			aMap := a.(map[int]int)

			// 눌렸을 경우 그 팀으로 변경해준다.
			aMap[squareNumber] = game_team

			// 스코어
			score := 0

			for _, value := range aMap {
				score += value
			}

			jsonData["status"] = float64(squareStatus)
			jsonData["game_status"] = game_status
			jsonData["game_square"] = aMap
			jsonData["score"] = score

			log.Println(clients)
			log.Println(clients[ws])
			log.Println(player_board[game_status-1])

			originalSlice := player_board[game_status-1]
			newSlice := originalSlice[:len(originalSlice)-1]

			// 현재 게임보드에 있는 플레이어에게만 게임판 정보를 준다.
			for client, value := range clients {
				for _, val := range newSlice {
					if value == val {
						err := client.WriteJSON(jsonData)
						if err != nil {
							log.Printf("Error broadcasting message: %v", err)
						}
					}
				}
			}
		}
	}

}

func main() {
	go startHTTPServer()
	startWebSocketServer()
}
