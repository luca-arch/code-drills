# Code Drills

This assessment is incomplete, it was estimated as a 3-hour task but it took me 4 hours to get to what's committed now.

## Usage

Just run `docker compose up` and wait for the build, it should not take too long.

Then navigate to [localhost:8080](http://localhost:8080/) and press the *Get Balance Sheet* button. Note the form does not do anything yet.

There are some handy commands inside the [Makefile](./Makefile) (like, the linter and test runners), run `make` to see the help screen.

## Frontend application

The content of the [frontend-app](./frontend-app) folder was bootstrapped with [Vite](https://vitejs.dev/).

To run the development version (once the Docker containers are running), use the following command:

```sh
cd frontend-app

npm run dev
```

The dev version should be reachable at [localhost:5173](http://localhost:5173/), with auto-reloading enabled.

## Testing

This repo only contains a minimal set of unit tests, most of which are rather trivial.
Please note, in order to run tests, all go and npm packages must be installed, because tests are executed on the host machine rather than Docker containers.

```sh
# Run TSX tests (host machine)
make tests-assets

# Run go tests (in container)
make tests-go
```

## TODOs

- [ ] Move test runners inside docker container
- [ ] Refactor `web.server.ListBalanceSheet()` to add automatic retries when the error is either `xero.ErrTooManyRequests` or `xero.ErrXeroDown`. See [backoff retries](https://encore.dev/blog/retries).
- [ ] Update backend's Dockerfile with [dockerize](https://github.com/jwilder/dockerize) and wait for `mock-xero:3000` before starting the webserver.
- [x] Run `make lint-go` and fix all warnings and errors where possible
- [x] Use Vite instead of react-scripts
- [x] Add `make lint-assets`
- [x] Add tests for the front-end!!!
- [ ] Refactor TS types to use camelCase starting with lowercase letters (maybe?).
- [x] Rebase commit history, possibly use [gitmoji](https://gitmoji.dev/)
