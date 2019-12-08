# exmonit
The service periodically updates currency pairs exchange rates and exposes data with API

## Install

Use provided `docker-compose.yml` and change `.env` file.
It starts PostgreSQL, Promethues and the application itself.

## Metrics & monitoring

Following metrics are published to Prometheus:

`db_requests_total` database requests counter with labels `type`(save operations and queries) and `status` (failed or success)

`db_requests_duration` database requests duration

`update_duration` duration of update cycle with label for each exchange

HTTP endpoint `/status` could be used for API server health checks
