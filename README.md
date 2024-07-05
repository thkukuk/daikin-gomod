# Daikin AC go library and control utility

This package provides a go library to query and control Daikin AirConditioner, an application to query and set current settings and an exporter for Prometheus.

## Features

* Library
  * Discover devices on the local network
  * Query current sensor values
  * Query power consumption of the current day
  * Query and set current operating parameters
* **daikin-ac-ctrl**
  * Discover devices on the local network if none specified
  * Print current sensor data, power consumption and control options
  * Power on and off
  * Set target temperatur
* **daikin-ac-exporter**
  * Discover devices on the local network if none specified
  * Export current sensor data, power consuption and control options as [Prometheus](https://prometheus.io) metrics


## API/Library

## Container

### Public Container Image

The command to run the public available image would be:

```bash
podman run -p 9071:9071 -v <path>/config.yaml:/config.yaml registry.opensuse.org/home/kukuk/containerfile/daikin-ac-exporter:latest
```

`podman` can be replaced with `docker` without any further changes.

### Build locally

To build the container image with the `daikin-ac-exporter` binary included run:

```bash
sudo podman build --rm --no-cache --build-arg VERSION=$(cat VERSION) --build-arg BUILDTIME=$(date +%Y-%m-%dT%TZ) -t daikin-ac-exporter .
```

## Configuration

daikin-ac-exporter can be configured via command line and configuration file.

### Commandline

Available options are:
```plaintext
Usage:
  daikin-ac-exporter [flags]

Flags:
  -a, --address string   Daikin aircon address
  -c, --config string    configuration file (default "config.yaml")
  -h, --help             help for daikin-ac-exporter
  -q, --quiet            don't print any informative messages
  -v, --verbose          become really verbose in printing messages
      --version          version for daikin-ac-exporter
```

### Configuration File

By default `daikin-ac-exporter` looks for the file `config.yaml` in the local directory. This can be overriden with the `--config` option.

```yaml
# Optional: address and port to listen on, default is port 9071
listen: ":9071"
# Optional: address of Daikin AC
#address: <IPv4 address>
```
