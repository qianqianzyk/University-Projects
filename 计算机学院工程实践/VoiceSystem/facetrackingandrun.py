import cv2
import time
import RPi.GPIO as GPIO

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

    def setpid(p, i, d):
        self.pid = PID.PositionalPID(p, i, d)

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


GPIO.setmode(GPIO.BCM)#Set GPIO Coding: BCM
GPIO.setwarnings(False)#Ignore Warning

#Define Pin
HeadSvPin = 23
CmrXPin = 11
CmrYPin = 9

#Initialize Servo
print('Initializing Servo')
cmrxSv = Angle(CmrXPin)
cmrySv = Angle(CmrYPin, -45, hbound = 28)

#Initialize motor
print("Initialize motor")
init()

#Prepare Video Stream
print('Preparing Videw Stream')
frame_x = 160
frame_y = 120
cap = cv2.VideoCapture(0)
cap.set(3, 320)
cap.set(4, 240)
cap.set(5, 15)
if not cap.isOpened():
    print("Error: Could not open camera.")
    exit()

#Load Module
print('Load Module')
face_cascade = cv2.CascadeClassifier('haarcascade_frontalface_default.xml')

print('Camera Preheating')
time.sleep(2)#Preheating
cmrxSv.set(0)
cmrySv.set(-60)
print('Ready')
#time_cnt = 0#Time Count
threshold_turn = 40   # 转向阈值
threshold_forward = 15  # 前进阈值
turn_speed = 20
straight_speed = 10

try:
    while True:#Main Circle
        time.sleep(0.05)
        ret, frame = cap.read()
        if not ret:
            continue
        gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
        #Get Faces
        faces = face_cascade.detectMultiScale(gray, scaleFactor=1.1, minNeighbors=5)
        #Draw Frame (Only one face)
        if len(faces) > 0:
        #    time_cnt += 1
            (x, y, w, h) = faces[0]
            #cv2.rectangle(frame, (x, y), (x+w, y+h), (255, 0, 0), 2)
            ct_x = x + w//2
            ct_y = y + h//2
            cv2.circle(frame, (ct_x, ct_y), 2, (0, 255, 0), -1)
            if (ct_x - frame_x > 60):
                cmrxSv.add(-6)
            elif (frame_x - ct_x > 60):
                cmrxSv.add(6)
            elif ct_x - frame_x > 20:
                cmrxSv.add(-3)
            elif frame_x - ct_x > 20:
                cmrxSv.add(3)
            if (ct_y - frame_y > 60):
                cmrySv.add(-5)
            elif (frame_y - ct_y > 60):
                cmrySv.add(5)
            elif ct_y - frame_y > 20:
                cmrySv.add(-2)
            elif frame_y - ct_y > 20:
                cmrySv.add(2)
            print(ct_x, ',', ct_y)

            # 电机控制
            offset_x = ct_x - frame_x


            if abs(offset_x) > threshold_turn:
                if offset_x > 0:
                    # 人脸在右侧，右转
                    right(turn_speed, turn_speed)
                    print("右转")
                else:
                    # 人脸在左侧，左转
                    left(turn_speed, turn_speed)
                    print("左转")
            elif abs(offset_x) > threshold_forward:
                # 偏移不大，微调前进
                run(straight_speed, straight_speed)
                print("前进")
            else:
                # 人脸居中，保持前进
                run(straight_speed, straight_speed)
                print("直行")
        else:
            # 未检测到人脸，停止
            brake()
            print("停止")
        #else:
        #    time_cnt = 0
        #if time_cnt == 60:
        #    cv2.imwrite('img\target.jpg', frame)

        cv2.imshow('Face Tracking', gray)
        if cv2.waitKey(1) & 0xFF == ord('q'):
            break

except KeyboardInterrupt:
    pass

finally:
    #Close
    cv2.destroyAllWindows()
    cap.release()
    GPIO.cleanup()
    print('Program Exit')
