import { AvroModel } from "../kafka/avro";

/**
 * Customer Avro Schema
 */
export const Schema = `
{
  "type": "record",
  "name": "Customer",
  "namespace": "org.ashkal.templates.data.avro",
  "fields": [
    {
      "name": "birthdt",
      "type": "string"
    },
    {
      "name": "fname",
      "type": "string"
    },
    {
      "name": "fullname",
      "type": "string"
    },
    {
      "name": "gender",
      "type": "string"
    },
    {
      "name": "id",
      "type": "string"
    },
    {
      "name": "lname",
      "type": "string"
    },
    {
      "name": "mname",
      "type": "string"
    },
    {
      "name": "suffix",
      "type": "string"
    },
    {
      "name": "title",
      "type": "string"
    },
    {
      "name": "largePayload",
      "type": "string"
    }
  ]
}
`;

export function buildCustomer(source: any): Customer {
  return { ...source, Schema: Schema };
}

/**
 * the Customer Avro Model
 */
export class Customer implements AvroModel {
  Schema: string = Schema;

  id: string;

  birthdt: string;

  fname: string;

  fullname: string;

  gender: string;

  lname: string;

  mname: string;

  suffix: string;

  title: string;

  largePayload: string;
}
