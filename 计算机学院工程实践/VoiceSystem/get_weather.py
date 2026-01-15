# -*- coding: UTF-8 -*-
# @author: qianqianzyk 
# 获取今日天气（默认杭州）
import requests
from bs4 import BeautifulSoup

DEFAULT_CITY_CODE = "101210101"  # 杭州城市代码
BASE_URL = "http://www.weather.com.cn/weather"
USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

def get_hangzhou_weather(city_code=DEFAULT_CITY_CODE):
    url = f"{BASE_URL}/{city_code}.shtml"
    headers = {
        'User-Agent': USER_AGENT
    }

    try:
        # 发送请求
        response = requests.get(url, headers=headers, timeout=10)
        response.raise_for_status()
        response.encoding = 'utf-8'

        # 解析 HTML 页面
        soup = BeautifulSoup(response.text, 'html.parser')
        today_weather = soup.find('ul', class_='t clearfix').find('li')

        # 提取日期
        date = today_weather.find('h1').get_text(strip=True)

        # 提取天气情况
        weather = today_weather.find('p', class_='wea').get_text(strip=True)

        # 提取温度
        temp_tag = today_weather.find('p', class_='tem')
        high_tag = temp_tag.find('span')
        high_temp = high_tag.get_text(strip=True) if high_tag else ''
        low_temp = temp_tag.find('i').get_text(strip=True)
        temperature = f"{low_temp}~{high_temp}" if high_temp else f"{low_temp}"

        # 提取风力风向
        wind_tag = today_weather.find('p', class_='win')
        wind_direction = wind_tag.find('span')['title']
        wind_force = wind_tag.find('i').get_text(strip=True)
        wind_info = f"{wind_direction}{wind_force}"

        # 拼接文本
        text = f"日期：{date}，天气：{weather}，温度：{temperature}，风力：{wind_info}。"

        return text

    except Exception as e:
        print(f"获取天气信息失败: {e}")
        return None