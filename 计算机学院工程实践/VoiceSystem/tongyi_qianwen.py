# -*- coding: UTF-8 -*-
# @author: qianqianzyk 
# 阿里云百炼通义千问智能问答
from dotenv import load_dotenv
import os
import requests
import re

load_dotenv()

class TongyiQianwen:
    def __init__(self):
        self.api_key = os.getenv('TYQW_API_KEY')
        self.endpoint = os.getenv("TYQW_API_ENDPOINT")
        self.model = os.getenv("TYQW_MODEL")

        self.qa_patterns = {
            r"(你是谁|你叫什么|你的名字)": [
                "我是你的智能助手小T,你可以叫我T博士,随时为你解答问题"
            ],
            r"(能做什么|你会什么|功能)": [
                "我能回答各类知识问题，百科全书式问答哦"
            ]
        }
    
    def ask(self, question):
        """智能问答入口"""
        answer = self._match_pattern(question)
        if answer:
            return answer
            
        prompt = self._build_prompt(question)
        try:
            response = requests.post(
                self.endpoint,
                headers={"Authorization": f"Bearer {self.api_key}"},
                json={
                    "model": self.model,
                    "input": {"messages": [
                        {"role": "system", "content": prompt},
                        {"role": "user", "content": question}
                    ]},
                    "parameters": {
                        "temperature": 0.9
                    }
                },
                timeout=10
            )
            return self._format_answer(response.json())
        except Exception as e:
            print(f"API错误: {e}")
            return "网络不稳定，请稍后再试"
    
    def _match_pattern(self, question):
        question = question.lower().strip("？?")
        for pattern in self.qa_patterns:
            if re.search(pattern, question):
                return self.qa_patterns[pattern][0]
        return None
    
    def _build_prompt(self, question):
        """构建动态Prompt"""
        return f"""
        你是一个通用智能助手，需满足：
        1. 回答包括常识、科技、历史等各类问题
        2. 字数严格≤50字,可以适当超出
        3. 对相似问题保持回答一致性
        
        用户问：{question}
        回答："""

    def _format_answer(self, response):
        """格式化API响应"""
        answer = response.get('output', {}).get('text', '').strip()
        answer = re.sub(r'\s+', ' ', answer)
        return answer if answer else "这个问题我需要再学习下"