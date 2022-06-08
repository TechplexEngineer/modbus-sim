[Modbus Simulator](github.com/techplexengineer/modbus-sim)
================

A simple modbus/tcp server based on the modbus server implementation from github.com/tbrandon/mbserver

The Holding Registers and Input Registers are defined as follows:

### Holding Registers

To read Holding Register `100` use `400100`

| Address | Description |
| :-----: | :---------: |
| 100     | 0xff00 ie. 65280 |
| 101     | 0xffff ie. 65535 or -1 |
| 102     | 0x0000 ie. 0 |

| 201     | artificially generates error: IllegalFunction |
| 202     | artificially generates error: IllegalDataAddress |
| 203     | artificially generates error: IllegalDataValue |
| 204     | artificially generates error: SlaveDeviceFailure |
| 205     | artificially generates error: AcknowledgeSlave |
| 206     | artificially generates error: SlaveDeviceBusy |
| 207     | artificially generates error: NegativeAcknowledge |
| 208     | artificially generates error: MemoryParityError |
| 210     | artificially generates error: GatewayPathUnavailable |
| 211     | artificially generates error: GatewayTargetDeviceFailedtoRespond |

| 300     | uptime msb |
| 301     | uptime lsb |
| 302     | application start time msb |
| 303     | application start time lsb |
| 400     | unixtime msb |
| 401     | unixtime lsb |
| 500     | math.pi msb |
| 501     | math.pi lsb |

### Input Registers

To read Input Register `100` use `300100`

| Address | Description |
| :-----: | :---------: |
| 100     | 0xff00 ie. 65280 |
| 101     | 0xffff ie. 65535 or -1 |
| 102     | 0x0000 ie. 0 |

| 201     | artificially generates error: IllegalFunction |
| 202     | artificially generates error: IllegalDataAddress |
| 203     | artificially generates error: IllegalDataValue |
| 204     | artificially generates error: SlaveDeviceFailure |
| 205     | artificially generates error: AcknowledgeSlave |
| 206     | artificially generates error: SlaveDeviceBusy |
| 207     | artificially generates error: NegativeAcknowledge |
| 208     | artificially generates error: MemoryParityError |
| 210     | artificially generates error: GatewayPathUnavailable |
| 211     | artificially generates error: GatewayTargetDeviceFailedtoRespond |

| 300     | uptime msb |
| 301     | uptime lsb |
| 302     | application start time msb |
| 303     | application start time lsb |
| 400     | unixtime msb |
| 401     | unixtime lsb |
| 500     | math.pi msb |
| 501     | math.pi lsb |

## Supported Architectures
Simply pulling `techplex/modbus-sim:latest` should retrieve the correct image for your arch.

The architectures supported by this image are:
| Architecture | Available |
| :----------: | :-------: |
| x86-64       | ✅        |
| arm64        | ✅        |
| armhf        | ✅        |


## Application Setup
The application can be accessed at tcp://yourhost:1502

## Usage
Here are some example snippets to help you get started creating a container.

### docker-compose

```yaml
---
version: "2.1"
services:
  modbus:
    image: techplex/modbus-sim:latest
    container_name: modbus
    ports:
      - 1502:1502
    restart: unless-stopped
```

### docker cli

```bash
docker run -d \
  --name=modbus \
  -p 1502:1502 \
  --restart unless-stopped \
  techplex/modbus-sim:latest
```

## Updating Info

Below are the instructions for updating containers:

### Via Docker Compose

* Update all images: `docker-compose pull`
  * or update a single image: `docker-compose pull techplex/modbus-sim`
* Let compose update all containers as necessary: `docker-compose up -d`
  * or update a single container: `docker-compose up -d techplex/modbus-sim`
* You can also remove the old dangling images: `docker image prune`

### Via Docker Run

* Update the image: `docker pull techplex/modbus-sim:latest`
* Stop the running container: `docker stop techplex/modbus-sim`
* Delete the container: `docker rm techplex/modbus-sim`
* Recreate a new container with the same docker run parameters as instructed above
* You can also remove the old dangling images: `docker image prune`

## Building locally

If you want to make local modifications to these images for development purposes or just to customize the logic:

```bash
git clone https://github.com/techplexengineer/modbus-sim.git
cd modbus-sim
docker build \
  --no-cache \
  --pull \
  -t techplex/modbus-sim:latest .