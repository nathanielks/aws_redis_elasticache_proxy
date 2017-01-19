# Redis elasticache proxy
This is a prototype of a multi-user proxy to be placed in front of many AWS Redis Elasticache instances.

Elasticache does not support authentication or TLS -  the AWS recommendation is that the security group containing the Elasticache should only be attached to instances that are permitted to access it, however this does not work when you have shared instances with mixed tenants on them.

YOU SHOULD NOT USE THIS PROXY AS IS. It has not been validated in production, and is a proof of concept more than anything else.

## How it works
* The proxy runs in the security group attached to all the Elasticache instances that have been started by tenants and is the only way to access them.
* Tenants are given a token via some means (such as a PaaS broker that also creates the Elasticache instance) along with the hostname and port of the proxy.
* The token is actually a base64 encoded value that looks like `ELASTICACHE_HOST:PORT HASH`, where the host/port are the ones belonging to their ES instance, and HASH is `sha256(ELASTICACHE_HOST:PORT SHARED_SECRET)`. The `SHARED_SECRET` is only known to the proxy and whatever issues the token.
* The tenant configures their Redis client to connect to the host:port they have been given, and uses the token as the redis server password. If the proxy is running in TLS mode they must also configure their client for TLS.
* The server parses and validates the token, connects to the Elasticache instance the token is valid for, then uses a Linux socket splice to attach the tenant and Elasticache sockets to one another.
* At this point all packet proxying is being handled by the Linux kernel and the proxy code sits back with an iced tea.

## How does one proxy
1. This should be running on an instance with access to the elasticache instances, but nothing else should be in that security group.
2. It can operate both in TLS and non-TLS modes, but assuming non TLS...
3. Start the proxy `./proxy 0.0.0.0:6379 SHARED_SECRET` - `SHARED_SECRET` is a shared secret that is used to issue tokens
4. Generate a token for the proxy using `./generate_token.sh ELASTICACHE_HOST:ELASTICACHE_PORT SHARED_SECRET`
5. Configure your redis client to point at the proxy host and use the token generated in (4) as the redis server password

Use of TLS mode is left as an excercise for the reader. As is writing some kind of broker to issue the tokens and create Elasticache instances to tenants.
