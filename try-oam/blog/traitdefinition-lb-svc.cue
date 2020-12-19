apiVersion: core.oam.dev/v1alpha2
kind: TraitDefinition
metadata:
  name: loadbalancer
  annotations:
    definition.oam.dev/description: "`LoadBalancer` is used to provider load balancer for applications and publish service."
spec:
  appliesToWorkloads:
    - webservice
  workloadRefPath: spec.workloadRef
  definitionRef:
    name: loadbalancer.extend.oam.dev
  extension:
    template: |
      output: {
        apiVersion: v1
        kind: Service
        spec:
          ports: parameter.ports
          selector: parameter.selector
          type: LoadBalancer
      }

      parameter: {
      	name: string
      	ports: [...{
            name: string
            port: *80 | int
            protocol: *"TCP" | string
            targetPort: int
      	}]
        selector?: [string]: string
      }
      
