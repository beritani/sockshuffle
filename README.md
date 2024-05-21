# sockshuffle

sockshuffle is a lightweight SOCKS5 proxy load balancer designed to efficiently distribute network traffic among multiple proxies.

## Run
```
% docker run -e PROXIES=socks5://localhost:1081,socks5://localhost:1082 --network host beritani/sockshuffle:latest
```

## Build and Run

```
% docker build -t sockshuffle .

% docker run -p 1080:1080 -e PROXIES=socks5://localhost:1081,socks5://localhost:1082 --network host sockshuffle
```

## WireGuard Proxies

Using the [wireproxy](https://github.com/pufferffish/wireproxy) you can host a SOCKS5 proxy using docker that routes traffic through a WireGuard connections.

```conf
# configs/wg0.conf

[Interface]
...

[Peer]
...

// Wireproxy config
[Socks5]
BindAddress = 0.0.0.0:1080
```

```yaml
# docker-compose.yml

version: '3'

services:
  wg0:
    container_name: wg0
    image: ghcr.io/pufferffish/wireproxy:latest
    volumes:
        - ./configs:/configs
    command: -c /configs/wg0.conf
    networks:
        - proxy
  
  wg1:
    container_name: wg1
    image: ghcr.io/pufferffish/wireproxy:latest
    volumes:
        - ./configs:/configs
    command: -c /configs/wg1.conf
    networks:
        - proxy

  sockshuffle:
    container_name: sockshuffle
    build: beritani/sockshuffle:latest
    ports:
        - 1080:1080
    environment:
        - PROXIES=socks5://wg0:1080,socks5://wg1:1080
        # - HOST=0.0.0.0
        # - PORT=1080
        # - USERNAME=admin
        # - PASSWORD=pass
    networks:
        - proxy

networks:
  proxy:
    driver: bridge
```

```bash
% docker compose up --build
```

### Auto-Generate Docker Compose

Using the [Pkl](https://pkl-lang.org) tool you can auto generate the docker-compose file from a directory containing all your wireguard config files (by default `compose.pkl` will look at `/configs`)

```
% pkl eval -f yaml compose.pkl > docker-compose.yml

% docker compose up
```

## Attribution

This project relies on the following libraries and projects.

- [go-socks5](https://github.com/things-go/go-socks5)
- [wireproxy](https://github.com/pufferffish/wireproxy)

## License

MIT License (MIT). Copyright (c) 2024 Sean N. (https://seann.co.uk)

For more information see [LICENSE](/LICENSE.md).