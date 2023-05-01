package org.ashkal.templates.data.customer;

import lombok.Getter;
import net.datafaker.Faker;
import org.ashkal.templates.data.avro.Customer;
import org.ashkal.templates.data.faker.BaseFaker;
import org.springframework.stereotype.Service;

import java.time.LocalDate;
import java.util.Arrays;
import java.util.List;

import static java.time.format.DateTimeFormatter.ISO_LOCAL_DATE;

@Service
@Getter
public class CustomerFaker extends BaseFaker {
  private static final LocalDate DATE_OLD = LocalDate.of(1995, 1, 1);
  private static final LocalDate DATE_YOUNG = LocalDate.of(2020, 1, 1);

  private static String largePayload;

  public CustomerFaker(Faker faker) {
    super(faker);
  }

  public Customer generate(int payloadSize) {
    return generate(randomId(), payloadSize);
  }

  private static String getLargePayload(int payloadSize) {
    if (largePayload == null) {
      byte[] bigdata = new byte[payloadSize];
      Arrays.fill(bigdata, (byte) '#');
      largePayload = new String(bigdata);
    }
    return largePayload;
  }

  public Customer generate(String customerId, int payloadSize) {
    LocalDate birthDate = randomDate(DATE_OLD, DATE_YOUNG);

    String firstName = faker.name().firstName();
    String middleName = faker.name().firstName();
    String lastName = faker.name().lastName();
    String fullName = firstName + " " + middleName + " " + lastName;
    return new Customer(
        birthDate.format(ISO_LOCAL_DATE),
        firstName,
        fullName,
        faker.options().nextElement(List.of("U", "N", "F", "M")),
        customerId,
        lastName,
        middleName,
        faker.name().suffix(),
        faker.name().title(),

        payloadSize == -1 ? fullName:getLargePayload(payloadSize));
  }

}
