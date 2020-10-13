# Try KEDA Cron

The PoC is based on [KEDA Cron](https://keda.sh/docs/2.0/scalers/cron/) and [KEDA cpu/memory resource metrics
which is introduced in KEDA v2.0-beta](https://keda.sh/blog/keda-2.0-beta/).

## Prerequisites
- [ ] HPA
- [ ] [Skaffold](https://docs.google.com/document/d/1laX-XK1gyziTLQ_ZAbRVcw0rMv6e6ILDUixmaVu2h4E/edit)
- [ ] [KEDA](https://keda.sh/docs/2.0/deploy/)

## Deploy deployment and KEDA
  Generate [skaffold.yaml](skaffold.yaml). For more details, please refer to [Try HPA](../../try-hpa-with-skaffold/cpu-utilization/README.md).

  Deploy an application (deployment) and KEDA `ScaledObject` by `skaffold dev`.
  
  After the deployment becomes stable, visit it with a heavy load.
  ```
  sudo k port-forward deploy/php-apache 80
  ab -n 10000 -c 100 http://127.0.0.1/
  ```

  Keda will scale the deployment with the combination of Cron and cpu utilization metrics.

  ```
  ➜  /Users/zhouzhengxi k get deployment php-apache --watch
  php-apache   0/2     2            0           0s
  php-apache   1/2     2            1           2s
  php-apache   2/2     2            2           2s
  php-apache   2/4     2            2           77s
  php-apache   2/4     2            2           77s
  php-apache   2/4     2            2           77s
  php-apache   2/4     4            2           77s
  php-apache   3/4     4            3           78s
  php-apache   4/4     4            4           78s
  NAME         READY   UP-TO-DATE   AVAILABLE   AGE
  php-apache   10/10   10           10          61m

  ➜  /Users/zhouzhengxi k get hpa --watch
  NAME                    REFERENCE               TARGETS                 MINPODS   MAXPODS   REPLICAS   AGE
  keda-hpa-resource-poc   Deployment/php-apache   2/1 (avg), <unknown>/5% 1         10        2          24s
  keda-hpa-resource-poc   Deployment/php-apache   1/1 (avg), <unknown>/5% 1         10        4          17s
  keda-hpa-resource-poc   Deployment/php-apache   500m/1 (avg), 251%/5%   1         10        8          61m
  keda-hpa-resource-poc   Deployment/php-apache   400m/1 (avg), 251%/5%   1         10        10         61m
  keda-hpa-resource-poc   Deployment/php-apache   400m/1 (avg), 251%/5%   1         10        10         61m
  keda-hpa-resource-poc   Deployment/php-apache   400m/1 (avg), 100%/5%   1         10        10         62m
  keda-hpa-resource-poc   Deployment/php-apache   400m/1 (avg), 100%/5%   1         10        10         62m
  keda-hpa-resource-poc   Deployment/php-apache   400m/1 (avg), 17%/5%    1         10        10         63m
  keda-hpa-resource-poc   Deployment/php-apache   400m/1 (avg), 17%/5%    1         10        10         63m
  keda-hpa-resource-poc   Deployment/php-apache   400m/1 (avg), 1%/5%     1         10        10         64m
  keda-hpa-resource-poc   Deployment/php-apache   400m/1 (avg), 1%/5%     1         10        10         64m
  keda-hpa-resource-poc   Deployment/php-apache   400m/1 (avg), 1%/20%    1         10        10         64m
  keda-hpa-resource-poc   Deployment/php-apache   400m/1 (avg), 1%/20%    1         10        10         69m
  keda-hpa-resource-poc   Deployment/php-apache   400m/1 (avg), 1%/20%    1         10        10         69m
  keda-hpa-resource-poc   Deployment/php-apache   1/1 (avg), 1%/20%       1         10        4          69m
  keda-hpa-resource-poc   Deployment/php-apache   1/1 (avg), 1%/20%       1         10        4          69m

  ➜  /Users/zhouzhengxi/Programming/golang/src/github.com/zzxwill/try-cloudnative/try-keda/resource git:(master) ✗ k get scaledobject.keda.sh --watch
  NAME           SCALETARGETKIND      SCALETARGETNAME   TRIGGERS   AUTHENTICATION   READY   ACTIVE   AGE
  resource-poc   apps/v1.Deployment   php-apache        cron                        True    True     62m
  resource-poc   apps/v1.Deployment   php-apache        cron                        True    True     62m
  resource-poc   apps/v1.Deployment   php-apache        cron                        True    True     62m
  resource-poc   apps/v1.Deployment   php-apache        cron                        True    True     62m
  
  ➜  /Users/zhouzhengxi/Programming/golang/src/github.com/zzxwill/try-cloudnative/try-keda/cron git:(master) ✗ k describe hpa keda-hpa-resource-poc
  Name:                                                        keda-hpa-resource-poc
  Namespace:                                                   default
  Labels:                                                      app.kubernetes.io/managed-by=skaffold
                                                               app.kubernetes.io/name=keda-hpa-resource-poc
                                                               app.kubernetes.io/part-of=resource-poc
                                                               app.kubernetes.io/version=v2
                                                               scaledObjectName=resource-poc
                                                               skaffold.dev/run-id=89adb07a-0c8e-43c1-a62a-a8bc98b39018
  Annotations:                                                 <none>
  CreationTimestamp:                                           Tue, 13 Oct 2020 14:59:45 +0800
  Reference:                                                   Deployment/php-apache
  Metrics:                                                     ( current / target )
    "cron-Asia-Shanghai-1xxxx-59xxxx" (target average value):  1 / 1
    resource cpu on pods  (as a percentage of request):        1% (1m) / 20%
  Min replicas:                                                1
  Max replicas:                                                10
  Deployment pods:                                             4 current / 4 desired
  Conditions:
    Type            Status  Reason              Message
    ----            ------  ------              -------
    AbleToScale     True    ReadyForNewScale    recommended size matches current size
    ScalingActive   True    ValidMetricFound    the HPA was able to successfully calculate a replica count from external metric cron-Asia-Shanghai-1xxxx-59xxxx(&LabelSelector{MatchLabels:map[string]string{scaledObjectName: resource-poc,},MatchExpressions:[]LabelSelectorRequirement{},})
    ScalingLimited  False   DesiredWithinRange  the desired count is within the acceptable range
  Events:
    Type    Reason             Age    From                       Message
    ----    ------             ----   ----                       -------
    Normal  SuccessfulRescale  8m39s  horizontal-pod-autoscaler  New size: 8; reason: cpu resource utilization (percentage of request) above target
    Normal  SuccessfulRescale  8m24s  horizontal-pod-autoscaler  New size: 10; reason: cpu resource utilization (percentage of request) above target
    Normal  SuccessfulRescale  43s    horizontal-pod-autoscaler  New size: 4; reason: All metrics below target
  ```
