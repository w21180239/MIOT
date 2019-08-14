docker stop licotek/magicscene-miot-adapter-service-prod
docker rmi licotek/magicscene-miot-adapter-service-prod
docker build --no-cache -t licotek/magicscene-miot-adapter-service-prod:"$1" .
docker push licotek/magicscene-miot-adapter-service-prod:"$1"
