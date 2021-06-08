## ContentService service

Сервис по хранению и менеджменту контента пользователей

#### Сборка

Выполнить make

```shell
make
```

Команда соберёт все зависимости и собранный бинарный файл положит в `bin/` 

(Опционально) Выполнить `make publish`, чтобы положить сервис в контейнер


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