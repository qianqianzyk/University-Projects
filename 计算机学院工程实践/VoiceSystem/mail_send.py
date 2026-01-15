# -*- coding: UTF-8 -*-
# @author: qianqianzyk 
# QQ邮箱发送功能
import os
import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from email.mime.image import MIMEImage
from email.header import Header
from dotenv import load_dotenv

# 加载环境变量
load_dotenv()

# 全局常量配置
SMTP_SERVER = os.getenv("SMTP_SERVER")
SMTP_PORT = os.getenv("SMTP_PORT")
EMAIL_SENDER = os.getenv("EMAIL_SENDER")
EMAIL_PASSWORD = os.getenv("EMAIL_KEY")
EMAIL_SUBJECT = "⚠️ 发现敌人！"
EMAIL_INFO = "发现一名可疑目标，请及时处理。以下是实时情报图像："

def send_alert_email(to_email, image_path):
    msg = MIMEMultipart("related")
    msg["From"] = Header(EMAIL_SENDER)
    msg["To"] = Header(to_email)
    msg["Subject"] = Header(EMAIL_SUBJECT, "utf-8")

    html = f"""
    <html>
    <body>
        <h2>{EMAIL_SUBJECT}</h2>
        <p>{EMAIL_INFO}</p>
        <p><img src="cid:target_image" style="max-width:600px; height:auto;"></p>
    </body>
    </html>
    """
    msg.attach(MIMEText(html, "html", "utf-8"))

    try:
        with open(image_path, "rb") as f:
            img = MIMEImage(f.read())
            img.add_header("Content-ID", "<target_image>")
            msg.attach(img)

        server = smtplib.SMTP(SMTP_SERVER, SMTP_PORT)
        server.starttls()
        server.login(EMAIL_SENDER, EMAIL_PASSWORD)
        server.sendmail(EMAIL_SENDER, [to_email], msg.as_string())
        server.quit()

        print("情报邮件发送成功")
        return True
    except Exception as e:
        print(f"邮件发送失败：{e}")
        return False