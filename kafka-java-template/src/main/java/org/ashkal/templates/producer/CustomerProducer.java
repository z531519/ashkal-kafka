package org.ashkal.templates.producer;

import lombok.extern.slf4j.Slf4j;
import org.ashkal.templates.data.avro.Customer;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.SendResult;
import org.springframework.stereotype.Component;
import org.springframework.util.concurrent.ListenableFuture;
import tech.allegro.schema.json2avro.converter.JsonAvroConverter;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;

@Slf4j
@Component
public class CustomerProducer {

    @Autowired
    KafkaTemplate<String, Customer> kafkaTemplate;

    /**
     * The target kafka topic message will be sent to
     */
    @Value("${topic.name}") String topic;

    /**
     * Your producer logic here
     * <p>
     * The producer template code simply reads a multi-line avro entries representing Customer records and
     *  sends it to the target kafka topic.
     */
    public void produce(InputStream in) throws IOException {
        JsonAvroConverter converter = new JsonAvroConverter();
        BufferedReader reader = new BufferedReader(new InputStreamReader(in));
        while(reader.ready()) {
            String jsonLine = reader.readLine();

            Customer customer = converter.convertToSpecificRecord(jsonLine.getBytes(), Customer.class, Customer.SCHEMA$);
            produce(customer);
            log.info("Message sent");
        }
    }

    public ListenableFuture<SendResult<String, Customer>>  produce(Customer customer) {
        return kafkaTemplate.send(topic, customer.getId(), customer);

    }
}
