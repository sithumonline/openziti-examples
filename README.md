# OpenZiti Examples

## Installation

Follow the [OpenZiti installation guide](https://openziti.io/docs/learn/quickstarts/network/local-no-docker) to set up your OpenZiti environment.


```bash
source /dev/stdin <<< "$(wget -qO- https://get.openziti.io/ziti-cli-functions.sh)"; expressInstall
```

```bash
startRouter
startController
```

Add following to your `~/.bashrc` or `~/.zshrc` to set the build directory for the SDKs:

```bash
export ZITI_BIN_DIR="/Users/username/.ziti/quickstart/machine.local/ziti-bin/ziti-v1.5.4"
export PATH="$ZITI_BIN_DIR:$PATH"
```

## gRPC Example

Steps:
1. Log into OpenZiti. The host:port and username/password will vary depending on your network.

       ziti edge login localhost:1280 -u admin -p admin

1. Run this script to create everything you need.

       echo Changing to build directory
       cd $ZITI_SDK_BUILD_DIR

       echo Create the service
       ziti edge create service grpc --role-attributes grpc-service

       echo Create three identities and enroll them
       ziti edge create identity device grpc.client -a grpc.clients -o grpc.client.jwt
       ziti edge create identity device grpc.server -a grpc.servers -o grpc.server.jwt
       ziti edge enroll --jwt grpc.server.jwt
       ziti edge enroll --jwt grpc.client.jwt

       echo Create service policies
       ziti edge create service-policy grpc.dial Dial --identity-roles '#grpc.clients' --service-roles '#grpc-service'
       ziti edge create service-policy grpc.bind Bind --identity-roles '#grpc.servers' --service-roles '#grpc-service'

       echo Run policy advisor to check
       ziti edge policy-advisor services

1. Run the server.

       ./grpc-server --identity grpc.server.json --service grpc 
1. Run the client

       ./grpc-client --identity grpc.client.json --service grpc --name World

### Example output

The following is the output you'll see from the server and client side after running the previous commands.
**Server**
```
$ ./grpc-server --identity grpc.server.json --service grpc
2022/10/21 11:17:34 server listening at grpc
2022/10/21 11:18:09 Received: World
```
**Client**
```
$ ./grpc-client --identity grpc.client.json --service grpc --name World
2022/10/21 13:26:19 Greeting: Hello World
```

## Teardown

Done with the example? This script will remove everything created during setup.
```
ziti edge login localhost:1280 -u admin -p admin

echo Removing service policies
ziti edge delete service-policy grpc.dial
ziti edge delete service-policy grpc.bind

echo Removing identities
ziti edge delete identity grpc.client
ziti edge delete identity grpc.server

echo Removing service
ziti edge delete service grpc
```

