kubectl delete -f deploy.yaml

kubectl delete -f gateway.yaml
kubectl delete -f virtual-service.yaml

kubectl delete -f issuer.yaml
kubectl delete -f cert.yaml

kubectl delete secret httpserver -n istio-system
kubectl delete secret istio-ca -n istio-system