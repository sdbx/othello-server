# 개요

인코딩은 모두 utf-8입니다. 인공지능 봇을 만들고 싶으신 분은 쭉 내려가셔서 게임부분부터 읽으시면 됩니다.

# 유저 관련

## rest 

### 요약
| 메소드 | 엔드포인트 | 설명 |
| --- | --- | --- |
| POST | **/register/test** | 테스트 유저 등록입니다 |
| POST | /register/naver | 네이버 유저 등록입니다 |
| GET | /users/me | 클라이언트의 정보를 가져옵니다 |

### /register/test

#### 바디

```
{
  username: 유저이름
}
```

#### 응답

```
{
  secret: 유저시크릿
  username: 유저이름
} 
```

### /users/me

네이버 토큰 또는 유저 시크릿을 이용해 클라이언트의 정보를 가져옵니다. 둘 중 하나만 있어도 동작하며 네이버 토큰의 경우 클라이언트가 등록되지 않았을 경우 자동으로 등록후 정보를 반환해줍니다.

#### 헤더

X-User-Secret : 유저 시크릿

#### 응답

```
{
  secret: 유저시크릿
  username: 유저이름
} 
```

# 방

## rest

### 헤더

X-User-Secret : 유저 시크릿

### 요약

| 메소드 | 엔드포인트 | 설명 |
| --- | --- | --- |
| GET | /rooms | 방 리스트를 구합니다 |
| GET | /rooms/{room}?password=비밀번호 | 특정 방의 정보를 가져옵니다 |

### /rooms

#### 응답

```
{
  rooms:[
    {
      name: 방이름
      king: 방장이름
      password: true or false
    }
    ...
  ]
}
```

### /rooms/{room}?password=비밀번호

#### 응답

```
{
  name: 방이름
  king: 방장이름
  password: true or false
  status: "ingame" or "outgame"
  game: 게임id or x 
}
```

## 웹소켓

존재하지 않는 방에 login 할 시 방이 자동으로 파집니다. 방은 게임이 시작되었을 경우 모든 클라이언트들이 방 웹소켓과 접속이 끊어져도 방은 유지되나 준비중의 방의 경우는 방이 사라집니다. 그렇기에 방과 게임 페이지를 분리해도 게임이 끝났을 때 방 페이지로 돌아와 다시 방 웹소켓에 연결한다는 보장이 있으면 문제가 없습니다.

### 송신

### ping

핑!

```
{
  type:"ping"
}
```

응답:
```
{
  type:"pong"
}
```

웹소켓의 특성상 지속적으로 메세지를 보내지 않으면 타임아웃에 걸려버립니다. ping메세지를 일정 주기로 보내주십시오.

예:
```js
function keepAlive() {
    var timeout = 20000;
    if (webSocket.readyState == webSocket.OPEN) {
        webSocket.send('{type:"ping"}');
    }
    timerId = setTimeout(keepAlive, timeout);
}
keepAlive();
```


### **enter**

서버에게 자신이 등록된 클라이언트임을 증명합니다. login이 성공적으로 이루어지지 않았을 경우 프로토콜 사용이 불가능합니다

```
{
  type:"login"
  secret:"시크릿"
  room:"방"
}
```

실패응답:
```
{
  type:"error"
  from:"login"
  msg:"에러메세지"
}
```

### **actions**

만약 방장일 경우 어떤 행동을 취합니다.

강퇴(브로드캐스트)

```
{
  type:"actions"
  target:"유저이름"
  action:"kick"
}
``` 

방장넘기기(브로드캐스트)

```
{
  type:"actions"
  target:"유저이름"
  action:"king"
}
``` 

흑으로 만들기(브로드캐스트)

```
{
  type:"actions"
  to:"유저이름"
  action:"color"
  color:"black"
}
``` 

백으로 만들기(브로드캐스트)

```
{
  type:"actions"
  to:"유저이름"
  action:"color"
  color:"white"
}
``` 

게임 타입 변경

```
{
  type:"actions"
  action:"typechange"
  initial: 게임보드
  size: {
    w: 가로
    h: 세로
  }
}
```

게임 시작하기

```
{
  type:"actions"
  action:"gamestart"
}
```

### 수신

### **disconnect**

누군가의 접속이 끊어졌을 때 생깁니다. 만약 방장의 접속이 끊어진 경우 다음 방장의 아이디도 포함됩니다.

```
{
  type:"disconnect"
  username:유저이름
  next_king:유저이름 or x
}
```

### **connect**

누군가의 접속했을 때 생깁니다

```
{
  type:"connect"
  username:유저이름
}
```

### **gamestart**

게임이 시작되었음을 의미합니다.

```
{
  type:"gamestart"
  game:게임id
}
```


### **gameend**

게임이 끝났음을 의미합니다.

```
{
  type:"gameend"
}
```

