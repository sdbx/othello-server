# 

| 메소드 | 엔드포인트 | 설명 |
| --- | --- | --- |
| POST | /register | 클라이언트를 등록시킵니다 |
| GET | /rooms | 방 리스트를 구합니다 |
| POST | /rooms/{room} | 방을 팝니다 |
| DELETE | /rooms/{room} | 방을 없앱니다. |
| GET | /rooms/{room} | 특정 방의 정보를 가져옵니다 |
| POST | /connect | 특정 방에 들어갑니다 |
| DELETE | /connect | 방에서 나옵니다 |
| UPDATE | /rooms/{room}/{user} | 특정 룸의 특정 유저를 관리합니다 |

# 게임

## rest

| 메소드 | 엔드포인트 | 설명 |
| --- | --- | --- |
| POST | /move | 수를 둡니다 |
| DELETE | /move | 수를 무릅니다 |
| GET | /board | 현재 보드를 가져옵니다 |
| GET | /history | 현재 히스토리를 가져옵니다 |
| GET | /initial | 초기 보드를 가져옵니다 |
| POST | /surrender | 서랜을 칩니다 |

## websocket

모든 메세지는 TEXT포멧으로만 오고 모두 JSON형식을 따르며 메세지의 종류를 의미하는 type필드가 있습니다. 아래 서브헤더들의 제목은 type필드의 값 즉, 메세지의 종류를 의미합니다.

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


### login

서버에게 자신이 등록된 클라이언트임을 증명합니다. login이 성공적으로 이루어지지 않았을 경우 프로토콜 사용이 불가능합니다

```
{
  type:"login",
  secret:"시크릿"
}
```

성공응답:
```
{
  type:"success",
  from:"login"
  username:"사용자이름"
}
```

실패응답:
```
{
  type:"error",
  from:"login",
  msg:"에러메세지"
}
```

### join

게임에 접속합니다

```
{
  type:"join",
  id:"방id",
  password:"비밀번호" or x
}
```

성공응답:
```
{
  type:"success",
  form:"join",
  initial:초기게임보드,
  history:히스토리
}
```

### 수신

### turn

턴이 넘어갔음을 의미합니다

```
{
  type:"turn",
  now:"black" or "white",
  move:수
}
```

### end

게임이 끝났음을 의미합니다.

```
{
  type:"end",
  winner:"black" or "white"
  cause:"원인"
}
```

### tick

한 초가 지났음을 의미합니다.

```
{
  type:"tick",
  black_time:흑 남은 시간,
  white_time:백 남은 시간
}
```
