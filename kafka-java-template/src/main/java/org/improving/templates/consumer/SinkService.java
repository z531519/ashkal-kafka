package org.ashkal.templates.consumer;

import org.springframework.stereotype.Component;

/**
 * A simple service that services as a final sink after our consumer handles the necessary transformations.
 * This will normally be an external service call or persisting data
 *
 */
@Component
public class SinkService {

    public void sink(String data) {
        System.out.println(data);
    }
}
