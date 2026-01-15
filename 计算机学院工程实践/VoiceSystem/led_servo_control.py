# _*_ coding: UTF-8 _*_
# 小车七彩灯控制模块
import RPi.GPIO as GPIO
import time

# RGB 灯与舵机引脚定义
LED_R = 22
LED_G = 27
LED_B = 24
ServoPin = 23

# 初始化 GPIO
def init_rgb_servo():
    global pwm_servo
    GPIO.setmode(GPIO.BCM)
    GPIO.setwarnings(False)

    GPIO.setup(LED_R, GPIO.OUT)
    GPIO.setup(LED_G, GPIO.OUT)
    GPIO.setup(LED_B, GPIO.OUT)
    GPIO.setup(ServoPin, GPIO.OUT)

    pwm_servo = GPIO.PWM(ServoPin, 50)
    pwm_servo.start(0)

# 七彩灯控制函数，根据位置设置不同颜色
def color_light(pos):
    if pos > 150:
        GPIO.output(LED_R, GPIO.HIGH)
        GPIO.output(LED_G, GPIO.LOW)
        GPIO.output(LED_B, GPIO.LOW)
    elif pos > 125:
        GPIO.output(LED_R, GPIO.LOW)
        GPIO.output(LED_G, GPIO.HIGH)
        GPIO.output(LED_B, GPIO.LOW)
    elif pos > 100:
        GPIO.output(LED_R, GPIO.LOW)
        GPIO.output(LED_G, GPIO.LOW)
        GPIO.output(LED_B, GPIO.HIGH)
    elif pos > 75:
        GPIO.output(LED_R, GPIO.HIGH)
        GPIO.output(LED_G, GPIO.HIGH)
        GPIO.output(LED_B, GPIO.LOW)
    elif pos > 50:
        GPIO.output(LED_R, GPIO.LOW)
        GPIO.output(LED_G, GPIO.HIGH)
        GPIO.output(LED_B, GPIO.HIGH)
    elif pos > 25:
        GPIO.output(LED_R, GPIO.HIGH)
        GPIO.output(LED_G, GPIO.LOW)
        GPIO.output(LED_B, GPIO.HIGH)
    elif pos > 0:
        GPIO.output(LED_R, GPIO.HIGH)
        GPIO.output(LED_G, GPIO.HIGH)
        GPIO.output(LED_B, GPIO.HIGH)
    else:
        GPIO.output(LED_R, GPIO.LOW)
        GPIO.output(LED_G, GPIO.LOW)
        GPIO.output(LED_B, GPIO.LOW)

# 舵机归中
def center_servo():
    pwm_servo.ChangeDutyCycle(2.5 + 10 * 90 / 180)
    time.sleep(0.5)

# 舵机归中并执行七彩灯扫描
def center_and_sweep_lights():
    center_servo()
    for pos in range(181):
        pwm_servo.ChangeDutyCycle(2.5 + 10 * pos / 180)
        color_light(pos)
        time.sleep(0.009)
    for pos in reversed(range(181)):
        pwm_servo.ChangeDutyCycle(2.5 + 10 * pos / 180)
        color_light(pos)
        time.sleep(0.009)

# 清理资源
def clean_up():
    pwm_servo.stop()
    GPIO.cleanup()
