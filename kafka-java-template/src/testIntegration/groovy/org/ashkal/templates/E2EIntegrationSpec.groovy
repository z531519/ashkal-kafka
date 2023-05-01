package org.ashkal.templates

import org.awaitility.Durations
import org.ashkal.templates.consumer.ConsumerProcessor
import org.ashkal.templates.consumer.SinkService
import org.ashkal.templates.producer.CustomerProducer
import org.mockito.Mockito
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.beans.factory.annotation.Value
import org.springframework.boot.test.mock.mockito.MockBean
import org.springframework.core.io.ClassPathResource

import static org.awaitility.Awaitility.await

class E2EIntegrationSpec extends BaseIntegrationSpec {

    @Autowired
    CustomerProducer customerProducer

    @Value("testdata.json")
    ClassPathResource testdata

    @MockBean
    SinkService sinkService


    @Autowired
    ConsumerProcessor consumerProcessor


    @Value('${topic.name}')
    String topic


    def 'e2e produce and consume integration test'() {
        when:

        customerProducer.produce(testdata.getInputStream())

        then:
        def expectedSinks = [
                'Customer: Willie Mable Beer, birthdt: 2001-03-09',
                'MINOR> Customer: Alec Brock Schroeder, birthdt: 2013-10-13',
                'MINOR> Customer: Lon Russel Rath, birthdt: 2016-08-17',
                'MINOR> Customer: Garth Altagracia Hauck, birthdt: 2008-05-20',
                'Customer: Elmo Rodrigo Hammes, birthdt: 1966-05-27',
                'MINOR> Customer: Carrol Layla Harris, birthdt: 2012-08-18',
                'MINOR> Customer: Conrad Kristel Tillman, birthdt: 2004-06-05',
                'MINOR> Customer: Reggie Andrew Kuphal, birthdt: 2012-11-23',
                'Customer: Iva Giovanni Crona, birthdt: 1997-01-11',
                'Customer: Rudy Michael Casper, birthdt: 1990-11-04']

        await().atMost(Durations.TEN_SECONDS).untilAsserted(() -> {
            expectedSinks.each {Mockito.verify(sinkService).sink(it)}
            Mockito.verify(sinkService, Mockito.times(10)).sink(Mockito.any())
        })
    }



}
