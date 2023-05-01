package org.ashkal.templates.api;

import com.fasterxml.jackson.core.JsonProcessingException;
import org.apache.kafka.clients.producer.RecordMetadata;
import org.ashkal.templates.data.avro.Customer;
import org.ashkal.templates.producer.CustomerProducer;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.kafka.support.SendResult;
import org.springframework.util.concurrent.ListenableFuture;
import org.springframework.web.bind.annotation.*;

import java.util.concurrent.ExecutionException;

/**
 * Endpoint for producing messages.  Sample request:
 * curl --location --request POST 'http://localhost:8080/customers' \
 * --header 'Content-Type: application/json' \
 * --data-raw '{
 *     "birthdt": "1900-12-30",
 *     "fname": "Samuel",
 *     "fullname": "Samuel Francine Leffler",
 *     "gender": "N",
 *     "id": "asd",
 *     "lname": "Leffler",
 *     "mname": "Francine",
 *     "suffix": "V",
 *     "title": "Global Applications Designer",
 *     "largePayload": "Samuel Francine Leffler"
 * }'
 *
 */
@RestController
@RequestMapping(value = "/customers")
public class CustomerController {

    @Autowired
    CustomerProducer producer;

    @GetMapping
    public ResponseEntity<String> ping() {
        return new ResponseEntity("Hello", HttpStatus.OK);
    }


    public record RecordMetadataDto(String topic, int topicPartition, long baseOffset, long timestamp,
                                    int serializedKeySize, int serializedValueSize) {}


    @PostMapping (consumes = MediaType.APPLICATION_JSON_VALUE, produces = MediaType.APPLICATION_JSON_VALUE)
    public ResponseEntity<RecordMetadataDto> intake(@RequestBody Customer customer) throws ExecutionException, InterruptedException, JsonProcessingException {
        ListenableFuture<SendResult<String, Customer>> future = producer.produce(customer);
        SendResult<String, Customer> result = future.get();

        RecordMetadata m = result.getRecordMetadata();
        RecordMetadataDto dto = new RecordMetadataDto(m.topic(), m.partition(), m.offset(), m.timestamp(), m.serializedKeySize(), m.serializedValueSize());

        return new ResponseEntity<>(dto,  HttpStatus.CREATED);
    }
}
