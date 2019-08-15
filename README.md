# MIOT协议适配器

## 初始化（第一次运行本项目）
```
python initialize.py
```
## 本地运行（建议使用IDE来debug）
```
go build main.go -o miot.exe
```
## 编译并上传到docker hub
先登录docker hub的账号
```
docker login
```
然后运行py
```
python auto.py {tag}
```
比如：
```
python auto.py v1.0.2
python auto.py debug
```
