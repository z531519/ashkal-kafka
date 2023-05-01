package org.ashkal.templates.consumer;

import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.ashkal.templates.data.avro.Customer;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

import java.time.LocalDate;

@Slf4j
@Component
public class ConsumerProcessor {

    @Autowired
    public SinkService sinkService;

    @KafkaListener( topics = {"${topic.name}" })
    public void consumer(ConsumerRecord<String, Customer> record) {
        handle(record.key(), record.value());
    }



    /**
     * Your processing logic here.
     * <p>
     * This template logic simply outputs details of the record and prepends `MINOR` when the customer age is below 21.
     */
    public void handle(String key, Customer customer) {
        LocalDate birthdt = LocalDate.parse(customer.getBirthdt());
        LocalDate cutoff = LocalDate.now().minusYears(21);
        String transformed = cutoff.isBefore(birthdt) ?
                String.format("MINOR> Customer: %s, birthdt: %s", customer.getFullname(), customer.getBirthdt()):
                String.format("Customer: %s, birthdt: %s", customer.getFullname(), customer.getBirthdt());
        sinkService.sink(transformed);
    }
}
