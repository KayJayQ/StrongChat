import socket
import sys

def netcat(hostname, port, content = "GET 127.0.0.1:8080 HTTP/1.1"):
    port = int(port)
    content = content.encode("utf-8")
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect((hostname, port))
    s.sendall(content)
    s.shutdown(socket.SHUT_WR)
    while True:
        data = s.recv(1024)
        if data == "":
            break
        print("Received:", repr(data))
    print("Connection closed.")
    s.close()

if __name__ == "__main__":
    if len(sys.argv) > 3:
        netcat(sys.argv[1], sys.argv[2], sys.argv[3])
    else:
        netcat(sys.argv[1], sys.argv[2])
