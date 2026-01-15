# -*- coding: UTF-8 -*-
# 避障控制模块
import RPi.GPIO as GPIO
import time
import threading
from voice_system import listen_and_recognize

IN1 = 20
IN2 = 21
IN3 = 19
IN4 = 26
ENA = 16
ENB = 13
key = 8
EchoPin = 0
TrigPin = 1

GPIO.setmode(GPIO.BCM)
GPIO.setwarnings(False)

def init():
    global pwm_ENA
    global pwm_ENB
    GPIO.setup(ENA, GPIO.OUT, initial=GPIO.HIGH)
    GPIO.setup(IN1, GPIO.OUT, initial=GPIO.LOW)
    GPIO.setup(IN2, GPIO.OUT, initial=GPIO.LOW)
    GPIO.setup(ENB, GPIO.OUT, initial=GPIO.HIGH)
    GPIO.setup(IN3, GPIO.OUT, initial=GPIO.LOW)
    GPIO.setup(IN4, GPIO.OUT, initial=GPIO.LOW)
    GPIO.setup(key, GPIO.IN)
    GPIO.setup(EchoPin, GPIO.IN)
    GPIO.setup(TrigPin, GPIO.OUT)
    pwm_ENA = GPIO.PWM(ENA, 2000)
    pwm_ENB = GPIO.PWM(ENB, 2000)
    pwm_ENA.start(0)
    pwm_ENB.start(0)

def run(leftspeed, rightspeed):
    GPIO.output(IN1, GPIO.HIGH)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.HIGH)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(leftspeed)
    pwm_ENB.ChangeDutyCycle(rightspeed)

def back(leftspeed, rightspeed):
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.HIGH)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.HIGH)
    pwm_ENA.ChangeDutyCycle(leftspeed)
    pwm_ENB.ChangeDutyCycle(rightspeed)

def left(leftspeed, rightspeed):
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.HIGH)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(leftspeed)
    pwm_ENB.ChangeDutyCycle(rightspeed)

def right(leftspeed, rightspeed):
    GPIO.output(IN1, GPIO.HIGH)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(leftspeed)
    pwm_ENB.ChangeDutyCycle(rightspeed)

def spin_left(leftspeed, rightspeed):
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.HIGH)
    GPIO.output(IN3, GPIO.HIGH)
    GPIO.output(IN4, GPIO.LOW)
    pwm_ENA.ChangeDutyCycle(leftspeed)
    pwm_ENB.ChangeDutyCycle(rightspeed)

def spin_right(leftspeed, rightspeed):
    GPIO.output(IN1, GPIO.HIGH)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.HIGH)
    pwm_ENA.ChangeDutyCycle(leftspeed)
    pwm_ENB.ChangeDutyCycle(rightspeed)

def brake():
    GPIO.output(IN1, GPIO.LOW)
    GPIO.output(IN2, GPIO.LOW)
    GPIO.output(IN3, GPIO.LOW)
    GPIO.output(IN4, GPIO.LOW)

def key_scan():
    while GPIO.input(key):
        pass
    while not GPIO.input(key):
        time.sleep(0.01)
        if not GPIO.input(key):
            time.sleep(0.01)
            while not GPIO.input(key):
                pass

# 测距函数
def Distance():
    GPIO.output(TrigPin, GPIO.LOW)
    time.sleep(0.000002)
    GPIO.output(TrigPin, GPIO.HIGH)
    time.sleep(0.000015)
    GPIO.output(TrigPin, GPIO.LOW)

    t3 = time.time()

    while not GPIO.input(EchoPin):
        t4 = time.time()
        if (t4 - t3) > 0.03:
            return -1

    t1 = time.time()
    while GPIO.input(EchoPin):
        t5 = time.time()
        if (t5 - t1) > 0.03:
            return -1

    t2 = time.time()
    time.sleep(0.01)
    return ((t2 - t1) * 340 / 2) * 100

def Distance_test():
    num = 0
    ultrasonic = []
    while num < 5:
        distance = Distance()
        while int(distance) == -1:
            distance = Distance()
            print("Tdistance is %f" % distance)
        while int(distance) >= 500 or int(distance) == 0:
            distance = Distance()
            print("Edistance is %f" % distance)
        ultrasonic.append(distance)
        num = num + 1
        time.sleep(0.01)
    print(ultrasonic)
    distance = (ultrasonic[1] + ultrasonic[2] + ultrasonic[3]) / 3
    print("distance is %f" % distance)
    return distance

# 避障主逻辑
def obstacle_avoidance_task(stop_event):
    try:
        init()
        time.sleep(2)
        while not stop_event.is_set():
            distance = Distance_test()
            if distance > 50:
                run(20, 20)
            elif 30 <= distance <= 50:
                run(15, 15)
            else:
                spin_right(23, 23)
                time.sleep(0.35)
                brake()
                time.sleep(0.001)
                d = Distance_test()
                if d >= 30:
                    run(13, 13)
                else:
                    spin_left(23, 23)
                    time.sleep(0.6)
                    brake()
                    time.sleep(0.001)
                    d = Distance_test()
                    if d >= 30:
                        run(13, 13)
                    else:
                        spin_left(23, 23)
                        time.sleep(0.3)
                        brake()
                        time.sleep(0.001)
    except Exception as e:
        print(f"[避障异常] {e}")
    finally:
        if pwm_ENA: pwm_ENA.stop()
        if pwm_ENB: pwm_ENB.stop()
        GPIO.cleanup()

def obstacle_interrupt_task(stop_event):
    while not stop_event.is_set():
        cmd = listen_and_recognize()
        if not cmd:
            continue
        if "中断" in cmd:
            print("检测到停止指令，中断避障")
            stop_event.set()

def handle_obstacle_avoidance_with_interrupt():
    stop_event = threading.Event()

    t1 = threading.Thread(target=obstacle_avoidance_task, args=(stop_event,))
    t2 = threading.Thread(target=obstacle_interrupt_task, args=(stop_event,))

    t1.start()
    t2.start()

    t1.join()
    t2.join()

    print("避障模式结束，回到主流程")
    
def simple_obstacle_avoidance():
    try:
        init()
        time.sleep(2)
        while True:
            distance = Distance_test()
            if distance > 50:
                run(20, 20)
            elif 30 <= distance <= 50:
                run(15, 15)
            else:
                spin_right(23, 23)
                time.sleep(0.35)
                brake()
                time.sleep(0.001)
                d = Distance_test()
                if d >= 30:
                    run(13, 13)
                else:
                    spin_left(23, 23)
                    time.sleep(0.6)
                    brake()
                    time.sleep(0.001)
                    d = Distance_test()
                    if d >= 30:
                        run(13, 13)
                    else:
                        spin_left(23, 23)
                        time.sleep(0.3)
                        brake()
                        time.sleep(0.001)
    except KeyboardInterrupt:
        print("已手动终止避障")
    finally:
        if pwm_ENA: pwm_ENA.stop()
        if pwm_ENB: pwm_ENB.stop()
        GPIO.cleanup()