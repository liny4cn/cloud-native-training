docker build -t ly4cn/go-native-cloud:week3 .
docker login -u $DOCKER_USER --password-stdin
docker push ly4cn/go-native-cloud:week3
docker logout