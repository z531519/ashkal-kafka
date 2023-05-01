import { eachMessageHandler } from "../../src/customer/consumer";
import { sink } from "../../src/customer/sink";

jest.mock("../../src/customer/sink");

describe("consumer handler", () => {
  const testCustomer = birthdt => ({
    fname: "fname",
    id: "1",
    birthdt,
    fullname: "fullname",
    gender: "M",
    lname: "lname",
    mname: "mname",
    suffix: "suffix",
    title: "title",
    largePayload: "largePayload"
  });

  it("Should process adult consumer record()", () => {
    eachMessageHandler(null, { value: testCustomer("1990-01-01") });
    expect(sink).toHaveBeenCalledWith(
      "Customer: fullname, birthdt: 1990-01-01"
    );
  });

  it("Should process minor consumer record()", async () => {
    eachMessageHandler(null, { value: testCustomer("2010-01-01") });
    expect(sink).toHaveBeenCalledWith(
      "MINOR> Customer: fullname, birthdt: 2010-01-01"
    );
  });
});
