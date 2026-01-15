# _*_ coding: UTF-8 _*_
# 小车控制模块
import RPi.GPIO as GPIO
import time

# 电机引脚定义
IN1 = 20  
IN2 = 21  
IN3 = 19  
IN4 = 26  
ENA = 16  
ENB = 13  

pwm_ENA = None
pwm_ENB = None

# 小车摄像头Y舵机引脚定义
CmrXPin = 11
CmrYPin = 9

# 摄像头云台类
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


# 初始化电机引脚及 PWM
def init_car():
    global pwm_ENA, pwm_ENB
    GPIO.setmode(GPIO.BCM)            # 设置引脚编号模式为 BCM
    GPIO.setwarnings(False)           # 关闭 GPIO 警告
    # 初始化摄像头云台
    # global cmrSv
    # cmrSv = CmrSv(CmrXPin, CmrYPin, y_max_bound = 28)

    # 设置电机控制引脚为输出模式
    GPIO.setup(ENA, GPIO.OUT, initial=GPIO.HIGH)
    GPIO.setup(IN1, GPIO.OUT, initial=GPIO.LOW)
    GPIO.setup(IN2, GPIO.OUT, initial=GPIO.LOW)
    GPIO.setup(ENB, GPIO.OUT, initial=GPIO.HIGH)
    GPIO.setup(IN3, GPIO.OUT, initial=GPIO.LOW)
    GPIO.setup(IN4, GPIO.OUT, initial=GPIO.LOW)

    # 初始化 PWM 控制，频率设为 2000Hz
    pwm_ENA = GPIO.PWM(ENA, 2000)
    pwm_ENB = GPIO.PWM(ENB, 2000)
    pwm_ENA.start(0)
    pwm_ENB.start(0)

# 小车前进
def car_run(delaytime=1):
    GPIO.output(IN1, GPIO.HIGH)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.HIGH)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(80)
    pwm_ENB.ChangeDutyCycle(80)
    time.sleep(delaytime)

# 小车后退
def car_back(delaytime=1):
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.HIGH)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.HIGH)
    pwm_ENA.ChangeDutyCycle(80)
    pwm_ENB.ChangeDutyCycle(80)
    time.sleep(delaytime)

# 小车左转
def car_left(delaytime=1):
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.HIGH)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(80)
    pwm_ENB.ChangeDutyCycle(80)
    time.sleep(delaytime)

# 小车右转
def car_right(delaytime=1):
    GPIO.output(IN1, GPIO.HIGH)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(80)
    pwm_ENB.ChangeDutyCycle(80)
    time.sleep(delaytime)

# 小车原地左转
def car_spin_left(delaytime=1):
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.HIGH)
    GPIO.output(IN3, GPIO.HIGH)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(80)
    pwm_ENB.ChangeDutyCycle(80)
    time.sleep(delaytime)

# 小车原地右转
def car_spin_right(delaytime=1):
    GPIO.output(IN1, GPIO.HIGH)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.HIGH)
    pwm_ENA.ChangeDutyCycle(80)
    pwm_ENB.ChangeDutyCycle(80)
    time.sleep(delaytime)

# 小车刹车
def car_brake(delaytime=1):
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(0)
    pwm_ENB.ChangeDutyCycle(0)
    time.sleep(delaytime)

# 小车跳舞
def car_dance():
    time.sleep(1)
    # 设置小车XY舵机
    GPIO.setup(CmrYPin, GPIO.OUT)
    GPIO.setup(CmrXPin, GPIO.OUT)
    #设置pwm引脚和频率为50hz
    pwm_servo = GPIO.PWM(CmrYPin, 50)
    pwm_servo.start(2.5)
    pwm_servox = GPIO.PWM(CmrXPin, 50)
    pwm_servox.start(2.5)
    for _ in range(6):
        car_spin_left(0.2)
        time.sleep(0.1)
        car_spin_right(0.2)
    for _ in range(4):
        car_run(0.2)
        time.sleep(0.1)
        car_back(0.2)
    car_brake(0.5)
    car_spin_right(4)
    car_spin_left(4)
    for _ in range(3):
        car_left(0.2)
        time.sleep(0.1)
        car_right(0.2)
    for _ in range(4):
        car_spin_left(0.2)
        time.sleep(0.1)
        car_spin_right(0.2)
    car_brake(1)
    for _ in range (4):
        for pos in range(91):
            pwm_servo.ChangeDutyCycle(2.5 + 10 * pos/180)
            
            time.sleep(0.009) 
        for pos in reversed(range(91)):
            pwm_servo.ChangeDutyCycle(2.5 + 10 * pos/180)
            
            time.sleep(0.009)
    pwm_servo.stop()
    pwm_servox.stop()
    
# 清理小车资源
def clean_up():
    pwm_ENA.stop()
    pwm_ENB.stop()
    GPIO.cleanup()