# -*- coding: UTF-8 -*-
# @author: qianqianzyk
# 百度云手势智能识别
import os
import cv2
import time
import numpy as np
from PIL import Image, ImageDraw, ImageFont
from aip import AipBodyAnalysis
from dotenv import load_dotenv
import demjson3 as demjson

# 手势字典
HAND_GESTURES = {
    'One': '数字一', 'Five': '数字五', 'Fist': '拳头', 'Ok': 'OK',
    'Prayer': '祈祷', 'Congratulation': '作揖', 'Honour': '作别',
    'Heart_single': '比心心', 'Thumb_up': '点赞', 'Thumb_down': 'Diss',
    'ILY': '我爱你', 'Palm_up': '掌心向上', 'Heart_1': '双手比心一',
    'Heart_2': '双手比心二', 'Heart_3': '双手比心三', 'Two': '数字二',
    'Three': '数字三', 'Four': '数字四', 'Six': '数字六', 'Seven': '数字七',
    'Eight': '数字八', 'Nine': '数字九', 'Rock': 'Rock', 'Insult': '竖中指',
    'Face': '脸'
}

load_dotenv()

APP_ID = os.getenv('RT_BD_APP_ID')
API_KEY = os.getenv('RT_BD_API_KEY')
SECRET_KEY = os.getenv('RT_BD_SECRET_KEY')

client = AipBodyAnalysis(APP_ID, API_KEY, SECRET_KEY)

def cv2_img_add_text(img, text, x, y, color=(0, 255, 0), size=20):
    if isinstance(img, np.ndarray):
        img = Image.fromarray(cv2.cvtColor(img, cv2.COLOR_BGR2RGB))
    draw = ImageDraw.Draw(img)
    font = ImageFont.truetype("/usr/share/fonts/truetype/droid/DroidSansFallbackFull.ttf", size)
    draw.text((x, y), text, fill=color, font=font)
    return cv2.cvtColor(np.asarray(img), cv2.COLOR_RGB2BGR)

def detect_gesture(frame):
    success, encoded = cv2.imencode('.jpg', frame)
    if not success:
        return []

    response = client.gesture(encoded.tobytes())
    result = demjson.decode(str(response))
    return result.get('result', [])

def run_gesture_module():
    camera = cv2.VideoCapture(int(os.getenv('CAMERA_INDEX')))
    camera.set(cv2.CAP_PROP_FRAME_WIDTH, 320)
    camera.set(cv2.CAP_PROP_FRAME_HEIGHT, 240)
    camera.set(cv2.CAP_PROP_FPS, 30)
    camera.set(cv2.CAP_PROP_FOURCC, cv2.VideoWriter.fourcc('M', 'J', 'P', 'G'))
    camera.set(cv2.CAP_PROP_BRIGHTNESS, 40)
    camera.set(cv2.CAP_PROP_CONTRAST, 50)

    detected = None

    try:
        while True:
            ret, frame = camera.read()
            if not ret:
                continue

            gestures = detect_gesture(frame)
            found_six = False

            if gestures:
                for item in gestures:
                    classname = item.get('classname')
                    label = HAND_GESTURES.get(classname, classname)
                    left = item.get('left', 0)
                    top = item.get('top', 0)
                    width = item.get('width', 0)
                    height = item.get('height', 0)

                    cv2.rectangle(frame, (left, top), (left + width, top + height), (255, 0, 0), 2)
                    frame = cv2_img_add_text(frame, label, left, top - 25, size=24)
                    print(f"识别到：{label}")

                    if classname == 'Six':
                        found_six = True
                        detected = label
            else:
                print("识别结果：什么也没识别到")
                frame = cv2_img_add_text(frame, "未识别", 30, 30, color=(0, 0, 255), size=30)

            cv2.imshow("手势识别", frame)

            if found_six:
                print("识别到Six，五秒后退出")
                cv2.imshow("手势识别", frame)
                cv2.waitKey(1)
                cv2_img_add_text(frame, "数字六", 10, 10, color=(255, 0, 0), size=20)
                cv2.imshow("手势识别", frame)
                cv2.waitKey(1)
                time.sleep(5)
                break
            if cv2.waitKey(10) & 0xFF == ord('q'):
                print("按下Q退出")
                break

    except KeyboardInterrupt:
        print("\n用户中断操作")
    finally:
        camera.release()
        cv2.destroyAllWindows()

    return detected