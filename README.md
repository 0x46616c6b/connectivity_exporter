# Connectivity Exporter

Connectivity Metrics Exporter for Prometheus, written in Go.

[![Build Status](https://travis-ci.org/0x46616c6b/connectivity_exporter.svg?branch=master)](https://travis-ci.org/0x46616c6b/connectivity_exporter) [![](https://images.microbadger.com/badges/version/0x46616c6b/connectivity_exporter.svg)](https://microbadger.com/images/0x46616c6b/connectivity_exporter "Get your own version badge on microbadger.com")

## Run

    docker run -d -p 9449:9449 0x46616c6b/connectivity_exporter:latest

## Configuration

    -http.hosts string
        Comma seperated list with hosts to check (default "google.com,facebook.com,github.com")
    -http.timeout string
        Timeout for the HTTP Checks (default "5s")
    -log.format value
        Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true" (default "logger:stderr")
    -log.level value
        Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal] (default "info")
    -web.listen-address string
        Address to listen on for web interface and telemetry. (default ":9449")
    -web.telemetry-path string
        Path under which to expose metrics. (default "/metrics")

## Metrics

    # Duration of the request
    connectivity_http_request_time_ns
    # Boolean value with 1 if the request was successful
    connectivity_http_request_successful
