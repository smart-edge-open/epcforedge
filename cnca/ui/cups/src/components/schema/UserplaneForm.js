// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

const userplaneForm = [
  {
    key: "uuid",
    style: {
      marginTop: "5%",
      marginBottom: "5%",
    },
  },
  "function",
];

const configForm = [
  "config",
];

const selectorsForm = [
  {
    key: "selectors",
    items: [
      {
        key: "selectors[].network",
        title: "Network"
      },
      {
        key: "selectors[].uli",
        title: "ULI"
      },
      {
        key: "selectors[].pdn",
        title: "PDN"
      },
    ]
  }
];

const entitlementsForm = [
  {
    key: "entitlements",
    items: [
      {
        key: "entitlements[].apns",
        title: "APNs"
      },
      {
        key: "entitlements[].imsis",
        title: "IMSIs"
      },
    ],
  },
];

const configFormSchema = {
  title: "Configuration",
  type: "object",
  properties: {
    config: {
      title: "Configuration",
      type: "object",
      properties: {
        sxa: {
          title: "Sxa",
          type: "object",
          properties: {
            cp_ip_address: {
              title: "CP IP Address",
              type: "string",
            },
            up_ip_address: {
              title: "UP IP Address",
              type: "string",
            },
          },
        },
        sxb: {
          title: "Sxb",
          type: "object",
          properties: {
            cp_ip_address: {
              title: "CP IP Address",
              type: "string",
            },
            up_ip_address: {
              title: "UP IP Address",
              type: "string",
            },
          },
        },
        s1u: {
          title: "S1-U",
          type: "object",
          properties: {
            up_ip_address: {
              title: "UP IP Address",
              type: "string",
            },
          },
        },
        s5u_sgw: {
          title: "S5-U (SGW)",
          type: "object",
          properties: {
            up_ip_address: {
              title: "UP IP Address",
              type: "string",
            },
          },
        },
        s5u_pgw: {
          title: "S5-U (PGW)",
          type: "object",
          properties: {
            up_ip_address: {
              title: "UP IP Address",
              type: "string",
            },
          },
        },
        sgi: {
          title: "SGi",
          type: "object",
          properties: {
            up_ip_address: {
              title: "UP IP Address",
              type: "string",
            },
          },
        },
        breakout: {
          title: "Breakout",
          type: "array",
          items: {
            title: "Userplane",
            type: "object",
            properties: {
              up_ip_address: {
                title: "IP Address",
                type: "string",
              },
            },
          },
        }
      },
    },
  },
};

const selectorsFormSchema = {
  type: "object",
  title: "Selectors",
  properties: {
    selectors: {
      title: "Selectors",
      type: "array",
      items: {
        title: "Selector",
        type: "object",
        properties: {
          id: {
            title: "ID",
            type: "string",
            readOnly: true,
          },
          network: {
            title: "network",
            type: "object",
            properties: {
              mcc: {
                title: "MCC",
                type: "string",
              },
              mnc: {
                title: "MNC",
                type: "string",
              },
            },
          },
          uli: {
            title: "ULI",
            type: "object",
            properties: {
              tai: {
                title: "TAI",
                type: "object",
                properties: {
                  tac: {
                    title: "TAC",
                    type: "number",
                  },
                },
              },
              ecgi: {
                title: "ECGI",
                type: "object",
                properties: {
                  eci: {
                    title: "ECI",
                    type: "number",
                  },
                },
              },
            },
          },
          pdn: {
            title: "PDN",
            type: "object",
            properties: {
              apns: {
                title: "APNs",
                type: "array",
                items: {
                  type: "string"
                },
              },
            },
          },
        },
      },
    },
  }
};

const entitlementsFormSchema = {
  type: "object",
  title: "Entitlements",
  properties: {
    entitlements: {
      title: "Entitlements",
      type: "array",
      items: {
        type: "object",
        properties: {
          id: {
            title: "ID",
            type: "string",
            readOnly: true,
          },
          apns: {
            title: "APNs",
            type: "array",
            items: {
              title: "APN",
              type: "string",
            },
          },
          imsis: {
            title: "IMSIs",
            type: "array",
            items: {
              title: "IMSI",
              type: "object",
              properties: {
                begin: {
                  title: "begin",
                  type: "string",
                },
                end: {
                  title: "end",
                  type: "string",
                },
              },
            },
          },
        },
      },
    }
  },
};

const userplaneFormSchema = {
  type: "object",
  title: "Userplane",
  required: [
    "function",
  ],
  properties: {
    id: {
      title: "ID",
      type: "string",
      readOnly: true,
    },
    uuid: {
      title: "UUID",
      type: "string",
      minLength: 36,
      maxLength: 36,
    },
    function: {
      title: "Function",
      type: "string",
      enum: [
        "NONE",
        "SGWU",
        "PGWU",
        "SAEGWU",
      ],
    },
  },
};


export {
  userplaneForm,
  configForm,
  selectorsForm,
  entitlementsForm,
  userplaneFormSchema,
  configFormSchema,
  selectorsFormSchema,
  entitlementsFormSchema,
};
