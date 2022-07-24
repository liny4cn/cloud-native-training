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

## 2. 第二部分

除了将 httpServer 应用优雅的运行在 Kubernetes 之上，我们还应该考虑如何将服务发布给对内和对外的调用方。
来尝试用 Service, Ingress 将你的服务发布给集群外部的调用方吧。
在第一部分的基础上提供更加完备的部署 spec，包括（不限于）：

* Service
* Ingress

可以考虑的细节：

*如何确保整个应用的高可用。
*如何通过证书保证 httpServer 的通讯安全。

---

__作业内容 (1)__

* Kubernetes API 文件

为了方便，把各种API都放到一个 `app.yml` 文件中。后续采用 `helm` 的时候，进行分离会更方便，更工程化。

* 通过 `kubectl` 命令进行部署

```
kubectl apply -f app.yml
```
使用 `apply` ，结合 RollingUpdate 策略，可以实现滚动更新部署。

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

---

__作业内容 (2)__

* 源代码

从 `week 3` 作业修改并合并源代码到 `http-server` 目录，主要变更：

```
1. 拆分文件，进行代码按功能重构
2. 增加 zerolog，实现 日志部分的代码与配置分离
3. 增加 envflag，实现 配置变量允许从环境变量和命令行参数获取，并提供默认值。
4. 增加优雅退出代码
```

* Kubernetes API 文件

继续增加 `app.yml`，增加 `service`, `ingress`, `secret` 配置。

注： 其中 `secret` 配置为 https 私发证书。正常证书应购买，或是从 `Let's Encrypt` 获取，可以增加相关的插件进行自动获取。

* 更新部署

```
kubectl apply -f app.yml
```

执行结果：

```
namespace/cncamp-8 unchanged
secret/go-http-server-secret unchanged
configmap/go-http-server-conf unchanged
deployment.apps/go-http-server configured
horizontalpodautoscaler.autoscaling/go-http-server-hpa configured
service/lb-azure-service created
service/http-services created
ingress/gateway created
```

查看并等待部署完成(可以在命令前面加 `watch` )：
```
kubectl get all -n cncamp-8
```