# 게임

/ws/games

흑과 백이 모두 한 유저여도 상관이 없습니다. 이 경우 홀수번째 수놓기는 흑이고 짝수번째 수놓기는 백입니다. 다만 이경우는 수무르기가 동작하지 않습니다.

## 자료형

게임보드는 정수들로 이루어진 2차원 배열입니다. 각 숫자가 의미하는 것은 아래와 같습니다

| 숫자 | 의미 |
| --- | --- |
| 0 | 흑돌 |
| 1 | 백돌 |
| 2 | 공백 |


수는 기보형식으로 된 돌의 위치를 의미합니다. 수는 앞에 a-z까지의 글자와 뒤에 숫자들로 이루어져 있습니다. a는 0을 의미하며 b는 2를 의미하며 ... z는 25를 의미합니다. 1은 0을 의미하며 2는 1을 의미하며 .... 26은 25를 의미합니다. 게임보드[숫자][알파벳]으로 이 위치의 돌을 구할 수 있습니다. 


히스토리는 게임수들로 이루어진 1차원 배열입니다. 짝수번째의 인덱스의 값들은 흑의 수를 의미하며 홀수번째의 인덱스의 값들은 백의 수를 의미합니다. 히스토리 안에서는 수가 none인 경우도 있는데 이는 둘 수 있는 수가 없어서 넘겨진 것으로 아무 곳에도 두지 않았다는 것입니다. 히스토리 안에서 none수가 있다고 클라이언트가 none수를 둘 수 있는 것은 아닙니다. 만약 이 것이 가능해진다면 오델로 공식 룰에 어긋나게 되기 때문입니다. 그렇기 때문에 히스토리 안이나 웹소켓으로 통보받은 none수는 처리하되 none수를 보내지는 않도록 설계해야합니다.

## rest

### 헤더

X-User-Secret : 유저 시크릿

### 요약 

| 메소드 | 엔드포인트 | 설명 |
| --- | --- | --- |
| GET | **/games/{id}** | 현재 게임에 대한 정보를 가져옵니다 |
| POST | **/games/{id}/actions** | 게임에 뭔짓을 합니다 |
| POST | **/games/{id}** | 게임을 만듭니다(테스트) |

### /games/{id}

#### 응답

```
{
  room:방이름
  board:현재게임보드
  history:히스토리
  initial:초기게임보드
  usernames:{
    black:흑 유저이름
    white:백 유저이름
  }
  times:{
    black:흑 남은시간(초)
    white:백 남은시간(초)
  }
  total: {
    black: 흑돌수
    white: 백돌수
  }
}
```

### /games/{id}/actions

#### **수놓기**

```
{
  type:"put"
  move:수(기보형식)
}
```

#### 수무르기 신청

```
{
  type:"undo"
}
```

#### 수무르기 응답

```
{
  type:"undo_answer"
  answer:"yes" or "no"
}
```

## websocket

처음 login을 성공적으로 마쳤을 경우 웹소켓은 게임의 변화를 통보합니다.  

### 송신

### **ping**

핑!

```
{
  type:"ping"
}
```

응답:
```
{
  type:"pong"
}
```

웹소켓의 특성상 지속적으로 메세지를 보내지 않으면 타임아웃에 걸려버립니다. ping메세지를 일정 주기로 보내주십시오.

예:
```js
function keepAlive() {
    var timeout = 20000;
    if (webSocket.readyState == webSocket.OPEN) {
        webSocket.send('{type:"ping"}');
    }
    timerId = setTimeout(keepAlive, timeout);
}
keepAlive();
```


### **enter**

서버에게 자신이 등록된 클라이언트임을 증명합니다.

```
{
  type:"login"
  secret:"시크릿"
  game:"게임id"
}
```

실패응답:
```
{
  type:"error"
  from:"login"
  msg:"에러메세지"
}
```

### 수신

### **turn**

턴이 넘어갔음을 의미합니다

```
{
  type:"turn"
  color:"black" or "white"
  move:수
}
```

### **end**

게임이 끝났음을 의미합니다.

```
{
  type:"end"
  winner:"black" or "white"
  cause:"원인"
  total: {
    black: 흑돌수
    white: 백돌수
  }
}
```

### **tick**

10초에 한번씩 보내집니다.

```
{
  type:"tick"
  black:흑 남은 시간(초)
  white:백 남은 시간(초)
}
```

### undo

수무르기를 신청했음을 의미합니다. 만약 색이 상대편의 색이라면 무언가 해줘야합니다.

```
{
  type:"undo"
  color:"black" or "white"
}

```

### undo_answer

수무르기에 대한 응답을 의미합니다. 만약 answer가 yes라면 수 index로 돌아가야 합니다.

```
{
  type:"undo"
  answer:"yes" or "no"
  color:"black" or "white"
  move: 수index
}
```