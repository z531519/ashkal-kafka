package org.ashkal.templates


import org.ashkal.templates.producer.KafkaProducerConfig
import org.junit.ClassRule
import org.springframework.boot.autoconfigure.kafka.KafkaAutoConfiguration
import org.springframework.boot.test.context.SpringBootTest
import org.testcontainers.containers.DockerComposeContainer
import org.testcontainers.containers.wait.strategy.Wait
import org.testcontainers.spock.Testcontainers
import spock.lang.Shared
import spock.lang.Specification

import java.time.Duration
/**
 * a Base Spec class that allows for running the integration test against testcontainers
 */
@SpringBootTest(classes = [TestConfig, KafkaAutoConfiguration, KafkaProducerConfig])
@Testcontainers
class BaseIntegrationSpec extends Specification{


    @ClassRule
    @Shared
    DockerComposeContainer environment = new DockerComposeContainer(new File("src/testIntegration/resources/docker-compose.yml"))
                    .withExposedService("schema-registry", 8081,
                            Wait.forListeningPort().withStartupTimeout(Duration.ofSeconds(60)))


}
