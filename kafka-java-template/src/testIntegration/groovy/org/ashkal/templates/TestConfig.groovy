package org.ashkal.templates

import org.ashkal.templates.producer.KafkaProducerProperties
import org.springframework.boot.ApplicationRunner
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.context.annotation.ComponentScan
import org.springframework.context.annotation.Configuration
import org.springframework.context.annotation.FilterType

@Configuration
@EnableConfigurationProperties(KafkaProducerProperties.class)
@ComponentScan(basePackages = ['org.ashkal.templates.consumer', 'org.ashkal.templates.producer'],
        excludeFilters = [
                @ComponentScan.Filter(type = FilterType.ASSIGNABLE_TYPE, value = ApplicationRunner.class)
        ]
)
class TestConfig {
}
