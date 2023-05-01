package org.ashkal.templates.consumer

import org.ashkal.templates.data.avro.Customer
import org.spockframework.spring.SpringBean
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest
import spock.lang.Specification

@SpringBootTest(classes = ConsumerProcessor.class)
class ConsumerProcessorSpec extends Specification {

    @SpringBean
    SinkService sinkStream = Mock()

    @Autowired
    ConsumerProcessor processor

    def 'consumer handling'() {
        when:
        def customer = new Customer(
                "2020-01-01",
                "firstName",
                "fullName",
                "M",
                "1",
                "lastName",
                "middleName",
                "suffix",
                "title",
                "largePayload"
        )
        processor.handle("a", customer)
        then:
        1 * sinkStream.sink("MINOR> Customer: fullName, birthdt: 2020-01-01")
    }

}
