## ContentService service

Service to store and manage user content

#### Build

Run make

```shell
make
```

Command build all dependencies and put binary file in `bin/`

(_Optional_) Run `make publish`, to dockerize service


#### Test

You can run unit-tests
```shell
make test
```

You can run linter
```shell
make check
```

You can run integration-tests

```shell
make build publish && ./bin/run-integraion-tests.sh
```