package main

import (
	"encoding/json"
	"log"
	"net/http"
	"multiflippaper/go_func"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]int)
var maxClientID = 0

var connectStatus = 0
var playStatus = 1
var teamWhereStatus = 2
var SquarStatus = 3

var team_color = map[int]int{
	0: -1,
	1: -1,
}

var team_square = map[int]int{
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

func removeClient(smallestClients []int, client *websocket.Conn) {
	for client := range clients {
		log.Println(client, smallestClients)
	}
	if _, exists := clients[client]; exists {
		// 클라이언트가 존재하는 경우에만 삭제
		delete(clients, client)
	}
}

/////////////////////////////////////////////////////////////
// 방에 2명이 들어오면 2명을 묶어서 게임 시작하게 해주기 수정해야함
/////////////////////////////////////////////////////////////

func getClientIDs(clients map[*websocket.Conn]int) []int {
	// 클라이언트 번호만 모아둘 슬라이스 생성
	clientIDs := make([]int, 0)

	// 모든 클라이언트 번호를 슬라이스에 추가
	for _, clientID := range clients {
		clientIDs = append(clientIDs, clientID)
	}

	// 슬라이스에서 가장 작은 2개의 클라이언트 번호 선택
	var smallest1, smallest2 int
	if len(clientIDs) > 0 {
		smallest1 = clientIDs[0]
		if len(clientIDs) > 1 {
			smallest2 = clientIDs[1]
		}
	}

	return []int{smallest1, smallest2}
}

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

	defer ws.Close()

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
			connectfunc.Connectfunc(jsonData)
		}
	}

}
