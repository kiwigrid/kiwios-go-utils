# KiwiOS Go Utilities

This repository offers various utility packages for Go artifacts in KiwiOS. It does not provide any direct executable.

## GatewayJWT

Provides the function `GetGatewayJWT(tokenURL string)` which fetches a JWT from the [Gateway Identity Service](https://bitbucket.dev.kiwigrid.com/projects/GM/repos/gateway-registry). This JWT can be used to authenticate against Kiwigrid Cloud services.

The `tokenURL` is usually retrieved via the [gateway-registry-bridge](https://bitbucket.dev.kiwigrid.com/projects/BOSS/repos/gateway-registry-bridge).
