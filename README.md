## Go Example

Code created in [GO](https://golang.org) using [POSTGRESQL](https://www.postgresql.org/)

To run you must have the [GO] configured on your box with access to a postgresql server.

## Install
	go get github.com/lib/pq
    go get github.com/BurntSushi/toml

## Configure
    Create user/database on postgres
    sudo su - postgres -c "psql --dbname=DBNAME" < schema/graph.schema

## Load
    go run import.go <URL | LOCALFILE>
    example:
      go run import.go data/graph1.xml

## Query
    go run query.go <LOCALFILE>
    example:
      go run query.go queries/q1.json
      cat queries/q1.json | go run query.go
