<!-- ## 1 -->
<!-- ## -->

docker network create postgres

docker run -d --name postgres-1 \
--network postgres \
-e POSTGRES_USER=postgresadmin \
-e POSTGRES_PASSWORD=admin123 \
-e POSTGRES_DB=postgresdb \
-e PGDATA="/data" \
-v ${PWD}/postgres-1/pgdata:/data \
-v ${PWD}/postgres-1/config:/config \
-v ${PWD}/postgres-1/archive:/mnt/server/archive \
-p 5001:5432 \
postgres:latest -c 'config_file=/config/postgresql.conf'

docker exec -ti postgres-1 bash

createuser -U postgresadmin -P -c 5 --replication replicationUser

<!-- add in postgres-1 (master) pg_hba.conf -->
host    replication     replicationUser             0.0.0.0/0            md5


<!-- add in postgres-1 (master) postgresql.conf -->
wal_level = replica
max_wal_senders = 3
archive_mode = on
archive_command = 'test ! -f /mnt/server/archive/%f && cp %p /mnt/server/archive/%f'

<!-- stop postgres-1 -->
docker container stop postgres-1

<!-- re run docker postgres-1 -->
docker container start postgres-1


<!-- ## 2 -->
<!-- ## -->


docker run -it --rm \
--net postgres \
-v ${PWD}/postgres-2/pgdata:/data \
--entrypoint /bin/bash postgres:latest


<!-- makesure the postgres-1 is running before executing this -->
pg_basebackup -h postgres-1 -p 5432 -U replicationUser -D /data/ -Fp -Xs -R


docker run -d --name postgres-2 \
--network postgres \
-e POSTGRES_USER=postgresadmin \
-e POSTGRES_PASSWORD=admin123 \
-e POSTGRES_DB=postgresdb \
-e PGDATA="/data" \
-v ${PWD}/postgres-2/pgdata:/data \
-v ${PWD}/postgres-2/config:/config \
-v ${PWD}/postgres-2/archive:/mnt/server/archive \
-p 5002:5432 \
postgres:latest -c 'config_file=/config/postgresql.conf'



<!-- ## 3 Testing the database -->
<!-- ## -->

docker exec -ti postgres-1 bash

psql --username=postgresadmin postgresdb

CREATE TABLE customers (firstname text, customer_id serial, date_created timestamp);

\dt



docker exec -ti postgres-2 bash

psql --username=postgresadmin postgresdb

\dt