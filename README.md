# oilapp-operator
基于CoreOS 的Operator Framework，创建一个简单的CRD（自定义资源），实现Tomcat应用在k8s集群的快速部署及端口暴露。

----

## 开发环境

- Golang v1.13
- Operator-sdk v0.18.1
- Kubernetes v1.18.3

## 测试验证

注册新的crd（OilappService）

```bash
export GOPROXY=https://goproxy.io
operator-sdk build hmxq.top/oilapp-operator:latest
kubectl create -f deploy/crds/oilapp.hmxq.top_oilappservices_crd.yaml
kubectl api-resources | grep oilapp
```

在本机测试，每一个oilappservice，将生成2个副本的deployment和1个NodeType类型的service

> 创建资源

```bash
export OPERATOR_NAME=oilapp-operator
operator-sdk run local --watch-namespace=default
#image: tomcat:8.5-jdk11-openjdk，nodePort: 32001
kubectl apply -f deploy/crds/oilapp.hmxq.top_v1_oilappservice_cr.yaml
#image: tomcat:10-jdk11-openjdk，nodePort: 32002
kubectl apply -f deploy/crds/oilapp.hmxq.top_v1_oilappservice_cr2.yaml
```

> 查看资源

```bash
kubectl get crd
kubectl get oilappservice
kubectl get deployment
kubectl get pod
kubectl get svc
```

> 清理资源

```bash
kubectl delete -f deploy/crds/oilapp.hmxq.top_v1_oilappservice_cr.yaml
kubectl delete -f deploy/crds/coilapp.hmxq.top_oilappservices_crd.yaml
```
