# week3

- 构建本地镜像

- 编写 Dockerfile 将练习 2.2 编写的 httpserver 容器化

- 将镜像推送至 docker 官方镜像仓库

```
$ echo ***** | DOCKER_USER=ly4cn ./build.sh
```

- 通过 docker 命令本地启动 httpserver

```
$ docker pull ly4cn/go-native-cloud:week3

$ docker run -d --name test -p 8080:80 ly4cn/go-native-cloud:week3

$ curl localhost:8080/healthz
```

- 通过 nsenter 进入容器查看 IP 配置

```
$ PID=$(docker inspect test --format "{{.State.Pid}}")

$ sudo nsenter -t $PID -n ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
4: eth0@if5: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default
    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.17.0.2/16 brd 172.17.255.255 scope global eth0
       valid_lft forever preferred_lft forever

```
