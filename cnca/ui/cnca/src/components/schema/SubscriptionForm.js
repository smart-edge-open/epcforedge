/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019-2020 Intel Corporation
 */

const subscriptionForm = [
  "afServiceId",
  "afAppId",
  "afTransId",
  "appReloInd",
  "dnn",
  "snssai",
  "externalGroupId",
  "anyUeInd",
  "subscribedEvents",
  "gpsi",
  "ipv4Addr",
  "ipv6Addr",
  "macAddr",
  "dnaiChgType",
  "notificationDestination",
  "requestTestNotification",
  "websockNotifConfig",
  "self",
];

const trafficFiltersForm = [
  {
    key: "trafficFilters",
    items: [
      {
        key: "trafficFilters[].flowId",
        title: "Flow ID",
      },
      {
        key: "trafficFilters[].flowDescriptions",
        title: "Flow Description",
      },
    ],
  },
];

const ethFiltersForm = [
  {
    key: "ethTrafficFilters",
    items: [
      {
        key: "ethTrafficFilters[].destMacAddr",
        title: "Destination MAC Address",
      },
      {
        key: "ethTrafficFilters[].ethType",
        title: "Ethernet Type",
      },
      {
        key: "ethTrafficFilters[].fDesc",
        title: "Flow Description",
      },
      {
        key: "ethTrafficFilters[].fDir",
        title: "Flow Direction",
      },
      {
        key: "ethTrafficFilters[].sourceMacAddr",
        title: "Source Mac Address",
      },
      {
        key: "ethTrafficFilters[].vlanTags",
        title: "VLAN Tags",
      },
    ],
  },
];

const trafficRoutesForm = [
  {
    key: "trafficRoutes",
    items: [
      {
        key: "trafficRoutes[].dnai",
        title: "DNAI",
      },
      {
        key: "trafficRoutes[].routeInfo",
        title: "Route Info",
      },
      {
        key: "trafficRoutes[].routeProfId",
        title: "Route Prof Id",
      },
    ],
  },
];

const tempValidateAndSupportForm = [
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
  "suppFeat",
];

const trafficFiltersFormSchema = {
  type: "object",
  title: "Traffic Filters",
  properties: {
    trafficFilters: {
      title: "Traffic Filters",
      type: "array",
      items: {
        title: "Traffic Filter",
        type: "object",
        properties: {
          flowId: {
            title: "Flow ID",
            type: "number",
          },
          flowDescriptions: {
            title: "Flow Descriptions",
            type: "array",
            items: {
              title: "Flow Description",
              type: "string",
            },
          },
        },
      },
    },
  }
};

const ethFiltersFormSchema = {
  type: "object",
  title: "Ethernet Filters",
  properties: {
    ethTrafficFilters: {
      title: "Ethernet Filters",
      type: "array",
      items: {
        title: "Ethernet Filter",
        type: "object",
        properties: {
          destMacAddr: {
            title: "Destination MAC Address",
            type: "string",
          },
          ethType: {
            title: "Ethernet Type",
            type: "string",
          },
          fDesc: {
            title: "Flow Description",
            type: "string",
          },
          fDir: {
            title: "Flow Direction",
            type: "string",
          },
          sourceMacAddr: {
            title: "Source MAC Address",
            type: "string",
          },
          vlanTags: {
            title: "VLAN Tags",
            type: "array",
            items: {
              title: "VLAN Tag",
              type: "string",
            },
          },
        },
      },
    },
  }
};

const trafficRoutesFormSchema = {
  type: "object",
  title: "Traffic Routes",
  properties: {
    trafficRoutes: {
      title: "Traffic Routes",
      type: "array",
      items: {
        title: "Traffic Route",
        type: "object",
        properties: {
          dnai: {
            title: "DNAI",
            type: "string",
          },
          routeInfo: {
            title: "Route Information",
            type: "object",
            properties: {
              ipv4Addr: {
                title: "IPv4 Address",
                type: "string",
              },
              ipv6Addr: {
                title: "IPv6 Address",
                type: "string",
              },
              portNumber: {
                title: "Port Number",
                type: "number",
              },
            },
          },
          routeProfId: {
            title: "Route Prof ID",
            type: "string",
          },
        },
      },
    },
  }
};

const tempValidateAndSupportFormSchema = {
  type: "object",
  title: "Temporal Validities",
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
    suppFeat: {
      title: "Supported Features",
      type: "string",
    },
  }
};

const subscriptionFormSchema = {
  type: "object",
  title: "Subscription",
  properties: {
    afServiceId: {
      title: "AF Service ID",
      type: "string",
    },
    afAppId: {
      title: "AF Application ID",
      type: "string",
    },
    afTransId: {
      title: "AF Transaction ID",
      type: "string",
    },
    appReloInd: {
      title: "Application Relocation ID",
      type: "boolean",
    },
    dnn: {
      title: "DNN",
      type: "string",
    },
    snssai: {
      title: "SNSSAI",
      type: "object",
      properties: {
        sst: {
          title: "SST",
          type: "number",
        },
        sd: {
          title: "SD",
          type: "string",
        },
      },
    },
    externalGroupId: {
      title: "External Group ID",
      type: "string",
    },
    anyUeInd: {
      title: "Any UE Indicated",
      type: "boolean",
    },
    subscribedEvents: {
      title: "Subscribed Events",
      type: "array",
      items: {
        title: "Subscribed Event",
        type: "string",
      },
    },
    gpsi: {
      title: "GPSI",
      type: "string",
    },
    ipv4Addr: {
      title: "IPv4 Address",
      type: "string",
    },
    ipv6Addr: {
      title: "IPv6 Address",
      type: "string",
    },
    macAddr: {
      title: "MAC Address",
      type: "string",
    },
    dnaiChgType: {
      title: "DNAI Change Type",
      type: "string",
    },
    notificationDestination: {
      title: "Notification Destination",
      type: "string",
    },
    requestTestNotification: {
      title: "Request Test Notification",
      type: "boolean",
    },
    websockNotifConfig: {
      title: "Websocket Notification Configuration",
      type: "object",
      properties: {
        websocketUri: {
          title: "Websocket URI",
          type: "string",
        },
        requestWebsocketUri: {
          title: "Request Websocket URI",
          type: "boolean",
        },
      },
    },
    self: {
      title: "Subscription URI",
      type: "string",
    },
  }
};

export {
  subscriptionForm,
  trafficFiltersForm,
  ethFiltersForm,
  trafficRoutesForm,
  tempValidateAndSupportForm,
  subscriptionFormSchema,
  trafficFiltersFormSchema,
  ethFiltersFormSchema,
  trafficRoutesFormSchema,
  tempValidateAndSupportFormSchema,
};
