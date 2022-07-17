# week8

## 第一部分

现在你对 Kubernetes 的控制面板的工作机制是否有了深入的了解呢？

是否对如何构建一个优雅的云上应用有了深刻的认识，那么接下来用最近学过的知识把你之前编写的 http 以优雅的方式部署起来吧，你可能需要审视之前代码是否能满足优雅上云的需求。

__作业要求：__

编写 Kubernetes 部署脚本将 httpserver 部署到 Kubernetes 集群，以下是你可以思考的维度。

* 优雅启动
* 优雅终止
* 资源需求和 QoS 保证
* 探活
* 日常运维需求，日志等级
* 配置和代码分离

---

__作业内容__

* Kubernetes API 文件

为了方便，把各种API都放到一个 `app.yml` 文件中。后续采用 `helm` 的时候，进行分离会更方便，更工程化。

* 通过 `kubectl` 命令进行部署

```
kubectl apply -f app.yml
```
使用 `apply` 可以实现重新部署和滚动更新。

执行结果：
```
namespace/cncamp-8 created
secret/go-http-server-secret created
configmap/go-http-server-conf created
deployment.apps/go-http-server created
```

查看并等待部署完成(可以在命令前面加 `watch` )：
```
kubectl get all -n cncamp-8
```
执行结果：
```
NAME                                 READY   STATUS    RESTARTS   AGE
pod/go-http-server-ccf7599f7-bzvqm   1/1     Running   0          40s
pod/go-http-server-ccf7599f7-tk82k   1/1     Running   0          40s

NAME                             READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/go-http-server   2/2     2            2           40s

NAME                                       DESIRED   CURRENT   READY   AGE
replicaset.apps/go-http-server-ccf7599f7   2         2         2       40s
```
