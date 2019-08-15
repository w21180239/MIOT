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
## 一些重要的常量
在*src/constant/constant.go*里
```	
    Output_debug_information   //输出debug信息
    Local_debug                //进行本地debug，在服务器上使用的事内部网址，在获取devices等信息时不用登陆，但是在本地进行debug的时候使用的是外部网址，需要登陆，所以需要和下面这个cookies配套使用
    Cookie                     //通过在Chrome登陆过后查看cookie并填到这里
```