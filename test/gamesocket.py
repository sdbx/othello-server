import asyncio
import websockets

game = input("game ")
secret = input("secret ")

async def recv():
    async with websockets.connect('ws://localhost:8080/ws/games') as websocket:
        await websocket.send('{"type":"login","secret":"' + secret + '","game":"' + game+ '"}')
        while True:
            text = await websocket.recv()
            print(text)


asyncio.get_event_loop().run_until_complete(recv())