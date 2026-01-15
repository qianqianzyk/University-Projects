import socket
import RPi.GPIO as GPIO
import time

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
        time.sleep(max(abs(x_agl-self.xagl), abs(y_agl - self.yagl)) * 0.01)
        xpwm.stop()
        ypwm.stop()
        self.xagl = x_agl
        self.yagl = y_agl

    def increase(self, x_val, y_val):
        self.aglset(self.xagl + x_val, self.yagl + y_val)

#小车电机引脚定义
IN1 = 20
IN2 = 21
IN3 = 19
IN4 = 26
ENA = 16
ENB = 13

#小车按键定义
key = 8

def init():
    global pwm_ENA
    global pwm_ENB
    GPIO.setup(ENA,GPIO.OUT,initial=GPIO.HIGH)
    GPIO.setup(IN1,GPIO.OUT,initial=GPIO.LOW)
    GPIO.setup(IN2,GPIO.OUT,initial=GPIO.LOW)
    GPIO.setup(ENB,GPIO.OUT,initial=GPIO.HIGH)
    GPIO.setup(IN3,GPIO.OUT,initial=GPIO.LOW)
    GPIO.setup(IN4,GPIO.OUT,initial=GPIO.LOW)
    GPIO.setup(key,GPIO.IN)
    #设置pwm引脚和频率为2000hz
    pwm_ENA = GPIO.PWM(ENA, 2000)
    pwm_ENB = GPIO.PWM(ENB, 2000)
    pwm_ENA.start(0)
    pwm_ENB.start(0)
	
#小车前进	
def run(leftspeed, rightspeed):
    GPIO.output(IN1, GPIO.HIGH)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.HIGH)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(leftspeed)
    pwm_ENB.ChangeDutyCycle(rightspeed)
    time.sleep(0.1)

#小车后退
def back(leftspeed, rightspeed):
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.HIGH)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.HIGH)
    pwm_ENA.ChangeDutyCycle(leftspeed)
    pwm_ENB.ChangeDutyCycle(rightspeed)
    time.sleep(0.1)

	
#小车左转	
def left(leftspeed, rightspeed):
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.HIGH)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(leftspeed)
    pwm_ENB.ChangeDutyCycle(rightspeed)
    time.sleep(0.1)


#小车右转
def right(leftspeed, rightspeed):
    GPIO.output(IN1, GPIO.HIGH)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(leftspeed)
    pwm_ENB.ChangeDutyCycle(rightspeed)
    time.sleep(0.1)

	
#小车原地左转
def spin_left(leftspeed, rightspeed):
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.HIGH)
    GPIO.output(IN3, GPIO.HIGH)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(leftspeed)
    pwm_ENB.ChangeDutyCycle(rightspeed)
    time.sleep(0.1)


#小车原地右转
def spin_right(leftspeed, rightspeed):
    GPIO.output(IN1, GPIO.HIGH)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.HIGH)
    pwm_ENA.ChangeDutyCycle(leftspeed)
    pwm_ENB.ChangeDutyCycle(rightspeed)
    time.sleep(0.1)


#小车停止	
def brake():
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.LOW)
    time.sleep(0.1)


turn_speed = 10
straight_speed = 10

print('Initializing Servo')
GPIO.setmode(GPIO.BCM)#Set GPIO Coding: BCM
GPIO.setwarnings(False)#Ignore Warning
cmrSv = CmrSv(11, 9, y_max_bound = 28)

#Initialize motor
print("Initialize motor")
init()

print("Setting Socket")
cmrSocket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
cmrSocket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
host = "192.168.50.1"
port = 1628
cmrSocket.bind((host, port))
cmrSocket.listen(4)

frame_x = 160
frame_y = 120
cmrSv.aglset(0, -60)

print("Ready for Connecting")
client, addr = cmrSocket.accept()
print("New Connection")
print("IP is %s" % addr[0])
print("port is %d\n" % addr[1])

last_face_time = time.time()
face_timeout = 2  # 2秒未检测到人脸则回中

try:
    while True:
        msg = client.recv(1024)
        msg = msg.decode("utf-8")
        if msg != "":
            if msg == "none":
                brake()
                print("未检测到人脸，停车")
                # 超时自动回中
                if time.time() - last_face_time > face_timeout:
                    cmrSv.aglset(0, -45)  # 回到初始角度
                    print("云台自动回中")
                    # 小车原地转两圈（不阻塞式）
                    spin_time = 2  # 总共转2圈，先假设1圈1秒明天调参测
                    step = 0.1     # 每次转0.1秒
                    steps = int(spin_time / step * 1)  # 2圈/0.1s
                    for _ in range(steps):
                        spin_right(10, 10)
                        # 检查是否有新消息
                        client.settimeout(0.01)
                        try:
                            new_msg = client.recv(1024)
                            if new_msg:
                                new_msg = new_msg.decode("utf-8")
                                if new_msg != "none":
                                    print("发现人脸，立即响应新指令")
                                    msg = new_msg
                                    break
                        except:
                            pass
                        time.sleep(step)
                    brake()
                continue
            # if msg == "none":
            #     brake()
            #     print("未检测到人脸，停车")
            #     # 超时自动回中
            #     if time.time() - last_face_time > face_timeout:
            #         cmrSv.aglset(0, -60)  # 回到初始角度
            #         print("云台自动回中")
            #         # 小车原地转两圈
            #         for _ in range(2):
            #             spin_right(10, 10)  # 速度可调整
            #             time.sleep(10)       # 1秒为一圈，视实际情况调整
            #         brake()
            #     continue
            last_face_time = time.time()
            mess = msg.split(' ')
            try:
              err_x = int(mess[0])
              err_y = int(mess[1])
              motocontrol_desicion = mess[2]
              err_x = frame_x - err_x
              err_y = frame_y - err_y
            except:
              err_x = 0
              err_y = 0
              motocontrol_desicion = ''
            if abs(err_x) > 25 or abs(err_y) > 18:
                cmrSv.increase(err_x*0.05, err_y*0.03)
            # 电机控制决策
            if abs(err_x) <= 25:
                if motocontrol_desicion == "00":
                    run(straight_speed, straight_speed)
                    print("前进")
                elif motocontrol_desicion == "01":
                    back(straight_speed, straight_speed)
                    print("后退")
            else:
                if motocontrol_desicion == "00" and err_x > 0:
                    spin_right(turn_speed,turn_speed)
                    run(straight_speed, straight_speed)
                    print("右转后直行前进")
                elif motocontrol_desicion =="00" and err_x < 0:
                    spin_left(turn_speed,turn_speed)
                    run(straight_speed, straight_speed)
                    print("左转后直行前进")
                elif motocontrol_desicion == "01" and err_x > 0:
                    spin_right(turn_speed,turn_speed)
                    back(straight_speed, straight_speed)
                    print("右转后，后退")
                elif motocontrol_desicion == "01" and err_x < 0:
                    spin_left(turn_speed,turn_speed)
                    back(straight_speed, straight_speed)
                    print("左转后，后退")
                elif motocontrol_desicion == "10":
                    spin_left(turn_speed, turn_speed)
                    print("左转")
                elif motocontrol_desicion == "11":
                    spin_right(turn_speed, turn_speed)
                    print("右转")
            cmrSv.aglset(0, -45)  # 回到初始角度
            

except KeyboardInterrupt:
    pass
finally:
    print('Program Exit')
