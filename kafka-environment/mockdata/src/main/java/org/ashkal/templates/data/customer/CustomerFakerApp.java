package org.ashkal.templates.data.customer;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.cli.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.scheduling.annotation.EnableScheduling;
import org.springframework.util.unit.DataSize;

import java.io.OutputStream;

@Slf4j
@SpringBootApplication
@EnableScheduling
@RequiredArgsConstructor
public class CustomerFakerApp implements ApplicationRunner {

  private final CustomerFaker customerFaker;

  @Autowired
  CustomerFakerBackground backgroundApp;

  @Value("${background.customer.endpoint}")
  String customerEndpoint;

  public static void main(String[] args) throws Exception {
    SpringApplication.run(CustomerFakerApp.class, args);
  }

  @Override
  public void run(ApplicationArguments args) throws Exception {

    Options options = new Options();

    options.addOption("url", true, "change target endpoint from default http://localhost:8080/customers when running in background mode");
    options.addOption("payloadSize", true, "total size of payload to build, can use values like 10, 250KB, 5MB (defaults to fullname)");

    CommandLineParser parser = new DefaultParser();
    try {
      CommandLine cmd = parser.parse(options, args.getSourceArgs());

      int payloadSize = -1;
      String backgroundUrl = customerEndpoint;
      OutputStream out = System.out;

      if (cmd.hasOption("url")) {
        backgroundUrl = cmd.getOptionValue("url");
      }

      if (cmd.hasOption("payloadSize")) {
        DataSize dataSize = DataSize.parse(cmd.getOptionValue("payloadSize"));
        payloadSize = Long.valueOf(dataSize.toBytes()).intValue();
      }

      backgroundApp.startBackground(backgroundUrl, payloadSize);

    } catch(ParseException e) {
      HelpFormatter formatter = new HelpFormatter();
      formatter.printHelp("mockdata", options);  
    }

  }

}
