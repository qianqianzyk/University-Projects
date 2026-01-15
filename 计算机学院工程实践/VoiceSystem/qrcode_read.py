# -*- coding: UTF-8 -*-
# @author: qianqianzyk 
# 二维码智能识别
import os
import cv2
from dotenv import load_dotenv
import pyzbar.pyzbar as pyzbar

load_dotenv()

def read_qrcode_text():
    camera = cv2.VideoCapture(int(os.getenv('CAMERA_INDEX')))

    if not camera.isOpened():
        print("无法打开摄像头")
        return None

    # 摄像头参数设置
    camera.set(cv2.CAP_PROP_FRAME_WIDTH, 320)
    camera.set(cv2.CAP_PROP_FRAME_HEIGHT, 240)
    camera.set(cv2.CAP_PROP_FPS, 30)
    camera.set(cv2.CAP_PROP_FOURCC, cv2.VideoWriter.fourcc('M', 'J', 'P', 'G'))
    camera.set(cv2.CAP_PROP_BRIGHTNESS, 40)
    camera.set(cv2.CAP_PROP_CONTRAST, 50)
    # camera.set(cv2.CAP_PROP_EXPOSURE, 156)  # 根据设备兼容性决定是否启用

    try:
        while True:
            ret, frame = camera.read()
            if not ret:
                print("无法读取摄像头图像")
                continue

            gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
            barcodes = pyzbar.decode(gray)

            for barcode in barcodes:
                (x, y, w, h) = barcode.rect
                barcode_data = barcode.data.decode("utf-8")
                barcode_type = barcode.type

                # 绘制二维码边框和信息
                cv2.rectangle(frame, (x, y), (x + w, y + h), (0, 255, 0), 2)
                text = f'{barcode_type}: {barcode_data}'
                cv2.putText(frame, text, (x, y - 10),
                            cv2.FONT_HERSHEY_SIMPLEX, 0.5, (0, 255, 0), 2)

                cv2.imshow("QR Code Scanner", frame)
                print("识别到二维码内容：", barcode_data)

                cv2.waitKey(2000)
                return barcode_data
 
            cv2.imshow("QR Code Scanner", frame)

            if cv2.waitKey(10) & 0xFF == ord('q'):
                break
            
    except KeyboardInterrupt:
        print("\n用户中断操作")
    finally:
        camera.release()
        cv2.destroyAllWindows()

    return None