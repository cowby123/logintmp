import requests
import json
import sys

def reg(ip):
    data = {
	'UserName':'hitmone0',
	'Password':  'hitmone1',
	'ChineseName': '腿骨',
	'Email': 'hitmone0@gmail.com',
	'Address':  '嘉義縣民雄鄉',
	'Phone': '0912345678',
	'LastLogin': 'NOW'
    }
    ## headers中添加上content-type这个参数，指定为json格式
    headers = {'Content-Type': 'application/json'}

    ## post的时候，将data字典形式的参数用json包转换成json格式。
    response = requests.post(url='http://' + ip + "/api/userreg", headers=headers, data=json.dumps(data))
    print(response.text)

f = open('ip','r')
ip = f.read()
reg(ip) 