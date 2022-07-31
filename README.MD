# week 10 作业

## 作业要求
---

1. 为 HTTPServer 添加 0-2 秒的随机延时；
1. 为 HTTPServer 项目添加延时 Metric；
1. 将 HTTPServer 部署至测试集群，并完成 Prometheus 配置；
1. 从 Promethus 界面中查询延时指标数据；
1. （可选）创建一个 Grafana Dashboard 展现延时分配情况。

---
## 作业内容
---

* 为 HTTPServer 添加 0-2 秒的随机延时

简单地采用 time.Sleep 进行延时

```
delay := rand.Intn(2000)
time.Sleep(time.Duration(delay+10) * time.Millisecond)
```

* 为 HTTPServer 项目添加延时 Metric

(1) 注册 metric

在 main.go 的合适位置，创建并注册 latencyFn，并在后面调用 Observe 进行指标收集

```
latencyFn := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "execution_latency_seconds",
			Help:      help,
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
		}, []string{"step"},
	)

prometheus.Register(latencyFn)

```

(2) 收集 metric

在 request handler 中，计算并收集指标

```
start := time.Now();
// ...
// 进行延时（用 time.Sleep）
// 执行请求
// ...
last := time.Now();
latencyFn.WithLabelValues("total").Observe(last.Sub(start).Seconds())
```

* 将 HTTPServer 部署至测试集群，并完成 Prometheus 配置

(1) 修改并增加 annotations 配置到 pod

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-http-server
  namespace: cncamp-8
spec:
  replicas: 1
  selector:
    matchLabels:
      app: go-http-server
#...
  template:
    metadata:
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "80"
      labels:
        app: go-http-server
    spec:
      imagePullSecrets:
        - name: go-http-server-secret
      restartPolicy: Always
      containers:
        - name: go-http-server
          image: ly4cn/go-native-cloud:week10
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 80
#...
```

(2) 部署至测试集群

```
kubectl apply -f app.yml
```


* 从 Prometheus 界面中查询延时指标数据；

Query查询 prometheus

```
histogram_quantile(0.95, sum(rate(httpserver_execution_latency_sends_bucket[5m])) by (le))
```