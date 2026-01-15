# _*_ coding: UTF-8 _*_
# @author: qianqianzyk 
# 同步式语音识别控制系统
import os
import re
import time
import wave
import string
import numpy as np
import subprocess
import traceback
from dotenv import load_dotenv
from aip import AipSpeech
from tongyi_qianwen import TongyiQianwen
from voice_synthesis import speak_text_file, play_mp3

# 初始化环境
load_dotenv()
os.makedirs('voice_recognition', exist_ok=True)
os.makedirs('voice_synthesis', exist_ok=True)

# 常量配置
VOICE_RECOGNITION_DIR = os.path.join(os.getcwd(), 'voice_recognition')
VOICE_SYNTHESIS_DIR = os.path.join(os.getcwd(), 'voice_synthesis')
AUDIO_FILE = os.path.join(VOICE_RECOGNITION_DIR, 'demo.wav')

# 音频参数常量
CHANNELS = 1          # 单声道录音
SAMPLE_RATE = 16000   # 16kHz采样率
SILENCE_THRESHOLD = 1200   # 静音检测阈值(0-32768)
MIN_VOICE_DURATION = 0.8  # 最小有效语音时长(秒)
MAX_SILENCE_DURATION = 1   # 最大允许静音间隔(秒)
CHUNK_DURATION = 0.02  # 每块音频时长20ms
CHUNK_SIZE = int(SAMPLE_RATE * CHUNK_DURATION) * 2  # 每块大小(样本数*2字节)

# 百度语音API配置常量
APP_ID = os.getenv('BD_APP_ID')
API_KEY = os.getenv('BD_API_KEY')
SECRET_KEY = os.getenv('BD_SECRET_KEY')

bd_client = AipSpeech(APP_ID, API_KEY, SECRET_KEY)
tyqw = TongyiQianwen()

# 语音文件映射表
VOICE_MAPPING = {
    "收到": "received",
    "中断": "interrupt",
    "进入人脸追踪模式": "enter_face_tracking",
    "发现目标": "target_found",
    "跳舞": "dance",
    "收到正在查询今日天气": "weather",
    "循迹": "tank_start",
    "大家好呀": "hello",
    "欢迎再次使用哦": "bye",
}

def remove_file(path):
    try:
        if os.path.exists(path):
            os.remove(path)
    except Exception as e:
        print(f"[清理失败] {path}: {e}")

def remove_punctuation(text):
    english_punct = string.punctuation
    chinese_punct = '，。！？【】（）《》“”‘’：；、——…￥'
    all_punct = english_punct + chinese_punct
    pattern = f"[{re.escape(all_punct)}]"
    return re.sub(pattern, "", text)

def listen_and_recognize():
    ok = AudioRecorder.record_with_voice_alsa(AUDIO_FILE)
    if not ok:
        return None

    text = SpeechRecognizer.recognize_speech(AUDIO_FILE)
    if not text:
        return None

    return remove_punctuation(text)

