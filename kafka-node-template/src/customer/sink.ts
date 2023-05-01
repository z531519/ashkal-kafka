/**
 * This is typically a service layer where we will send our data to after the consumer does
 * the necessary transformations.  Examples of this will be persisting data or calling out
 * API services.
 */

export function sink(data: string) {
  console.log(data);
}
