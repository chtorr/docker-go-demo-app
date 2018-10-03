#!/bin/bash

export SERVICE_NAME="demo"
export GO_PATH="github.com/chtorr/docker-go-demo-app"
export POSTGRES_PASSWORD="postgres"

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

test() {
    stop
    docker-compose up --build -d
    until docker-compose -f docker-compose-test.yml run db-wait psql -q -h "demo-db" -c '\q'; do
        >&2 echo "Postgres is unavailable - sleeping"
        sleep 2
    done
    echo "Postgres is ready"
    docker-compose -f docker-compose-test.yml run tester go test -v ./...
    stop
}

dbconsole() {
	docker-compose exec -u postgres db psql ${SERVICE_NAME}
}

$*
