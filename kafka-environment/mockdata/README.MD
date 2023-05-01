# MOCKDATA

This is a sample mockdata generator for the Customer Model.  You can extend this to build out more complex data mocking scenarios.

## Customer Model

A simple data model representing a [customer](./src/main/avro/customer.avsc).  This model contains a property `largePayload` which we use to simulate a large payload record for Kafka.

## Running Mock Data Generator

Use this project to exercise the following Kafka Producer Template Apps:
- [Java](../../../../kafka-java-template)
- [Golang](../../../../kafka-go-template)
- [NodeJS](../../../../kafka-node-template)
- [DotNet](../../../../kafka-dotnet-template)


Start up one of the template applications so it can start listening for traffic on http://ocalhost:8080/customers (all template apps default to port 8080).

Run the mockdata generator with
```shell
./gradlew bootRun
# you can change the payload size
./gradlew bootRun --args="-payloadSize 5MB"
```

The background generator will send faked customer records out every 10 seconds (POSTing to http://ocalhost:8080/customers).  The logs should show something like
```text
Customer sent: Chang Phyllis Fadel
Customer sent: Edgardo Margart O'Kon
```
if the records were posted successfully.


If you don't have Java/Gradle and have no reason to install them, then simply run this within Docker
```shell
docker run --rm -u gradle -v $(pwd):/home/gradle gradle:jdk17 gradle bootRun --args="-url http://host.docker.internal:8080/customers"
```

`host.docker.internal` is the hostname assigned to the docker host machine (your localhost).


Below shows all the options:
```text
usage: mockdata
 -payloadSize <arg>   total size of payload to build, can use values like
                      10, 250KB, 5MB (defaults to fullname)
 -url                 change target endpoint from default
                      http://localhost:8080/customers when running in
                      background mode

```
