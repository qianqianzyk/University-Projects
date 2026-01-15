import socket
import RPi.GPIO as GPIO
import time
import struct
from mail_send import send_alert_email

class CmrSv:
    def __init__(self, xpin, ypin,
                 x_min_bound = -90, x_max_bound = 90,
                 y_min_bound = -90, y_max_bound = 90):
        self.xpin = xpin
        self.ypin = ypin
        self.xagl = 0
        self.yagl = 0
        self.xmin = x_min_bound
        self.xmax = x_max_bound
        self.ymin = y_min_bound
        self.ymax = y_max_bound
        GPIO.setup(xpin, GPIO.OUT)
        GPIO.setup(ypin, GPIO.OUT)
        
    def aglset(self, x_agl, y_agl):
        if x_agl < self.xmin:
            x_agl = self.xmin
        elif x_agl > self.xmax:
            x_agl = self.xmax
        if y_agl < self.ymin:
            y_agl = self.ymin
        elif y_agl > self.ymax:
            y_agl = self.ymax
        xpwm = GPIO.PWM(self.xpin, 50)
        ypwm = GPIO.PWM(self.ypin, 50)
        xpwm.start(7.5 + x_agl/18)
        ypwm.start(7.5 + y_agl/18)
        time.sleep(max(abs(x_agl-self.xagl), abs(y_agl - self.yagl)) * 0.008)
        xpwm.stop()
        ypwm.stop()
        self.xagl = x_agl
        self.yagl = y_agl

    def increase(self, x_val, y_val):
        self.aglset(self.xagl + x_val, self.yagl + y_val)

def receive_img(skt, img_addr):
    file_size_data = skt.recv(1024)
    file_size = struct.unpack('!I', file_size_data)[0]
    img = b''
    remain = file_size
    while remain > 0:
        chunk = skt.recv(min(4096, remain))
        if not chunk:
            print('Connect Interrupted')
            break
        img += chunk
        remain -= len(chunk)
    with open(img_addr, 'wb') as f:
        f.write(img)

print('Initializing Servo')
GPIO.setmode(GPIO.BCM)#Set GPIO Coding: BCM
GPIO.setwarnings(False)#Ignore Warning
cmrSv = CmrSv(11, 9, y_min_bound = -60, y_max_bound = 28)

print("Setting Socket")
cmrSocket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
cmrSocket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
host = "192.168.50.1"
port = 1629
cmrSocket.bind((host, port))
cmrSocket.listen(4)

frame_x = 160
frame_y = 120
cmrSv.aglset(0, -50)

print("Ready for Connecting")
client, addr = cmrSocket.accept()
print("New Connection")
print("IP is %s" % addr[0])
print("port is %d\n" % addr[1])
try:
    while True:
        msg = client.recv(1024)
        msg = msg.decode("utf-8")
        if msg != "":
            if msg == 'Stop':
                break
            mess = msg.split(' ')
            err_x = int(mess[0])
            err_y = int(mess[1])
            err_x = frame_x - err_x
            err_y = frame_y - err_y
            if abs(err_x) > 25 or abs(err_y) > 20:
                cmrSv.increase(err_x*0.04, err_y*0.04)
    receive_img(client, 'enermy_img.jpg')
    send_alert_email('pisy-2022@foxmail.com', 'enermy_img.jpg')
except KeyboardInterrupt:
    print('KeyboardInterrupt')
finally:
    print('Program Exit')
