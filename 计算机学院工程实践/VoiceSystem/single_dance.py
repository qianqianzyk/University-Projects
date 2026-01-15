#-*- coding:UTF-8 -*-
import RPi.GPIO as GPIO
import time
import threading
import subprocess


#小车电机引脚定义
IN1 = 20
IN2 = 21
IN3 = 19
IN4 = 26
ENA = 16
ENB = 13

# 小车摄像头Y舵机引脚定义
CmrYPin = 9
CmrXPin = 11
# 小车摄像头舵机类
class Angle:
    def __init__(self, pin, agl = 0, lbound = -90, hbound = 90):
        self.pin = pin      #Pin
        self.agl = agl      #Angle
        GPIO.setup(pin, GPIO.OUT)
        self.pwm = GPIO.PWM(pin, 50)#PWM Signal
        self.pwm.start(0)
        self.low = lbound   #Min Bound
        self.high = hbound  #Max Bound
        self.pid = 0

    # def setpid(p, i, d):
    #     self.pid = PID.PositionalPID(p, i, d)

    def set(self, agl):
        if self.low <= agl <= self.high:
            self.pwm.ChangeDutyCycle(7.5 + agl/18)
            time.sleep(abs(agl - self.agl) * 0.005)
            self.pwm.ChangeDutyCycle(0)
            self.agl = agl

    def add(self, val):
        agl = self.agl + val
        if agl > self.high:
            agl = self.high
        elif agl < self.low:
            agl = self.low
        self.set(agl)

    def __del__(self):
        self.pwm.stop()

#设置GPIO口为BCM编码方式
GPIO.setmode(GPIO.BCM)

#忽略警告信息
GPIO.setwarnings(False)

#电机引脚初始化操作
def motor_init():
    global pwm_ENA
    global pwm_ENB
    global delaytime
    global pwm_servo
    GPIO.setup(ENA,GPIO.OUT,initial=GPIO.HIGH)
    GPIO.setup(IN1,GPIO.OUT,initial=GPIO.LOW)
    GPIO.setup(IN2,GPIO.OUT,initial=GPIO.LOW)
    GPIO.setup(ENB,GPIO.OUT,initial=GPIO.HIGH)
    GPIO.setup(IN3,GPIO.OUT,initial=GPIO.LOW)
    GPIO.setup(IN4,GPIO.OUT,initial=GPIO.LOW)
    #设置pwm引脚和频率为2000hz
    pwm_ENA = GPIO.PWM(ENA, 2000)
    pwm_ENB = GPIO.PWM(ENB, 2000)
    pwm_ENA.start(0)
    pwm_ENB.start(0)
    # 设置小车XY舵机
    GPIO.setup(CmrYPin, GPIO.OUT)
    GPIO.setup(CmrXPin, GPIO.OUT)
    #设置pwm引脚和频率为50hz
    pwm_servo = GPIO.PWM(CmrYPin, 50)
    pwm_servo.start(2.5)
    pwm_servox = GPIO.PWM(CmrXPin, 50)
    pwm_servox.start(2.5)

# Y舵机点头
def servo_control():
    for _ in range (4):
        for pos in range(91):
            pwm_servo.ChangeDutyCycle(2.5 + 10 * pos/180)
            
            time.sleep(0.009) 
        for pos in reversed(range(91)):
            pwm_servo.ChangeDutyCycle(2.5 + 10 * pos/180)
            
            time.sleep(0.009)

#小车前进	
def run(delaytime):
    GPIO.output(IN1, GPIO.HIGH)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.HIGH)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(80)
    pwm_ENB.ChangeDutyCycle(80)
    time.sleep(delaytime)

#小车后退
def back(delaytime):
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.HIGH)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.HIGH)
    pwm_ENA.ChangeDutyCycle(80)
    pwm_ENB.ChangeDutyCycle(80)
    time.sleep(delaytime)

#小车左转	
def left(delaytime):
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.HIGH)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(80)
    pwm_ENB.ChangeDutyCycle(80)
    time.sleep(delaytime)

#小车右转
def right(delaytime):
    GPIO.output(IN1, GPIO.HIGH)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(80)
    pwm_ENB.ChangeDutyCycle(80)
    time.sleep(delaytime)

#小车原地左转
def spin_left(delaytime):
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.HIGH)
    GPIO.output(IN3, GPIO.HIGH)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(80)
    pwm_ENB.ChangeDutyCycle(80)
    time.sleep(delaytime)

#小车原地右转
def spin_right(delaytime):
    GPIO.output(IN1, GPIO.HIGH)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.HIGH)
    pwm_ENA.ChangeDutyCycle(80)
    pwm_ENB.ChangeDutyCycle(80)
    time.sleep(delaytime)

#小车停止	
def brake(delaytime):
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(80)
    pwm_ENB.ChangeDutyCycle(80)
    time.sleep(delaytime)

def dancemove():
    time.sleep(1)
    # 原地左右转切换
    for _ in range(4):
        spin_left(0.2)
        time.sleep(0.1)
        spin_right(0.2)
        
    # 快速前后移动
    for _ in range(3):
        run(0.2)
        time.sleep(0.1)
        back(0.2)

    # 停顿
    brake(0.5)

    # 原地右转2圈
    spin_right(4)
    # 原地右转2圈
    spin_left(4)

    # 左右摆动
    for _ in range(3):
        left(0.2)
        time.sleep(0.1)
        right(0.2)

    # 快速旋转
    for _ in range(2):
        spin_left(0.2)
        time.sleep(0.1)
        spin_right(0.2)

    # 停顿结束舞蹈
    brake(1)
    
def play_music():
    subprocess.run(['mpg123', 'dancing.mp3'])

#延时2s
time.sleep(2)

#try/except语句用来检测try语句块中的错误，
#从而让except语句捕获异常信息并处理。
#小车循环前进1s，后退1s，左转2s，右转2s，原地左转3s
#原地右转3s，停止1s。

try:
    motor_init()
    while True:
        
        # 创建音乐线程
        music_thread = threading.Thread(target=play_music)
       

        # 启动音乐播放
        music_thread.start()
        
        # 启动舞蹈（与音乐同步）
        dancemove()
        servo_control()

        # 等待音乐播放完成
        music_thread.join()

        time.sleep(1)
except KeyboardInterrupt:
    pass
    
pwm_ENA.stop()
pwm_ENB.stop()
GPIO.cleanup()

