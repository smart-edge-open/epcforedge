/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2020 Intel Corporation
 */

const pfdForm = [
  "*",
];

const pfdAppForm = [
  "*",
];

const pfdAppFormSchema = {
  type: "object",
  title: "Af Application",
  properties: {
    externalAppID: {
      title: "AF application ID",
      type: "string",
    },
    allowedDelay: {
      title: "Allowed delay (in seconds)",
      type: "number",
    },
    cachingTime: {
      title: "Caching time (in seconds)",
      type: "number",
    },
    pfds: {
      title: "PFD's",
      type: "array",
      minItems: 1,
      maxItems: 10,
      items: {
        title: "",
        type: "object",
        properties: { 
          pfd: {
            notitle: "true",
            type: "object",
            properties: {
             pfdID: {
              title: "PFD ID",
              type: "string",
            },
            flowDescriptions: {
              title: "Flow Descriptions Values",
              type: "array",
              items: {
                title: "Flow Description Value",
                type: "string",
              }
            },
            urls: {
              title: "urls",
              type: "array",
              items: {
                title: "url",
                type: "string",
              }
            },
            domainNames: {
              title: "Domain Names",
              type: "array",
              items: {
                title: "Domain Name",
                type: "string",
              }
            }
            }
          } // pfd end
        }
      }
    } // pfds end
  }
};


const pfdFormSchema = {
  type: "object",
  title: "Packet Flow Descriptor",
  properties: {
    pfdDatas:{
      title: "PFD Transactions",
      type: "array",
      items: {
        title: "",
        type: "object",
        properties: {
          apps: {
            title: "Af Application",
            type: "object",
            properties: {
              externalAppID: {
                title: "AF application ID",
                type: "string",
              },
              allowedDelay: {
                title: "Allowed delay (in seconds)",
                type: "number",
              },
              cachingTime: {
                title: "Caching time (in seconds)",
                type: "number",
              },
              pfds: {
                title: "PFD's",
                type: "array",
                minItems: 1,
                maxItems: 10,
                items: {
                  title: "",
                  type: "object",
                  properties: { 
                    pfd: {
                      notitle: "true",
                      type: "object",
                      properties: {
                       pfdID: {
                        title: "PFD ID",
                        type: "string",
                      },
                      flowDescriptions: {
                        title: "Flow Descriptions Values",
                        type: "array",
                        items: {
                          title: "Flow Description Value",
                          type: "string",
                        }
                      },
                      urls: {
                        title: "urls",
                        type: "array",
                        items: {
                          title: "url",
                          type: "string",
                        }
                      },
                      domainNames: {
                        title: "Domain Names",
                        type: "array",
                        items: {
                          title: "Domain Name",
                          type: "string",
                        }
                      }
                    }
                  } // pfd end
                }
              }
            } // pfds end
          }
        } // apps end
      }
    }
  } // pfdDatas end
 }
};


export {
  pfdAppForm,
  pfdForm,
  pfdAppFormSchema,
  pfdFormSchema,
};
