#!/bin/bash

SERVICE_NAME="demo"

trap stop INT

stop() {
    docker-compose stop
}

remove() {
    docker-compose rm -f
}

clean() {
    stop
    remove
    rm -rf tmp
}

up() {
    stop
    docker-compose up --build
}

dbconsole() {
	docker-compose exec -u postgres db psql $SERVICE_NAME
}

$*