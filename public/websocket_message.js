function connect_data(jsonData){
    team_color = jsonData["team_color"]
    if (team_color["0"] != -1) {
        document.getElementById('redTeamButton').classList.add('selected');
        document.getElementById('redTeamButton').innerText = '빨간 팀 선택됨';
    } 
    if (team_color["1"] != -1) {
        document.getElementById('blueTeamButton').classList.add('selected');
        document.getElementById('blueTeamButton').innerText = '파란 팀 선택됨';
    }
}

function teamremove(jsonData){
    document.getElementById('redTeamButton').classList.remove('selected');
    document.getElementById('redTeamButton').innerText = '빨간 팀';
    document.getElementById('blueTeamButton').classList.remove('selected');
    document.getElementById('blueTeamButton').innerText = '파란 팀';
}

function teamSelect(team_color){
    if (team_color["0"] != -1) {
        document.getElementById('redTeamButton').classList.add('selected');
        document.getElementById('redTeamButton').innerText = '빨간 팀 선택됨';
         
    } 
    if (team_color["1"] != -1) {
        document.getElementById('blueTeamButton').classList.add('selected');
        document.getElementById('blueTeamButton').innerText = '파란 팀 선택됨';
    }
}

function createsquare(jsonData){
    var gSquareValues = jsonData.game_square;

    for (let i = 0; i < 25; i++) { // 5x5 격자, 총 25개의 작은 정사각형
        const smallSquare = document.createElement('div');
        smallSquare.classList.add('small-square');
        smallSquare.classList.add(`square${i+1}`);
        // 작은 정사각형에 클릭 이벤트 추가

        function handleTouchOrMouse(event) {
            selectsquare(i + 1); // i값을 전달하여 해당 작은 정사각형의 번호를 사용
        }
        smallSquare.addEventListener('touchstart', handleTouchOrMouse);
        smallSquare.addEventListener('mousedown', handleTouchOrMouse);
        if (i % 2 == 0){
            smallSquare.classList.add("red");
        }else{
            smallSquare.classList.add("blue");
        }
        bigSquare.appendChild(smallSquare); // 큰 정사각형에 추가
    }
}

function gameSquare(jsonData){
    // console.log(jsonData.g_square);
    var gSquareValues = jsonData.game_square;

    bluescore = jsonData.score;

    var redscoreArea = document.getElementById('redscoreArea');
    redscoreArea.textContent = "레드팀 점수 : " + (25 - bluescore); 

    var bluescoreArea = document.getElementById('bluescoreArea');
    bluescoreArea.textContent = "블루팀 점수 : " + bluescore; 

    for (var i = 1; i <= 25; i++) {
        // var square = document.getElementsByClassName('square' + i);
        var square = document.querySelector('.square' + i);
        square.classList.remove('blue');
        square.classList.remove('red');

        // 해당 번호 (number)에 대한 g_square 값 가져오기
        var gSquareValue = gSquareValues[i];
        // g_square 값에 따라 색상 업데이트
        if (gSquareValue === 0) {
            square.classList.add('red');
        } else if (gSquareValue === 1) {
            square.classList.add('blue');
        }
    }
}


// 팀 선택 버튼을 눌렸을때
function selectTeam(team) {
    // MyClientID도 함께 보내고 싶으면 아래와 같이 처리할 수 있습니다
    if (team == "red"){
        myteam = 0
        document.getElementById('teamCursor').classList.remove('black');
        document.getElementById('teamCursor').classList.add('red');
    }else{
        myteam = 1
        document.getElementById('teamCursor').classList.remove('black');
        document.getElementById('teamCursor').classList.add('blue');
    }
    socket.send(JSON.stringify({ status : teamWhereStatus, team : team}));
}

// 사각형을 눌렸을때
function selectsquare(number){
    number = parseInt(number, 10);
    socket.send(JSON.stringify({ status : squareStatus, game_status: game_status, myteam : myteam, number: number }));
}