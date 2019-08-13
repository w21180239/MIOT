docker stop licotek-aligenie
docker rmi licotek-aligenie
docker build -t licotek-aligenie:latest .
docker run --name licotek-aligenie -d -p 8003:8080 --rm licotek-aligenie