class AudioRecorder:
    """音频录制功能类"""
    @staticmethod
    def record_with_voice_alsa(output_file):
        process = None
        try:
            # 构建arecord命令
            cmd = [
                'arecord',  # 录音命令
                '-D', 'plughw:2,0',  # 可以使用arecord -l查看
                '-f', 'S16_LE',  # 16位小端格式
                '-r', str(SAMPLE_RATE),  # 采样率
                '-c', str(CHANNELS),  # 声道数
                '-t', 'raw'  # 原始数据格式
            ]
            # 启动子进程进行录音
            process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
            
            print("\n等待语音输入...", end='', flush=True)
            
            # 计算音频块参数
            voice_chunks = []  # 存储有效语音数据块
            silent_chunks = 0  # 静音块计数器
            voice_detected = False  # 语音检测标志
            
        
            while True:
                # 读取音频数据块
                data = process.stdout.read(CHUNK_SIZE)
                if not data:
                    continue

                # 将音频数据转换为numpy数组
                audio_data = np.frombuffer(data, dtype=np.int16)
                # 计算当前块的最大振幅
                current_max = np.max(np.abs(audio_data))

                # 语音活动检测
                if current_max > SILENCE_THRESHOLD:
                    silent_chunks = 0
                    voice_chunks.append(data)
                    if not voice_detected:
                        print("\n检测到语音开始...", end='', flush=True)
                        voice_detected = True
                elif voice_detected:
                    silent_chunks += 1
                    voice_chunks.append(data)
                    if silent_chunks > int(MAX_SILENCE_DURATION / CHUNK_DURATION):
                        break

            # 检查是否满足最小录音时长要求
            if len(voice_chunks) >= int(MIN_VOICE_DURATION / CHUNK_DURATION):
                with wave.open(output_file, 'wb') as wf:
                    wf.setnchannels(CHANNELS)
                    wf.setsampwidth(2)
                    wf.setframerate(SAMPLE_RATE)
                    wf.writeframes(b''.join(voice_chunks))
                print("录音成功")
                return True
            else:
                print("语音太短，已忽略")
                return False
                
        except Exception as e:
            print(f"\n录音异常: {str(e)}")
            traceback.print_exc()
            return False
        finally:
            if process:
                process.terminate()
                process.wait()

class SpeechRecognizer:
    """语音识别功能类"""
    @staticmethod
    def recognize_speech(audio_file):
        try:
            # 读取音频文件
            with open(audio_file, 'rb') as f:
                audio_data = f.read()
            
            print("识别中...", end='', flush=True)
            # 调用百度语音识别API
            result = bd_client.asr(audio_data, 'wav', 16000, {'dev_pid': 1536})
            
            # 解析识别结果
            if result.get('err_msg') == 'success.':
                text = result['result'][0]
                print(f"识别成功: {text}")
                return text
            else:
                print("识别失败")
                return None
                
        except Exception as e:
            print(f"\n识别异常: {str(e)}")
            traceback.print_exc()
            return None

class CommandProcessor:
    """指令处理功能类"""
    @staticmethod
    def process_command(text):
        try:
            answer = tyqw.ask(text[2:] if text.startswith("小车") else text)
            print(f"AI回复: {answer}")

            # 语音合成并播放回复
            response_path = os.path.join(VOICE_SYNTHESIS_DIR, 'ai_response.txt')
            os.makedirs(os.path.dirname(response_path), exist_ok=True)
        
            with open(response_path, 'w', encoding='utf-8') as f:
                f.write(answer)

            mp3_path = speak_text_file(response_path)
            if mp3_path and os.path.exists(mp3_path):
                play_mp3(mp3_path)
            else:
                print("语音合成失败")
        except Exception as e:
            print(f"\n[AI处理异常] {e}")
            traceback.print_exc()

    @staticmethod
    def play_response(text):
        """播放指定文本对应的 MP3 音频。若无预置音频，则自动合成"""
        try:
            filename = VOICE_MAPPING.get(text)
            if not filename:
                filename = str(int(time.time() * 1000))
            mp3_path = os.path.join(VOICE_SYNTHESIS_DIR, f"{filename}.mp3")

            if not os.path.exists(mp3_path):
                txt_path = os.path.join(VOICE_SYNTHESIS_DIR, f"{filename}.txt")
                with open(txt_path, 'w', encoding='utf-8') as f:
                    f.write(text)

                mp3_path = speak_text_file(txt_path)

            if mp3_path and os.path.exists(mp3_path):
                play_mp3(mp3_path)
            else:
                print(f"[错误] 无法播放语音：{text}")

        except Exception as e:
            print(f"[play_response 错误] {e}")
            traceback.print_exc()

__all__ = [
    'AUDIO_FILE',
    'VOICE_MAPPING',
    'remove_file',
    'remove_punctuation',
    'listen_and_recognize',
    'AudioRecorder',
    'SpeechRecognizer',
    'CommandProcessor'
]