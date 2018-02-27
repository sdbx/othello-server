# rest

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

# 프로토콜

모든 메세지는 TEXT포멧으로만 오고 모두 JSON형식을 따르며 메세지의 종류를 의미하는 type필드가 있습니다. 아래 서브헤더들의 제목은 type필드의 값 즉, 메세지의 종류를 의미합니다.

## ping

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

## 게임

인간 게임 플레이의 경우 login을 하고 join을 하면 됩니다. 만약 인공지능을 데리고 올 경우, 인공지능은 login후 join을 통해 플레이어로써 접속을 하고 클라이언트의 경우 그냥 join을 함으로써 관전자로 접속하는 식으로 하여야합니다.

### login

서버에게 자신이 등록된 클라이언트임을 증명합니다. login이 성공적으로 이루어지지 않았을 경우 게임플레이가 불가능합니다.

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

### set

오델로 말을 둡니다

```
{
  type:"set",
  move:수
}
```

성공응답:
```
{
  type:"success",
  from:"set"
}
```

실패응답:
```
{
  type:"error",
  from:"set",
  error:"에러메세지"
}
```

브로드캐스트:
```
{
  type:"set",
  move:수,
  player:"black" or "white"
}
```

### get

현재 보드판을 반환합니다

```
{
  type:"get"
}
```

성공응답:
```
{
  type:"success",
  from:"get",
  board:현재게임보드
}
```

실패응답:
```
{
  type:"error",
  from:"get",
  error:"에러메세지"
}
```

### surrender

서랜을 칩니다.

```
{
  type:"surrender"
}
```

성공응답:
```
{
  type:"success",
  from:"surrender"
}
```

실패응답:
```
{
  type:"error",
  from:"surrender",
  error:"에러메세지"
}
```

브로드캐스트:
```
{
  type:"surrender",
  by:"black" or "white"
}
```

### possible

가능한 수들을 반환합니다.

```
{
  type:"possible"
}
```

성공응답:
```
{
  type:"success",
  from:"possible",
  moves:수들
}
```

실패응답:
```
{
  type:"error",
  from:"possible",
  error:"에러메세지"
}
```

### undo

수를 무릅니다. 상대방이 수를 두기전에 빠릴 하셔야 합니다.

```
{
  type:"undo"
}
```

성공응답:
```
{
  type:"success",
  from:"undo"
}
```

실패응답:
```
{
  type:"error",
  from:"undo",
  error:"에러메세지"
}
```

브로드캐스트:
```
{
  type:"undo",
  by:"black" or "white"
}
```

## 수신

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
