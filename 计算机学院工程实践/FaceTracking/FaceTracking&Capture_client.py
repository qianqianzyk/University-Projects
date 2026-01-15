import socket
import cv2
import mediapipe as mp
import numpy as np
import time
import struct

class connect_Raspberry():
    def __init__(self,host,port):
        print("Client Start")
        self.skt = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.host = host
        self.port = port
        try:
            self.skt.connect((host, port))  #Connect
            print("Connection Success")
        except:
            print('Connection Fail')

    def send(self, words):
        # 发送消息
        msg = words
        # 编码发送
        self.skt.send(msg.encode("utf-8"))
        # print("成功发送消息")

    def close(self):
        self.skt.close()
        print("Connection Closed\n")

def trans_img(skt, img_addr):
    with open(img_addr, 'rb') as f:
        image_data = f.read()
        file_size = len(image_data)
        skt.send(struct.pack('!I', file_size))
        skt.sendall(image_data)

#Detective Module
mp_face_detection = mp.solutions.face_detection
mp_drawing = mp.solutions.drawing_utils

#Socket Connection
RaspConnection = connect_Raspberry('192.168.50.1', 1629)

#Camera Set
capture = cv2.VideoCapture('http://192.168.50.1:2031/?action=stream')

capture.set(cv2.CAP_PROP_BUFFERSIZE, 1)
time_cnt = 0

try:
    while(True):
        #time.sleep(0.01)
        ref, frame = capture.read()
        if not ref:
            continue
        h,w,_ = np.shape(frame)
        image = cv2.cvtColor(frame,cv2.COLOR_BGR2RGB)

        #Face Detect
        with mp_face_detection.FaceDetection(
            model_selection=0,
            min_detection_confidence=0.85
        ) as face_detection:

            results = face_detection.process(image)

            if results.detections:
                if len(results.detections) > 0:
                    time_cnt += 1
                    if time_cnt == 60:
                        print('Face Detected, Preparing to Capture')
                        cv2.imwrite('enermy_img.jpg', frame)
                        RaspConnection.send('Stop')
                        print('Capture Comleted')
                        break
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
                
                #Servo Control Message
                lctx = int((cx+cw/2)*w)
                lcty = int((cy+ch/2)*h)
                if (0 <= lctx <= 320 and 0 <= lcty <= 280):
                    msg = str(lctx) + " " + str(lcty)
                    RaspConnection.send(msg)

        cv2.imshow("video",frame)
        if cv2.waitKey(1) & 0xff == ord('q'):
            break
    trans_img(RaspConnection.skt, '/enermy_img.jpg')
    print('Image Transported')
except:
    pass
finally:
    print("Video Detection Done")
    capture.release()
    cv2.destroyAllWindows()
    RaspConnection.close()

