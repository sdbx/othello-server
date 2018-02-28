# 유저 관련

## rest 

| 메소드 | 엔드포인트 | 설명 |
| --- | --- | --- |
| POST | /register | 클라이언트를 등록시킵니다 |

# 방

## rest

| 메소드 | 엔드포인트 | 설명 |
| --- | --- | --- |
| GET | /rooms | 방 리스트를 구합니다 |
| POST | /rooms/{room} | 방을 팝니다 |
| GET | /rooms/{room} | 특정 방의 정보를 가져옵니다 |

## websocket

모든 메세지는 TEXT포멧으로만 오고 모두 JSON형식을 따르며 메세지의 종류를 의미하는 type필드가 있습니다. 아래 서브헤더들의 제목은 type필드의 값 즉, 메세지의 종류를 의미합니다.

/ws/rooms/{room}로 접속합니다

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

### action

만약 방장일 경우 어떤 행동을 취합니다.

강퇴(브로드캐스트)

```
{
  type:"action",
  target:"유저이름",
  action:"kick"
}
``` 

방장넘기기(브로드캐스트)

```
{
  type:"action",
  target:"유저이름",
  action:"king"
}
``` 

흑으로 만들기(브로드캐스트)

```
{
  type:"action",
  target:"유저이름",
  action:"black"
}
``` 

백으로 만들기(브로드캐스트)

```
{
  type:"action",
  target:"유저이름",
  action:"white"
}
``` 

게임 시작하기

```
{
  type:"action",
  action:"gamestart"
}
```

### 수신

### info

이 방에 대한 정보를 알려줍니다

```
{
  type:"info",
  participants:참가자들 유저아이디,
  king:방장 유저아이디,
  type:게임 타입,
  status:"ingame" or "preparing"
}
```

### disconnect

누군가의 접속이 끊어졌을 때 생깁니다. 만약 방장의 접속이 끊어진 경우 다음 방장의 아이디도 포함됩니다.

```
{
  type:"disconnect",
  who:유저아이디,
  next_king:유저아이디 or x
}
```

### connect

누군가의 접속했을 때 생깁니다

```
{
  type:"disconnect",
  who:유저아이디
}
```

### gamestart

게임이 시작되었음을 의미합니다.

```
{
  type:"gamestart"
}
```


### gameend

게임이 끝났음을 의미합니다.

```
{
  type:"gameend"
}
```

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

### 송신

### ping

위와 동일

### login

위와 동일

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
