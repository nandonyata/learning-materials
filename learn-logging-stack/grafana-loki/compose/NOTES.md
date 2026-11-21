# Docker Compose stack

This version starts the same flow with one command:

```text
Go app -> app.log -> Alloy -> Loki -> Grafana UI
```

Run from this `compose` directory:

```bash
docker compose up -d
```

Generate log entries:

```bash
cd ../../zerolog
go run .
```

Open http://localhost:3000 and log in with `admin` / `admin`. The Loki datasource
is already configured. In **Explore**, select Loki and run:

```logql
{service="go-app"}
```

Loki is exposed at http://localhost:3101 and Alloy's UI at http://localhost:12345.
Stop this example with:

```bash
docker compose down
```
