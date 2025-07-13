# Test Guide for Brick Smart Template

This document describes the automated API testing process for all smart device services in the Brick Smart Template project.

## Overview

The test script (`scripts/test.sh`) provides a unified way to test the API endpoints of all example devices (cleaner, lighting, thermostat). It covers service health checks, configuration, start/stop/restart, status polling, and data validation.

## How to Run Tests

```bash
make test
# or
./scripts/test.sh
```

## What the Test Script Does

For each device (cleaner, lighting, thermostat):

1. **Clean & Start**: Ensures the container is in a clean state, then starts it.
2. **Health Check**: Verifies the service is running and healthy via `/health` endpoint.
3. **Configure & Start App**: Sends configuration and start commands to the app via `/app/configure` and `/app/start`.
4. **Status & Data Polling**: Polls `/app/status` and `/app/data` endpoints to check running state and output.
5. **Stop & Restart**: Stops the app, then restarts it, verifying state transitions.
6. **Final Data Check**: Ensures the app resumes correctly after restart.

## API Endpoints Covered

- `GET /health` — Service health check
- `POST /app/configure` — Configure the app
- `POST /app/start` — Start the app
- `POST /app/stop` — Stop the app
- `POST /app/restart` — Restart the app
- `GET /app/status` — Get current status
- `GET /app/data` — Get current data

## Output

The script prints color-coded logs for each step, including API responses and status summaries. At the end, a summary of all device tests is shown.

## Notes

- Make sure all containers are built and running before testing.
- The script is idempotent and can be run multiple times.
- For troubleshooting, check the logs printed by the script or inspect the running containers with `docker ps` and `docker logs <container>`. 