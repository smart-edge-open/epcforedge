/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

const subscriptionModifyForm = [
  "appReloInd",
];

const temporalValidityForm = [
  {
    key: "tempValidities",
    items: [
      {
        key: "tempValidities[].startTime",
        title: "Start Time",
      },
      {
        key: "tempValidities[].stopTime",
        title: "Stop Time",
      },
    ],
  },
  "validGeoZoneIds",
];

const temporalValidityFormSchema = {
  type: "object",
  title: "Temporal Validity",
  properties: {
    tempValidities: {
      title: "Temporal Validities",
      type: "array",
      items: {
        title: "Temporal Validity",
        type: "object",
        properties: {
          startTime: {
            title: "Start Time",
            type: "string",
          },
          stopTime: {
            title: "Stop Time",
            type: "string",
          },
        },
      },
    },
    validGeoZoneIds: {
      title: "Valid Geo Zone IDs",
      type: "array",
      items: {
        title: "Valid Geo Zone ID",
        type: "string",
      },
    },
  }
};

const subscriptionModifyFormSchema = {
  type: "object",
  title: "Subscription",
  properties: {
    appReloInd: {
      title: "Application Relocation ID",
      type: "boolean",
    },
  }
};

export {
  subscriptionModifyForm,
  temporalValidityForm,
  subscriptionModifyFormSchema,
  temporalValidityFormSchema,
};
