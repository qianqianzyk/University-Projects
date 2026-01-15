# -*- coding: UTF-8 -*-
# @author: qianqianzyk
# 百度云图像智能识别
import os
import cv2
import base64
import requests
import json
from dotenv import load_dotenv
import demjson3 as demjson

load_dotenv()

API_KEY = os.getenv('TX_BD_API_KEY')
SECRET_KEY = os.getenv('TX_BD_SECRET_KEY')

# 获取Access Token
def get_access_token():
    url = f"https://aip.baidubce.com/oauth/2.0/token?client_id={API_KEY}&client_secret={SECRET_KEY}&grant_type=client_credentials"
    headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
    }
    response = requests.post(url, headers=headers, data=json.dumps(""))
    if response.status_code == 200:
        return response.json().get("access_token")
    else:
        print("获取access_token失败：", response.text)
        return None

# 摄像头拍照并保存图片
def capture_image(filename="img/captured.jpg"):
    camera = cv2.VideoCapture(int(os.getenv('CAMERA_INDEX')))
    camera.set(cv2.CAP_PROP_FRAME_WIDTH, 320)
    camera.set(cv2.CAP_PROP_FRAME_HEIGHT, 240)
    camera.set(cv2.CAP_PROP_FPS, 30)
    camera.set(cv2.CAP_PROP_FOURCC, cv2.VideoWriter.fourcc('M', 'J', 'P', 'G'))
    camera.set(cv2.CAP_PROP_BRIGHTNESS, 40)
    camera.set(cv2.CAP_PROP_CONTRAST, 50)

    try:
        while True:
            ret, frame = camera.read()
            if not ret:
                continue

            cv2.imshow("拍照窗口", frame)

            key = cv2.waitKey(1)
            if key == 32:  # 空格键
                cv2.imwrite(filename, frame)
                print("已拍照并保存为", filename)
                return filename
            elif key == 27:  # ESC键
                print("用户取消拍照")
                return None

    except KeyboardInterrupt:
        print("\n中断拍照")
    finally:
        camera.release()
        cv2.destroyAllWindows()

# 图像Base64编码
def encode_image_base64(image_path):
    with open(image_path, 'rb') as f:
        return base64.b64encode(f.read()).decode()

# 调用图像识别API
def recognize_image(image_path):
    access_token = get_access_token()
    if not access_token:
        print("无法获取access_token")
        return

    url = f"https://aip.baidubce.com/rest/2.0/image-classify/v2/advanced_general?access_token={access_token}"
    headers = {'Content-Type': 'application/x-www-form-urlencoded'}
    image_base64 = encode_image_base64(image_path)
    data = {
        'image': image_base64,
        'baike_num': 1
    }

    response = requests.post(url, data=data, headers=headers)
    result = demjson.decode(response.text)

    if 'result' in result:
        print("\n识别结果：")
        for item in result['result']:
            print(f"- 关键词: {item['keyword']}，置信度: {item['score']:.2f}")
            if 'baike_info' in item:
                desc = item['baike_info'].get('description', '')
                if desc:
                    print(f"  百科描述：{desc[:100]}...")
    else:
        print("识别失败：", result)