package org.ashkal.templates;


import lombok.extern.slf4j.Slf4j;
import org.ashkal.templates.producer.KafkaProducerProperties;
import org.springframework.boot.ApplicationRunner;
        import org.springframework.boot.SpringApplication;
        import org.springframework.boot.autoconfigure.SpringBootApplication;
        import org.springframework.boot.context.properties.EnableConfigurationProperties;
        import org.springframework.context.annotation.ComponentScan;
        import org.springframework.context.annotation.FilterType;

@Slf4j
@SpringBootApplication
@EnableConfigurationProperties(KafkaProducerProperties.class)
@ComponentScan( excludeFilters = {
        @ComponentScan.Filter(type= FilterType.ASSIGNABLE_TYPE, value= ApplicationRunner.class)
})
public class DemoApp {
    /**
     * The Producer Application
     */
    public static void main(final String[] args) {
        SpringApplication.run(DemoApp.class, args);
    }
}
