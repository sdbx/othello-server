import asyncio
import websockets
import signal
import sys

game = input("game ")
secret = input("secret ")

async def recv(uri):
    async with websockets.connect(uri) as websocket:
        
        await websocket.send('{"type":"login","secret":"' + secret + '","game":"' + game+ '"}')
        while True:
            text = await websocket.recv()
            print(text)

def signal_handler(signal, frame):
        print('You pressed Ctrl+C!')
        sys.exit(0)
signal.signal(signal.SIGINT, signal_handler)
print('Press Ctrl+C')
asyncio.get_event_loop().run_until_complete(
    recv('ws://localhost:8080/ws/games'))