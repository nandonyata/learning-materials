# Go logs in Loki

This folder has two equivalent examples. Pick one at a time:

- [`manual/`](manual/NOTES.md): start Loki, Alloy, and Grafana separately with `docker run`.
- [`compose/`](compose/NOTES.md): start the same three services with one `docker compose up -d` command.

Both examples expect the Go application in `../zerolog` to write `app.log`.
