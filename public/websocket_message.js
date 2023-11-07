function connect_data(jsonData){
    console.log("연결됨");
    myClientId = jsonData.clientID;
    myClients = jsonData.clients;
    teamSelect(jsonData);
    var dataTextArea = document.getElementById('dataTextArea');
    var idtext = JSON.stringify(jsonData.clientID, null, 2);
    dataTextArea.textContent = idtext + "님 안녕하세요"; // JSON 데이터를 예쁘게 포맷팅하여 출력

    // 게임에 입장했을때 현재 진행중인 게임 데이터 보여주기 위해서
    var gSquareValues = jsonData.g_square;

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

function square(jsonData){
    var gSquareValues = jsonData.g_square;

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