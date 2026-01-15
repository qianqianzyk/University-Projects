# _*_ coding: UTF-8 _*_
# 树莓派G1-坦克履带车
# 主程序入口
import time
import avoid
import track
import car_control
import led_servo_control
from voice_system import AudioRecorder, SpeechRecognizer, CommandProcessor, AUDIO_FILE, remove_file, remove_punctuation
from face_sample import handle_face_tracking_with_interrupt
from get_weather import get_hangzhou_weather
from qrcode_read import read_qrcode_text
from gesture_recognition import run_gesture_module
from image_recognition import recognize_image, capture_image

def main():
    """主程序：录音 → 识别 → 本地指令分发或 AI 处理"""
    try:
        while True:
            print("\n" + "=" * 40 + "\n等待新一轮指令…")

            # 录音
            ok = AudioRecorder.record_with_voice_alsa(AUDIO_FILE)
            if not ok:
                time.sleep(1)
                continue

            # 识别指令
            text = SpeechRecognizer.recognize_speech(AUDIO_FILE)
            if not text:
                time.sleep(1)
                continue
            cmd = remove_punctuation(text)

            # 本地指令分发
            if "向前移动" in cmd:
                CommandProcessor.play_response("收到")
                time.sleep(0.5)
                try:
                    car_control.init_car()
                    car_control.car_run()
                except Exception as e:
                    print(f"发生异常：{e}")
                finally:
                    car_control.clean_up()
            
            elif "向后移动" in cmd:
                CommandProcessor.play_response("收到")
                time.sleep(0.5)
                try:
                    car_control.init_car()
                    car_control.car_back()
                except Exception as e:
                    print(f"发生异常：{e}")
                finally:
                    car_control.clean_up()

            elif "向左转向" in cmd:
                CommandProcessor.play_response("收到")
                time.sleep(0.5)
                try:
                    car_control.init_car()
                    car_control.car_left()
                except Exception as e:
                    print(f"发生异常：{e}")
                finally:
                    car_control.clean_up()

            elif "向右转向" in cmd:
                CommandProcessor.play_response("收到")
                time.sleep(0.5)
                try:
                    car_control.init_car()
                    car_control.car_right()
                except Exception as e:
                    print(f"发生异常：{e}")
                finally:
                    car_control.clean_up()

            elif "原地左转" in cmd:
                CommandProcessor.play_response("收到")
                time.sleep(0.5)
                try:
                    car_control.init_car()
                    car_control.car_spin_left()
                except Exception as e:
                    print(f"发生异常：{e}")
                finally:
                    car_control.clean_up()

            elif "原地右转" in cmd:
                CommandProcessor.play_response("收到")
                time.sleep(0.5)
                try:
                    car_control.init_car()
                    car_control.car_spin_right()
                except Exception as e:
                    print(f"发生异常：{e}")
                finally:
                    car_control.clean_up()

            elif "立即停止" in cmd:
                CommandProcessor.play_response("收到")
                time.sleep(0.5)
                try:
                    car_control.init_car()
                    car_control.car_brake()
                except Exception as e:
                    print(f"发生异常：{e}")
                finally:
                    car_control.clean_up()
            
            elif "跟踪敌人" in cmd: 
                CommandProcessor.play_response("进入人脸追踪模式")
                time.sleep(2)
                handle_face_tracking_with_interrupt()

            elif "跳舞" in cmd:
                CommandProcessor.play_response("收到")
                time.sleep(0.5)
                CommandProcessor.play_response("跳舞")
                try:
                    car_control.init_car()
                    car_control.car_dance()
                except Exception as e:
                    print(f"发生异常：{e}")
                finally:
                    car_control.clean_up()

            elif "查询今日天气" in cmd:
                CommandProcessor.play_response("收到正在查询今日天气")
                time.sleep(2)
                weather_info = get_hangzhou_weather()
                CommandProcessor.play_response(weather_info)
                time.sleep(10)

            elif "避开障碍物" in cmd:
                CommandProcessor.play_response("收到")
                time.sleep(0.5)
                avoid.handle_obstacle_avoidance_with_interrupt()

            elif "沿着黑线走" in cmd:
                CommandProcessor.play_response("收到")
                time.sleep(0.5)
                CommandProcessor.play_response("循迹")
                track.handle_tracking_with_interrupt()

            elif "启动二维码识别" in cmd:
                CommandProcessor.play_response("收到")
                time.sleep(0.5)
                text = read_qrcode_text()
                text = f"识别出{len(text)}个字：{text}"

            elif "启动手势识别" in cmd:
                CommandProcessor.play_response("收到")
                time.sleep(0.5)
                run_gesture_module()
                CommandProcessor.play_response("识别到手势六退出循环")
                time.sleep(3)

            elif "启动图像识别" in cmd:
                CommandProcessor.play_response("收到")
                time.sleep(0.5)
                img = capture_image()
                recognize_image(img)

            elif "给大家打个招呼" in cmd:
                CommandProcessor.play_response("大家好呀")
                try:
                    led_servo_control.init_rgb_servo()
                    led_servo_control.center_and_sweep_lights()
                except Exception as e:
                    print(f"发生异常：{e}")
                finally:
                    led_servo_control.clean_up()

            # 如果都不是本地指令，交给大模型来处理并进行音频播放
            elif "小车" in cmd:
                CommandProcessor.process_command(cmd)

            else:
                print("非小车指令，忽略")

            time.sleep(3)

    except KeyboardInterrupt:
        CommandProcessor.play_response("欢迎再次使用哦")
        print("\n收到终止信号,程序退出。")

    finally:
        remove_file(AUDIO_FILE)

if __name__ == "__main__":
    main()