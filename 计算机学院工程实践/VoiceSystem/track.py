# -*- coding: UTF-8 -*-
# 循迹控制模块
import RPi.GPIO as GPIO
import threading
import time
from voice_system import listen_and_recognize

#小车电机引脚定义
IN1 = 20
IN2 = 21
IN3 = 19
IN4 = 26
ENA = 16
ENB = 13

#小车按键定义
key = 8

#循迹红外引脚定义
TrackSensorLeftPin1  =  3
TrackSensorLeftPin2  =  5
TrackSensorRightPin1 =  4
TrackSensorRightPin2 =  18

GPIO.setmode(GPIO.BCM)
GPIO.setwarnings(False)

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
    GPIO.setup(TrackSensorLeftPin1,GPIO.IN)
    GPIO.setup(TrackSensorLeftPin2,GPIO.IN)
    GPIO.setup(TrackSensorRightPin1,GPIO.IN)
    GPIO.setup(TrackSensorRightPin2,GPIO.IN)

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

#小车后退
def back(leftspeed, rightspeed):
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.HIGH)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.HIGH)
    pwm_ENA.ChangeDutyCycle(leftspeed)
    pwm_ENB.ChangeDutyCycle(rightspeed)
	
#小车左转	
def left(leftspeed, rightspeed):
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.HIGH)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(leftspeed)
    pwm_ENB.ChangeDutyCycle(rightspeed)

#小车右转
def right(leftspeed, rightspeed):
    GPIO.output(IN1, GPIO.HIGH)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(leftspeed)
    pwm_ENB.ChangeDutyCycle(rightspeed)
	
#小车原地左转
def spin_left(leftspeed, rightspeed):
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.HIGH)
    GPIO.output(IN3, GPIO.HIGH)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(leftspeed)
    pwm_ENB.ChangeDutyCycle(rightspeed)

#小车原地右转
def spin_right(leftspeed, rightspeed):
    GPIO.output(IN1, GPIO.HIGH)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.HIGH)
    pwm_ENA.ChangeDutyCycle(leftspeed)
    pwm_ENB.ChangeDutyCycle(rightspeed)

#小车停止	
def brake():
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.LOW)

def track_task(stop_event):
    try:
        init()
        time.sleep(2)
        while not stop_event.is_set():
            TrackSensorLeftValue1  = GPIO.input(TrackSensorLeftPin1)
            TrackSensorLeftValue2  = GPIO.input(TrackSensorLeftPin2)
            TrackSensorRightValue1 = GPIO.input(TrackSensorRightPin1)
            TrackSensorRightValue2 = GPIO.input(TrackSensorRightPin2)

            if (TrackSensorLeftValue1 == False or TrackSensorLeftValue2 == False) and TrackSensorRightValue2 == False:
                spin_right(40, 40)
                time.sleep(0.08)
 
            elif TrackSensorLeftValue1 == False and (TrackSensorRightValue1 == False or TrackSensorRightValue2 == False):
                spin_left(35, 35)
                time.sleep(0.08)
  
            elif TrackSensorLeftValue1 == False:
                spin_left(30, 30)

            elif TrackSensorRightValue2 == False:
                spin_right(30, 30)
   
            elif TrackSensorLeftValue2 == False and TrackSensorRightValue1 == True:
                left(0, 35)
   
            elif TrackSensorLeftValue2 == True and TrackSensorRightValue1 == False:
                right(35, 0)

            elif TrackSensorLeftValue2 == False and TrackSensorRightValue1 == False:
                run(40, 40)
    except KeyboardInterrupt:
        print("循迹被手动终止")
    finally:
        pwm_ENA.stop()
        pwm_ENB.stop()
        GPIO.cleanup()
        print("GPIO 资源已清理")

def track_interrupt_listener(stop_event):
    while not stop_event.is_set():
        cmd = listen_and_recognize()
        if not cmd:
            continue
        if "中断" in cmd:
            print("检测到语音中断指令")
            stop_event.set()

def handle_tracking_with_interrupt():
    stop_event = threading.Event()

    t1 = threading.Thread(target=track_task, args=(stop_event,))
    t2 = threading.Thread(target=track_interrupt_listener, args=(stop_event,))

    t1.start()
    t2.start()

    t1.join()
    t2.join()

    print("循迹模式结束，回到主流程")
