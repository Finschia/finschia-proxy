# Finschia-Proxy

Finschia-Proxy is forked from [finschia](https://github.com/Finschia/finschia).
Finschia-Proxy is the mainnet proxy app implementation filtering tx. Other functions are the same as those of finschia.
Tx filtering is the function allowing only register tx type in allow list.

**Node**: Requires [Go 1.22+](https://golang.org/dl/)

**Warnings**: Initial development is in progress, but there has not yet been a stable.

# Quick Start

## Docker
**Build Docker Image**
```
make docker-build                # build docker image
```
or
```
make docker-build WITH_CLEVELDB=yes GITHUB_TOKEN=${YOUR_GITHUB_TOKEN}  # build docker image with cleveldb
```

_Note1_

If you are using M1 mac, you need to specify build args like this:
```
make docker-build ARCH=arm64
```

**Configure**
```
sh init_single.sh docker          # prepare keys, validators, initial state, etc.
```
or
```
sh init_single.sh docker testnet  # prepare keys, validators, initial state, etc. for testnet
```

**Run**
```
docker run -i -p 26656:26656 -p 26657:26657 -v ${HOME}/.fnsap:/root/.fnsap finschia/finschianode-proxy fnsad-proxy start
```

## Local

**Build**
```
make build
make install 
```

**Configure**
```
sh init_single.sh
```
or
```
sh init_single.sh testnet  # for testnet
```

**Run**
```
fnsad-proxy start                # Run a node
```

**visit with your browser**
* Node: http://localhost:26657/

## Localnet with 4 nodes

**Run**
```
make localnet-start
```

**Stop**
```
make localnet-stop
```


# How to contribute
check out [CONTRIBUTING.md](CONTRIBUTING.md) for our guidelines & policies for how we develop Finschia. Thank you to all those who have contributed!

