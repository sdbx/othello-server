import asyncio
import websockets
import requests
import json

secret = input("secret ")
r = requests.post('http://127.0.0.1:8080/abc', data=json.dumps({
    'secret': secret
}), headers={
    'Content-Type': 'application/json'
})

print(r.status_code)
print(r.text)

async def test():
    async with websockets.connect('ws://127.0.0.1:8080/ws/games') as websocket:
        while True:
            content = input("content ")
            if content == "exit":
                break
            await websocket.send(content)
            print("> {}".format(content))

            back = await websocket.recv()
            print("< {}".format(back))

asyncio.get_event_loop().run_until_complete(test())