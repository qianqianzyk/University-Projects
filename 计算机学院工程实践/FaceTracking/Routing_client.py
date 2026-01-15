import socket
import cv2
import mediapipe as mp
import numpy as np
import time

class connect_Raspberry():
    def __init__(self,host,port):
        print("Client Start")
        self.mySocket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        try:
            self.mySocket.connect((host, port))  #Connect
            print("Connection Success")
        except:
            print('Connection Fail')

    def send(self, words):
        # 发送消息
        msg = words
        # 编码发送
        self.mySocket.send(msg.encode("utf-8"))
        # print("成功发送消息")

    def close(self):
        self.mySocket.close()
        print("Connection Interupted\n")
        exit()


#检测脸部
mp_face_detection = mp.solutions.face_detection
mp_drawing = mp.solutions.drawing_utils

#通信传输
RaspConnection = connect_Raspberry('192.168.50.1', 1628)

#Camera Set
capture = cv2.VideoCapture('http://192.168.50.1:2030/?action=stream')


threshold_turn = 40   # 转向阈值
threshold_forward = 15  # 前进阈值
turn_speed = 10
straight_speed = 10
area_too_close = 18000   # 距离过近阈值
area_too_far = 2500      # 距离过远阈值

try:
    while(True):
        time.sleep(0.01)
        ref, frame = capture.read()
        if not ref:
            continue
        h,w,_ = np.shape(frame)
        image = cv2.cvtColor(frame,cv2.COLOR_BGR2RGB)

        #Face Detect
        with mp_face_detection.FaceDetection(
            model_selection=0,
            min_detection_confidence=0.9
        ) as face_detection:

            results = face_detection.process(image)

            if results.detections:
                if len(results.detections) > 0:
                    detection = results.detections[0]
                    box=detection.location_data.relative_bounding_box
                    #cx,cy,cw,ch=box
                    cx=box.xmin
                    cy=box.ymin
                    cw=box.width
                    ch=box.height
                    cv2.circle(frame, (int(w/2), int(h/2)), 1, (0, 0, 255), -1)
                    cv2.rectangle(frame, (int(cx*w) , int(cy*h)), (int((cx+cw)*w) , int((cy+ch)*h)),(0, 255, 0), 2)
                    cv2.circle(frame, (int((cx+cw/2)*w), int((cy+ch/2)*h)), 2, (255, 0, 0), -1)
                    # 电机控制
                    ct_x = h//2#屏幕中心点
                    face_cet_x = cx*h+(cw*h)/2#脸中心点
                    offset_x = ct_x - face_cet_x
                    face_area = (cw*w) * (ch*h)
                    motocontrol="" #00前进01后退10左转11右转
                    if face_area > area_too_close:
                        # 人脸太大，距离太近，停止
                        motocontrol="01"
                        # brake()
                        print("距离太近，后退")
                    elif face_area < area_too_far:
                        # 人脸太小，距离太远，可以选择前进或加速
                        if abs(offset_x) > threshold_turn:
                            if offset_x > 0:
                                motocontrol="10"
                                # right(turn_speed, turn_speed)
                                print("左转（远）")
                            else:
                                motocontrol="11"
                                # left(turn_speed, turn_speed)
                                print("右转（远）")
                        else:
                            motocontrol="00"
                            # run(straight_speed + 10, straight_speed + 10)
                            print("加速前进（远）")
                    else:
                        # 距离合适，正常跟随
                        if abs(offset_x) > threshold_turn:
                            if offset_x > 0:
                                motocontrol="10"
                                # right(turn_speed, turn_speed)
                                print("左转")
                            else:
                                motocontrol="11"
                                # left(turn_speed, turn_speed)
                                print("右转")
                        elif abs(offset_x) > threshold_forward:
                            motocontrol="00"
                            # run(straight_speed, straight_speed)
                            print("前进")
                        else:
                            motocontrol="00"
                            # run(straight_speed, straight_speed)
                            print("直行")                    
                #控制云台
                lctx = int((cx+cw/2)*w)
                lcty = int((cy+ch/2)*h)
                if (0 <= lctx <= 320 and 0 <= lcty <= 280):
                    msg = str(lctx) + " " + str(lcty) + ' ' + motocontrol
                    RaspConnection.send(msg)
            RaspConnection.send("none")

        cv2.imshow("video",frame)
        if cv2.waitKey(1) & 0xff == ord('q'):
            break

finally:
    print("Video Detection Done")
    cv2.destroyAllWindows()

