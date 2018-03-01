# 유저생성
POST http://127.0.0.1:8080/register

## header
'Content-Type': 'application/json'

## body
```
{
    'username': username
}
```

## return
```
{
    'secret':시크릿
}
```

# 게임 만들기
POST http://127.0.0.1:8080/games/{game}
## header
'Content-Type': 'application/json'
## body
```
{
    'secret': 시크릿,
    'white': white유저아이디,
    'black': black유저아이디
}
```

# 게임 정보 가져오기
GET http://127.0.0.1:8080/games/{game}

## return
```
{
    "black":   흑유저아이디,
    "white":   백유저아이디,
    "board":   현재보드,
    "history": 히스토리,
    "initial": 처음보드,
    "list" 디버깅으로 만듬 삭제 예정
}
```

# 수 놓기
POST http://127.0.0.1:8080/games/{game}/actions

## header
'Content-Type': 'application/json'
## body
```
{
    'secret': 시크릿,
    'move': 수,
    'type': 'put'
}
```

# 웹소켓
http://127.0.0.1:8080/ws/games

## 로그인(전송)
```
{
    "type":"login",
    "secret":시크릿,
    "game": 게임
}
```

## 수가 놓아졋을시(수신)
```
{
    "type":"put",
    "move":수,
    "color":색깔
}
```