/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

const serviceForm = [
  "dnai",
  "dnn",
  "tac",
  "priDns",
  "secDns",
  "upfIp",
  "snssai",
];


const serviceFormSchema = {
  type: "object",
  title: "Registered Services",
  properties: {
    dnai: {
      title: "DNAI",
      type: "string",
    },
    dnn: {
      title: "DNN",
      type: "string",
    },
    tac: {
      title: "TAC",
      type: "number",
    },
    priDns: {
      title: "Primary DNS Address",
      type: "string",
    },
    secDns: {
      title: "Secondary DNS Address",
      type: "string",
    },
    upfIp: {
      title: "UPF IP Address",
      type: "string",
    },
    snssai: {
      title: "SNSSAI",
      type: "string",
    },
  }
};

export {
  serviceForm,
  serviceFormSchema,
};
