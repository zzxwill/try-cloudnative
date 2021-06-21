

## installation
https://istio.io/latest/docs/setup/install/helm/

```shell
➜  /Users/zhouzhengxi/Programming/golang/src/github.com/zzxwill/try-cloudnative/ingress/istio/istio-1.8.2 git:(master) ✗ helm install -n istio-system istio-base manifests/charts/base
NAME: istio-base
LAST DEPLOYED: Mon Jan 25 22:49:50 2021
NAMESPACE: istio-system
STATUS: deployed
REVISION: 1
TEST SUITE: None
➜  /Users/zhouzhengxi/Programming/golang/src/github.com/zzxwill/try-cloudnative/ingress/istio/istio-1.8.2 git:(master) ✗ helm install --namespace istio-system istiod manifests/charts/istio-control/istio-discovery \
    --set global.hub="docker.io/istio" --set global.tag="1.8.2"

NAME: istiod
LAST DEPLOYED: Mon Jan 25 22:50:15 2021
NAMESPACE: istio-system
STATUS: deployed
REVISION: 1
TEST SUITE: None
➜  /Users/zhouzhengxi/Programming/golang/src/github.com/zzxwill/try-cloudnative/ingress/istio/istio-1.8.2 git:(master) ✗ helm install --namespace istio-system istio-ingress manifests/charts/gateways/istio-ingress \
    --set global.hub="docker.io/istio" --set global.tag="1.8.2"

NAME: istio-ingress
LAST DEPLOYED: Mon Jan 25 22:50:37 2021
NAMESPACE: istio-system
STATUS: deployed
REVISION: 1
TEST SUITE: None
➜  /Users/zhouzhengxi/Programming/golang/src/github.com/zzxwill/try-cloudnative/ingress/istio/istio-1.8.2 git:(master) ✗ k get all -n istio-system
NAME                                        READY   STATUS              RESTARTS   AGE
pod/istio-ingressgateway-648cb98899-5t7nd   0/1     ContainerCreating   0          12s
pod/istiod-5c6b8c6c49-pv5pg                 0/1     Pending             0          33s

NAME                           TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)                                                      AGE
service/istio-ingressgateway   LoadBalancer   172.21.5.55     <pending>     15021:32535/TCP,80:31819/TCP,443:31223/TCP,15443:30076/TCP   13s
service/istiod                 ClusterIP      172.21.12.228   <none>        15010/TCP,15012/TCP,443/TCP,15014/TCP                        34s

NAME                                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/istio-ingressgateway   0/1     1            0           13s
deployment.apps/istiod                 0/1     1            0           34s

NAME                                              DESIRED   CURRENT   READY   AGE
replicaset.apps/istio-ingressgateway-648cb98899   1         1         0       13s
replicaset.apps/istiod-5c6b8c6c49                 1         1         0       34s

NAME                                                       REFERENCE                         TARGETS         MINPODS   MAXPODS   REPLICAS   AGE
horizontalpodautoscaler.autoscaling/istio-ingressgateway   Deployment/istio-ingressgateway   <unknown>/80%   1         5         0          13s
horizontalpodautoscaler.autoscaling/istiod                 Deployment/istiod                 <unknown>/80%   1         5         1          34s
```

生成
service/istio-ingressgateway
deployment.apps/istio-ingressgateway

