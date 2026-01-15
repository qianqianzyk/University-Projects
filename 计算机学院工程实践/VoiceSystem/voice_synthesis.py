# -*- coding: UTF-8 -*-
# @author: qianqianzyk 
# 百度云语音合成
from __future__ import print_function
import sys
from dotenv import load_dotenv
from aip import AipSpeech
import os
import subprocess

IS_PY3 = sys.version_info[0] == 3

load_dotenv()

APP_ID = os.getenv('BD_APP_ID')
API_KEY = os.getenv('BD_API_KEY')
SECRET_KEY = os.getenv('BD_SECRET_KEY')

client = AipSpeech(APP_ID, API_KEY, SECRET_KEY)

def read_text(file_path):
    try:
        if IS_PY3:
            with open(file_path, 'r', encoding='utf-8') as f:
                return f.read().strip()
        else:
            with open(file_path, 'rb') as f:
                return f.read().decode('utf-8').strip()
    except UnicodeDecodeError:
        # 如果UTF-8解码失败，尝试GBK编码
        if IS_PY3:
            with open(file_path, 'r', encoding='gbk') as f:
                return f.read().strip()
        else:
            with open(file_path, 'rb') as f:
                return f.read().decode('gbk').strip()

def speak_text_file(file_path):
    if not os.path.exists(file_path):
        print(u'文件不存在: {}'.format(file_path))
        return None
    
    text = read_text(file_path)

    if not text:
        print(u'文件内容为空')
        return None

    if len(text.encode('utf-8')) > 1024:
        print(u'文本内容超过1024字节限制')
        return None

    result = client.synthesis(
        text, 
        'zh',   # 中文
        1,      # 普通话
        {
            'vol': 10,  # 音量 0-15
            'per': 4   # 发音人 0为女声，1为男声，3为度逍遥，4为度丫丫
        }
    )

    if not isinstance(result, dict):
        dir_path = os.path.dirname(file_path)
        try:
            os.makedirs(dir_path)
        except OSError:
            pass
        
        base_name = os.path.splitext(os.path.basename(file_path))[0]
        mp3_path = os.path.join(dir_path, '{}.mp3'.format(base_name))
        
        with open(mp3_path, 'wb') as f:
            f.write(result)
        
        print(u'语音合成成功，文件保存在: {}'.format(mp3_path))
        os.remove(file_path)
        print(u'已删除文本文件: {}'.format(file_path))
        return mp3_path
    else:
        print(u'语音合成失败，错误信息:', result)
        return None

def play_mp3(path):
    if not os.path.exists(path):
        print(u"播放文件不存在: {}".format(path))
        return
    
    players = [
        ('mpg123', ['-q']),
        ('omxplayer', ['-o', 'local', '--no-keys'])
    ]
    
    for player, args in players:
        try:
            subprocess.Popen([player] + args + [path])
            break
        except (OSError, subprocess.CalledProcessError):
            continue
    else:
        print(u"未找到可用的播放器")