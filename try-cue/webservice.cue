#parameter: {
        // +usage=Which image would you like to use for your service
        // +short=i
        image: string

        // +usage=Commands to run in the container
        cmd?: [...string]

        // +usage=Which port do you want customer traffic sent to
        // +short=p
        port: *80 | int
        // +usage=Define arguments by using environment variables
        env?: [...{
                // +usage=Environment variable name
                name: string
                // +usage=The value of the environment variable
                value?: string
                // +usage=Specifies a source the value of this var should come from
                valueFrom?: {
                        // +usage=Selects a key of a secret in the pod's namespace
                        secretKeyRef: {
                                // +usage=The name of the secret in the pod's namespace to select from
                                name: string
                                // +usage=The key of the secret to select from. Must be a valid secret key
                                key: string
                        }
                }
        }]
        // +usage=Number of CPU units for the service, like `0.5` (0.5 CPU core), `1` (1 CPU core)
        cpu?: string
}
