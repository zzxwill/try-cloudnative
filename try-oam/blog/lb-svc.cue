output: {
  apiVersion: "v1"
  kind: "Service"
  spec:
    ports: parameter.ports
    selector: context.name
    type: "LoadBalancer"
}

parameter: {
  ports: [...{
    name: string
    port: *80 | int
    protocol: *"TCP" | string
    targetPort: int
  }]
}
