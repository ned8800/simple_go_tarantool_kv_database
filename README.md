# simple_go_tarantool_kv_database

В этом репозитории лежит исходный код простого http сервера для связи с tarantool 

API взаимодействия представлен следующими правилами:

API:
 
* POST /kv body: {key: "test", "value": {SOME ARBITRARY JSON}} 
* PUT kv/{id} body: {"value": {SOME ARBITRARY JSON}}
* GET kv/{id} 
* DELETE kv/{id}

 - POST  возвращает 409 если ключ уже существует, 

 - POST, PUT возвращают 400 если боди некорректное

 - PUT, GET, DELETE возвращает 404 если такого ключа нет

 - все операции логируются

Работающий http сервер расположен по адресу http://217.16.20.177:8090/kv .

Работоспрособность проверялась через curl.
Примеры запросов к серверу:

 * GET:

``` bash
curl -X GET -H 'Content-Type: application/json' http://217.16.20.177:8090/kv/test
```

 * POST:

``` bash
curl -X POST -H 'Content-Type: application/json' -d '{ "key": "test", "value": {"SOME": "ARBITRARY JSON"} }' http://217.16.20.177:8090/kv
```

 * PUT:

``` bash
curl -X PUT -H 'Content-Type: application/json' -d '{ "key": "test", "value": {"SOME": "NEW ARBITRARY JSON"} }' http://217.16.20.177:8090/kv/test
```

 * DELETE:

``` bash
curl -X DELETE -H 'Content-Type: application/json' http://217.16.20.177:8090/kv/test
```
