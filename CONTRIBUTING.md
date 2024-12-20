## Contributing to Flight Summarizer

Thanks for using Flight Summarizer and contributing to open-source.

### Tools You'll Need

- [`just`](https://github.com/casey/just): It's like Make, but as a command
  runner. It's flippin' awesome.
- [`docker`](https://get.docker.com): Build and tests are containerized for
  consistency between macOS and Linux.

### How to contribute

1. [Create an
   issue](https://github.com/carlosonunez/flight-summarizer/issues/new)
   describing your change.

2. [Fork this repo](https://github.com/carlosonunez/flight-summarizer/fork),
   then create a pull request from it.

3. Make your change. **Tests are required**.

> **Pro-tip**: I recommend aliasing `git` to `git -c core.hooksPath='.githooks'`
> so that you can run some of the same tests that CI will run against your PR.

### Testing Flight Summarizer

```sh
just test
```

You can use these environment variables to customize the test suite:

- `TEST_EXCLUDE`: A regex of tests to exclude. Remember to add `E2E` to avoid
  running slow e2e tests.
- `TEST_FILTER`: A regex of test _functions_ to include, like
  `TestOriginAirport`.

### Building Flight Summarizer

```sh
just build
```

The binary will be located in `out/`.

You can use these environment variables to customize the build:

- `GOOS`: The OS to target.
- `GOARCH`: The architecture to target.
