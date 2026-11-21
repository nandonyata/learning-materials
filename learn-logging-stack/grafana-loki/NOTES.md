# Source Repo
https://github.com/grafana/loki
```
-> Go app (zerolog) → stdout/file → collector/agent → log storage → web UI
```
```
-> Grafana Alloy (collector) + Loki (storage) + Grafana (UI)
```

# Loki
Install Source: https://grafana.com/docs/loki/latest/setup/install/docker/

## Get the config file
```
wget https://raw.githubusercontent.com/grafana/loki/v3.7.0/cmd/loki/loki-local-config.yaml -O loki-config.yaml
```

## Command to run loki:
```
docker run --name loki -d -v $(pwd):/mnt/config -p 3100:3100 grafana/loki:3.7.0 -config.file=/mnt/config/loki-config.yaml
```

## Check
Verify that Loki is up and running.

To view readiness, navigate to http://localhost:3100/ready.
To view metrics, navigate to http://localhost:3100/metrics.


---

# Alloy
Install Source: https://grafana.com/docs/loki/latest/send-data/alloy/

## Add the config to a file '/config.alloy'
```
logging {
  level  = "info"
  format = "logfmt"
}
```

## Run container
```
docker run \
  -v "$(pwd)/learn-logging-stack/grafana-loki/config.alloy:/etc/alloy/config.alloy" \
  -p 12345:12345 \
  grafana/alloy:latest \
    run --server.http.listen-addr=0.0.0.0:12345 --storage.path=/var/lib/alloy/data \
    /etc/alloy/config.alloy
```
