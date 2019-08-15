import os
import platform
import sys


if len(sys.argv)!=2:
    print('Wrong version, please try again!')
    print('Here is a right example:\npython auto.py v0.0.3')
    exit(-1)
version = sys.argv[1]
if platform.system() == 'Windows':
    p = os.popen('echo %GOPATH%')
    re = f'set GOPATH={os.getcwd()};'
    x = p.read()[:-1]
    re += x
else:
    re = 'export GOPATH=$(dirname $(readlink -f $0))'
re += f'&&set CGO_ENABLED=0&&set GOOS=linux&&go build -a -installsuffix cgo -o bin/main .&&docker build --no-cache -t licotek/magicscene-miot-adapter-service-prod:{version} .&&docker push licotek/magicscene-miot-adapter-service-prod:{version}'
os.system(re)
