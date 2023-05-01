package org.ashkal.templates.data.customer;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.ashkal.templates.data.avro.Customer;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.web.client.RestTemplateBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.scheduling.TaskScheduler;
import org.springframework.scheduling.support.PeriodicTrigger;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClientException;
import org.springframework.web.client.RestTemplate;

import java.time.Duration;

@Slf4j
@Service
@RequiredArgsConstructor
public class CustomerFakerBackground implements Runnable {

  private final CustomerFaker customerFaker;

  RestTemplate restTemplate;
  int payloadSize = -1;

  String customerEndpoint;

  @Value("${background.duration}")
  Duration duration;

  @Bean
  RestTemplate restTemplate() {
    RestTemplateBuilder builder = new RestTemplateBuilder();
    this.restTemplate = builder.build();
    return this.restTemplate;
  };
  @Autowired
  TaskScheduler taskScheduler;

  public void startBackground(String customerEndpoint, int payloadSize) {
    log.info("Running in background mode");
    this.customerEndpoint = customerEndpoint;
    this.payloadSize = payloadSize;
    taskScheduler.schedule( this,
            new PeriodicTrigger(this.duration.toMillis()));

  }

  @Override
  public void run() {
    Customer customer = customerFaker.generate(payloadSize);
    HttpHeaders headers = new HttpHeaders();
    headers.setContentType(MediaType.APPLICATION_JSON);
    HttpEntity<String> request = new HttpEntity<>(customer.toString(), headers);
    try {
      restTemplate.postForLocation(customerEndpoint, request);
      System.out.printf("Customer sent: %s\n", customer.getFullname());
    } catch(RestClientException e) {
      e.printStackTrace();
    }
  }
}
