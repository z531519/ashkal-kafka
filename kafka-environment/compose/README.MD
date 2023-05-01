
# Compose

This sets up the Docker containers for the local Kafka environment.  The environment consists of
- zookeeper
- Kafka Brokers (3)
- Schema Registry
- KPow (https://kpow.io/)

You will need to make sure you have docker in your machine.  If not, https://docs.docker.com/get-docker/.  Once installed, you may have to bump up the resources (cpu, memory, disk) allocations to make sure the containers have enough to work with.

To work with this docker-compose setup, use the following commands under the `compose` folder.

To start the containers
```
docker compose up -d
```
This will spin up all the containers on the background.  This may take awhile if this is your first time doing this since docker will pull all the images from the docker registry.

To verify all the running containers
```
docker compose ps
```
Sample Output:
```
docker compose ps
NAME                COMMAND                  SERVICE             STATUS              PORTS
broker-1            "/etc/confluent/dock…"   broker-1            running (healthy)   9092/tcp, 0.0.0.0:19092->19092/tcp, 0.0.0.0:17071->7071/tcp
broker-2            "/etc/confluent/dock…"   broker-2            running (healthy)   9092/tcp, 0.0.0.0:29092->29092/tcp
broker-3            "/etc/confluent/dock…"   broker-3            running (healthy)   9092/tcp, 0.0.0.0:39092->39092/tcp
connect-1           "/etc/confluent/dock…"   connect-1           running (healthy)   9092/tcp, 0.0.0.0:17072->7071/tcp, 0.0.0.0:18083->8083/tcp
kpow                "/bin/sh -c 'java $J…"   kpow                running             0.0.0.0:3000->3000/tcp
schema-registry     "/etc/confluent/dock…"   schema-registry     running (healthy)   0.0.0.0:8081->8081/tcp
zookeeper           "/etc/confluent/dock…"   zookeeper           running (healthy)   2888/tcp, 0.0.0.0:2181->2181/tcp, 3888/tcp
```

To stop the containers
```
docker compose down
```
This will spin down all the containers.  Note that this will not remove the volumes used by the containers.  So the next time you compose up the containers, they will re-mount the same volumes allowing you to resume from your previous data state.



## Docker Compose with Gradle

Spin up the kafka container with docker compose (Kafka Environment Repo).  A gradle task is provided
```shell
./gradlew :kafka:composeUp
```
A successful run should result into this showing all the container/ports made available to you:
```text
+-----------------+----------------+-----------------+
| Name            | Container Port | Mapping         |
+-----------------+----------------+-----------------+
| zookeeper       | 2181           | localhost:2181  |
+-----------------+----------------+-----------------+
| broker-1        | 19092          | localhost:19092 |
+-----------------+----------------+-----------------+
| broker-2        | 29092          | localhost:29092 |
+-----------------+----------------+-----------------+
| broker-3        | 39092          | localhost:39092 |
+-----------------+----------------+-----------------+
| connect-1       | 8083           | localhost:18083 |
+-----------------+----------------+-----------------+
| kpow            | 3000           | localhost:3000  |
+-----------------+----------------+-----------------+
| schema-registry | 8081           | localhost:8081  |
+-----------------+----------------+-----------------+

```