import asyncio
import websockets
import requests
import json

mode = input("make user? (y/n)")
if mode == "y":
    username = input("username ")
    r = requests.post('http://127.0.0.1:8080/register', data=json.dumps({
        'username': username
    }), headers={
        'Content-Type': 'application/json'
    })

    print('status:' + str(r.status_code))
    print('secret:' + r.json()['secret'])
    secret = r.json()['secret']
else:
    secret = input("secret ")
    
game = input("gamename ")
white = input("white ")
black = input("black ")
r = requests.post('http://127.0.0.1:8080/games/'+game, data=json.dumps({
    'secret': secret,
    'white': white,
    'black': black
}), headers={
    'Content-Type': 'application/json'
})

print(r.status_code)
print(r.text)
