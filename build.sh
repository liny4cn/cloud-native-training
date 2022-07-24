docker build -t ly4cn/go-native-cloud:week8 .
docker login -u $DOCKER_USER --password-stdin
docker push ly4cn/go-native-cloud:week8
docker logout