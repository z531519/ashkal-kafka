# Kafka Environment

Kafka Environment is used to support the following Kafka Template Apps
- [Java](../kafka-java-template)
- [Golang](../kafka-go-template)
- [NodeJS](../kafka-node-template)
- [DotNet](../kafka-dotnet-template)

The Kafka Template Apps are basically similar implementations using different languages.  
- A simple API to allow POSTing of Customer record and sending it out to a Kafka Topic
- A Kafka Consumer that reads off the Customer records off a Kafka Topic and prints out the customer detail
  - add the prefix "MINOR", if the customer is less than 21yrs old
- The apps shoulds be able to handle large payload sizes (1.5MB), and be able to use compression
- A simple command line utility app loads an input file containing Customer record(s) into a Kafka Topic

Each of the Template App projects will also have
- Dockerfile
- Helm Chart
which will be used for deployment into a Kubernetes environment.


**This repo hosts the following**
- [compose](./compose/) - docker compose for standing up Kafka services
- [mockdata](./mockdata/) - Customer model mock data generation

In order to run the template apps, they need the Kafka services to talk to. Use the provided [docker compose](./compose/) for setting up your local Kafka.

You will need some sample data to use with the template apps for producing and consuming records.  Use the [mockdata](./mockdata/) to generate sample files.

At this point, you are ready to play with the template apps!

## Kafka Tools

There are variety of tools available for inspecting/managing Kafka.  
- [KPow](https://kpow.io/) installed in our docker-compose, launch via (http://localhost:3000)
- [kcat](https://docs.confluent.io/platform/current/app-development/kafkacat-usage.html) aka kafkacat
  
  command-line utility that you can use to test and debug Apache Kafka® deployments. You can use kcat to produce, consume, and list topic and partition information for Kafka. Described as “netcat for Kafka”, it is a swiss-army knife of tools for inspecting and creating data in Kafka.

- [Confluent Tools](https://docs.confluent.io/platform/current/installation/cli-reference.html) all the command line utilities to manage your Kafka

  Install intructions [here](https://docs.confluent.io/platform/current/installation/installing_cp/zip-tar.html#install-cp-using-zip-and-tar-archives)