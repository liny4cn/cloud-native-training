# 模块十二作业

把我们的 httpserver 服务以 Istio Ingress Gateway 的形式发布出来。以下是你需要考虑的几点：

* 如何实现安全保证；
* 七层路由规则；
* 考虑 open tracing 的接入。

# 作业内容

## (0) 部署Istio 及 域名解析

* 部署 Istio

使用 `istioctl install` 进行部署， 并检查:

```shell
$ kubectl get all -n istio-system
NAME                                        READY   STATUS    RESTARTS   AGE
pod/istio-egressgateway-575d8bd99b-rvtn9    1/1     Running   0          4h35m
pod/istio-ingressgateway-6668f9548d-lxpz7   1/1     Running   0          4h35m
pod/istiod-8495d444bb-gpxbb                 1/1     Running   0          4h35m

NAME                           TYPE           CLUSTER-IP     EXTERNAL-IP     PORT(S)
                                  AGE
service/istio-egressgateway    ClusterIP      10.0.245.76    <none>          80/TCP,443/TCP
                                  4h35m
service/istio-ingressgateway   LoadBalancer   10.0.223.131   20.239.12.195   15021:32461/TCP,80:31444/TCP,443:30133/TCP,31400:31657/TCP,15443:32595/TCP   5h23m
service/istiod                 ClusterIP      10.0.138.166   <none>          15010/TCP,15012/TCP,443/TCP,15014/TCP                                        5h23m

NAME                                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/istio-egressgateway    1/1     1            1           4h35m
deployment.apps/istio-ingressgateway   1/1     1            1           5h23m
deployment.apps/istiod                 1/1     1            1           5h23m

NAME                                              DESIRED   CURRENT   READY   AGE
replicaset.apps/istio-egressgateway-575d8bd99b    1         1         1       4h35m
replicaset.apps/istio-ingressgateway-6668f9548d   1         1         1       4h35m
replicaset.apps/istio-ingressgateway-778f44479    0         0         0       5h23m
replicaset.apps/istiod-6d67d84bc7                 0         0         0       5h23m
replicaset.apps/istiod-8495d444bb                 1         1         1       4h35m
```

* 解析 `cnc12.gocloudnative.work` 域名

将域名解析到上面的 `istio-ingressgateway` 的 LoadBalancer 地址： `20.239.12.195`


## (1) 部署 http-server

为了快速实践，将week10的 api 文件进行简化，只留下必须的内容（deployment + service)。

`deploy.yaml` 文件内容如下：

```yaml
# Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-http-server
  namespace: week12
spec:
  replicas: 1
  selector:
    matchLabels:
      app: go-http-server
  template:
    metadata:
      #整个命名空间允许注入或仅指定POD注入：
      #annotations:
      #  sidecar.istio.io/inject: "true"
      labels:
        app: go-http-server
    spec:
      restartPolicy: Always
      containers:
        - name: go-http-server
          image: ly4cn/go-native-cloud:week10
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 80
---
# Service
apiVersion: v1
kind: Service
metadata:
  name: http-server
  namespace: week12
spec:
  type: ClusterIP
  ports:
    - port: 80
      protocol: TCP
      name: http
  selector:
    app: go-http-server

```

执行部署命令：

```shell
$ kubectl create namespace week12
$ kubectl label ns week12 istio-injection=enabled
$ kubectl apply -f deploy.yaml

```
等待部署完成：

```shell
$ kubectl get all -n week12
NAME                                  READY   STATUS    RESTARTS   AGE
pod/go-http-server-784cfd6cc7-fnvkc   2/2     Running   0          27s

NAME                  TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)   AGE
service/http-server   ClusterIP   10.0.72.116   <none>        80/TCP    15m

NAME                             READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/go-http-server   1/1     1            1           15m

NAME                                        DESIRED   CURRENT   READY   AGE
replicaset.apps/go-http-server-784cfd6cc7   1         1         1       15m
```

## (2) 部署 Gateway 和 Virtual Service

`gateway.yaml` 文件：

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: httpserver-gw
  namespace: week12
spec:
  selector:
    istio: ingressgateway # use istio default ingress gateway
  servers:
    - port:
        number: 443
        name: https
        protocol: HTTPS
      tls:
        mode: SIMPLE
        credentialName: httpserver # must be the same as secret
      hosts:
        - cnc12.gocloudnative.work
```

`virtual-service.yaml` 文件:

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: httpserver-vs
  namespace: week12
spec:
  hosts:
    - "cnc12.gocloudnative.work"
  gateways:
    - httpserver-gw
  http:
    - match:
        - uri:
            prefix: /
      route:
        - destination:
            port:
              number: 80
            host: http-server
```

执行部署命令：

```shell
$ kubectl apply -f gateway.yaml
gateway.networking.istio.io/httpserver-gw created

$ kubectl get gw  -n week12
NAME            AGE
httpserver-gw   12m

$ kubectl apply -f virtual-service.yaml
virtualservice.networking.istio.io/httpserver-vs configured

$ kubectl get vs  -n week12
NAME            GATEWAYS            HOSTS                          AGE
httpserver-vs   ["httpserver-gw"]   ["cnc12.gocloudnative.work"]   12m
```


## (3) HTTPS 证书

通过 https 来实现安全保证。

#### Step 1: 签发域名证书

* 安装  cert-manager.io

cert-manager 提供了各种安装方式，可以从 [https://cert-manager.io/docs/installation/]() 了解。

最简单的方式就是采用 `kubectl apply` 安装，需要可配置安装则 采用 `helm` 。

```
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.9.1/cert-manager.yaml
```
默认安装到 `cert-manager` 命名空间中。

手工进行安装验证：
```
$ kubectl get pods --namespace cert-manager
NAME                                      READY   STATUS    RESTARTS   AGE
cert-manager-55649d64b4-7gqm2             1/1     Running   0          104m
cert-manager-cainjector-666db4777-47jjp   1/1     Running   0          104m
cert-manager-webhook-6466bc8f4-w4s5f      1/1     Running   0          104m
```

* 创建 cert-manager 颁发者

`issuer.yaml` 文件

```yaml
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: letsencrypt-ca
  namespace: istio-system
spec:
  acme:
    email: ly4cn@126.com
    preferredChain: ""
    privateKeySecretRef:
      name: letsencrypt-ca
    server: https://acme-v02.api.letsencrypt.org/directory
    solvers:
      - http01:
          ingress:
            class: istio
```

由于配置与我的服务器一致，因此不需要进行修改，可以直接使用：

```
$ kubectl apply -f issuer.yaml

issuer.cert-manager.io/letsencrypt-ca created
```

* 颁发域名证书

`cert.yaml` 内容

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: httpserver
  namespace: istio-system
spec:
  dnsNames:
    - cnc12.gocloudnative.work
  issuerRef:
    group: cert-manager.io
    kind: Issuer
    name: letsencrypt-ca
  secretName: httpserver
  usages:
    - digital signature
    - key encipherment
```

```shell
$ kubectl create -f cert.yaml

certificate.cert-manager.io/httpserver created
```

## (4) 验证结果

```shell
$ curl https://cnc12.gocloudnative.work/healthz
ok
```

