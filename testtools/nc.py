import socket
import sys

def netcat(hostname = "127.0.0.1", port = 8080, content = "GET 127.0.0.1:8080 HTTP/1.1\r\n"):
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
        print("Received:", str(data, encoding="utf-8"))
    print("Connection closed.")
    s.close()

if __name__ == "__main__":

    if len(sys.argv) == 4:
        netcat(sys.argv[1], sys.argv[2], sys.argv[3])
    if len(sys.argv) == 3:
        netcat(sys.argv[1], sys.argv[2])
    if len(sys.argv) == 1:
        netcat()
