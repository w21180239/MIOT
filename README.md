# 小度在家协议适配器

## 非容器使用

```bash
cd $GOPATH/src/gitlab.com
git clone git@gitlab.com:LICOTEK/DuerOSAdapter.git
cd DuerOSAdapter
go get -v -d
go run main.go
```

## 容器

运行`docker.sh`可在本地建立镜像并运行，容器采用基于scratch，只包含可执行文件`bin/main`构建的最小镜像。