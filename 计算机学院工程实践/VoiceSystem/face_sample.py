# _*_ coding: UTF-8 _*_
# 人脸识别测试样例
import os
import threading
import time
from voice_system import CommandProcessor, listen_and_recognize
from mail_send import send_alert_email

IMAGE_DIR = os.path.join(os.getcwd(), "img")
ALERT_IMAGE_PATH = os.path.join(IMAGE_DIR, "target.jpg")
ALERT_RECEIVER = "zyk18868377764@163.com"

def face_tracking_task(stop_event):
    while not stop_event.is_set():
        print("识别人脸中...")
        time.sleep(2)

        CommandProcessor.play_response("发现目标")
        time.sleep(1)

        if send_alert_email(ALERT_RECEIVER, ALERT_IMAGE_PATH):
            print("邮件已成功发送")
        else:
            print("邮件发送失败")

        time.sleep(1)
    print("人脸识别线程已退出")

def interrupt_listen_task(stop_event):
    while not stop_event.is_set():
        cmd = listen_and_recognize()
        if not cmd:
            continue
        if "中断" in cmd:
            print("检测到停止指令，中断追踪")
            stop_event.set()

def handle_face_tracking_with_interrupt():
    stop_event = threading.Event()

    t1 = threading.Thread(target=face_tracking_task, args=(stop_event,))
    t2 = threading.Thread(target=interrupt_listen_task, args=(stop_event,))

    t1.start()
    t2.start()

    t1.join()
    t2.join()

    print("追踪模式结束，回到主流程")
