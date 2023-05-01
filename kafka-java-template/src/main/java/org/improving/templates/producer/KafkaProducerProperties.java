package org.ashkal.templates.producer;

import lombok.Getter;
import lombok.Setter;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.context.properties.PropertyMapper;
import org.springframework.context.annotation.Configuration;
import org.springframework.util.unit.DataSize;

import java.time.Duration;
import java.util.Map;

/**
 * Provides extension to {@link org.springframework.boot.autoconfigure.kafka.KafkaProperties}.
 * <p>
 * Additional configuration properties used by our Kafka Producer to support large payloads.
 */
@Getter
@Setter
@Configuration
@ConfigurationProperties(prefix = "spring.kafka.producer")
public class KafkaProducerProperties {

    private DataSize maxRequestSize;
    private Duration lingerMs;

    public Map<String, Object> extendedProperties(Map<String, Object> properties) {

        PropertyMapper map = PropertyMapper.get().alwaysApplyingWhenNonNull();
        map.from(this::getLingerMs).asInt(Duration::toMillis)
                .to((value) -> properties.put(ProducerConfig.LINGER_MS_CONFIG, value));
        map.from(this::getMaxRequestSize).asInt(DataSize::toBytes)
                .to((value) -> properties.put(ProducerConfig.MAX_REQUEST_SIZE_CONFIG, value));

        return properties;
    }
}
