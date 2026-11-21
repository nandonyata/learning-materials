# Manual containers

This is the explicit version of the stack. Each service is started separately:

```text
Go app -> app.log -> Alloy -> Loki -> Grafana UI
```

Run every command from this `manual` directory.

## 1. Start Loki

```bash
docker run --name learn-loki -d \
  -v "$(pwd)/loki-config.yaml:/etc/loki/config.yaml:ro" \
  -p 3100:3100 \
  grafana/loki:3.7.0 -config.file=/etc/loki/config.yaml
```

## 2. Start Alloy

Alloy mounts the Go project's directory read-only and tails `app.log`. On macOS,
`host.docker.internal` in `config.alloy` reaches the Loki port on your host.

```bash
docker run --name learn-alloy -d \
  -v "$(pwd)/config.alloy:/etc/alloy/config.alloy:ro" \
  -v "$(pwd)/../../zerolog:/var/log/zerolog:ro" \
  -p 12345:12345 \
  grafana/alloy:latest \
  run --server.http.listen-addr=0.0.0.0:12345 --storage.path=/var/lib/alloy/data \
  /etc/alloy/config.alloy
```

## 3. Start Grafana

```bash
docker run --name learn-grafana -d -p 3000:3000 grafana/grafana:latest
```

Open http://localhost:3000 and log in with `admin` / `admin`. Add a Loki datasource
with URL `http://host.docker.internal:3100`, then use **Explore** and query:

```logql
{service="go-app"}
```

## 4. Generate logs

```bash
cd ../../zerolog
go run .
```

Stop the manual example when finished:

```bash
docker rm -f learn-grafana learn-alloy learn-loki
```
