package org.ashkal.templates.producer;

import org.ashkal.templates.data.avro.Customer;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.autoconfigure.kafka.KafkaProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.core.DefaultKafkaProducerFactory;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.core.ProducerFactory;

import java.util.Map;

/**
 * The default autoconfiguration with Spring Boot Kafka Producer does not support all kafka producer settings.
 * Here we are using @{@link KafkaProducerProperties} to add the additional configurations to support
 * large payloads.
 * <p>
 * NOTE: In most cases, the default autoconfiguration is all you need for standing up your producer and consumers.
 * This is just an example on how you can inject additional configuration settings.
 */
@Configuration
public class KafkaProducerConfig {

    @Autowired
    KafkaProperties kafkaProperties;
    @Autowired
    KafkaProducerProperties kafkaProducerProperties;


    @Bean
    public ProducerFactory<String, Customer> producerFactory(KafkaProperties kafkaProperties) {
        Map<String, Object> props = kafkaProperties.buildProducerProperties();
        props.putAll(kafkaProducerProperties.extendedProperties(props));
        return new DefaultKafkaProducerFactory<>(props);
    }

    @Bean
    public KafkaTemplate<String, Customer> customerKafkaTemplate() {
        return new KafkaTemplate<>(producerFactory(kafkaProperties));
    }
}
