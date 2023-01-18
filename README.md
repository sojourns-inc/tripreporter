# tripreporter

## Building and running

### Required

Copy `.env.example` to `.env`, and `config/redis.conf.example` to `config/redis.conf`, and (optionally) modify them. These files are required for building and running the tripreporter.

### Non-Docker

Choose one of the following methods to run, `http://localhost:3000` should be accessible afterwards.

If you have issues when first running development, try running a production build first.

```bash
# Build production without Docker. This is required to have caches of static files for `make dev-server`.
make

# Build and run development UI, and server.
# Run both in two different terminals.
make dev-ui
make dev-server
```

### Docker

This is only for production use.
This will run on `http://localhost:3000` by default, modify `.env` and re-run if you would like to change the port or address.

```bash
# This will do all the building and running for you.
# If you want to do this manually, look at what run.sh is doing.
./scripts/run.sh
```
