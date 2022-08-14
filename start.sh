# kubectl create namespace week12
# kubectl label ns week12 istio-injection=enabled
kubectl apply -f deploy.yaml

kubectl apply -f gateway.yaml
kubectl apply -f virtual-service.yaml

kubectl apply -f issuer.yaml
kubectl apply -f cert.yaml