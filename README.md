# docker-go-demo-app
Repo to test out some local dev Docker based setups.  This does not represent best practices for building Go services, but is rather a simple test bed.

Currently this uses the embedded DNS service in Docker user defined networks to resolve container name -> IP.  It is assumed services are available on consistent ports you know ahead of time.  More advanced service discovery TBD.


# Usage

1. Install Docker for Mac/Windows/Linux
2. Git clone this repo into your $GOPATH && cd into it
3. `./build/run.sh up` (first time db startup is spammy)

From another terminal window try some things out:
```sh
curl http://127.0.0.1:5000/
# Hello from demo

curl http://127.0.0.1:5000/demo
# Value from postgres: test name

./build/run.sh dbconsole
# docker-compose exec -u postgres db psql demo
# psql (10.5 (Debian 10.5-1.pgdg90+1))
# Type "help" for help.

# demo=# \d
#                 List of relations
#  Schema |       Name        |   Type   |  Owner
# --------+-------------------+----------+----------
#  public | demo              | table    | postgres
#  public | demo_id_seq       | sequence | postgres
#  public | schema_migrations | table    | postgres
# (3 rows)

# notice the user defined network name
docker network ls
# NETWORK ID          NAME                          DRIVER              SCOPE
# 422b1ed09bee        bridge                        bridge              local
# 0ba4336f94f6        docker-go-demo-app_internal   bridge              local
# ab0811323388        host                          host                local
# ...
```

Docker compose will keep a data volume around after stopping.  To fully reset everything:
```sh
./build/run.sh clean
```
