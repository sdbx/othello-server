import websocket
try:
    import thread
except ImportError:
    import _thread as thread
import time

game = input("game ")
secret = input("secret ")

def on_message(ws, message):
    print(message)

def on_error(ws, error):
    print(error)

def on_close(ws):
    print("### closed ###")

def run(*args):
    time.sleep(1)
    ws.send('{"type":"enter","secret":"' + secret + '","room":"' + game+ '"}')
    print("thread terminating...")

thread.start_new_thread(run, ())

if __name__ == "__main__":
    websocket.enableTrace(True)
    ws = websocket.WebSocketApp("ws://127.0.0.1:8080/ws/rooms",
    on_message = on_message,
    on_error = on_error,
    on_close = on_close)
    ws.run_forever()