// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

export default [
  {
    id: "1",
    uuid: "a8a25f02-b20c-4a29-a42c-46087c64fb66",
    function: "SGWU",
    config: {
      sxa: {
        cp_ip_address: "1.1.1.1",
        up_ip_address: "1.1.1.1",
      },
      sxb: {
        cp_ip_address: "1.1.1.1",
        up_ip_address: "1.1.1.1",
      },
      s1u: {
        up_ip_address: "1.1.1.1",
      },
      s5u_sgw: {
        up_ip_address: "1.1.1.1",
      },
      s5u_pgw: {
        up_ip_address: "1.1.1.1",
      },
      sgi: {
        up_ip_address: "1.1.1.1",
      },
      breakout: [
        {
          up_ip_address: "1.1.1.1",
        }
      ],
      dns: [
        {
          up_ip_address: "1.1.1.1",
        }
      ],
    },
    selectors: [
      {
        id: "1",
        network: {
          mcc: "",
          mnc: "",
        },
        uli: {
          tai: {
            "tac": 0
          },
          ecgi: {
            "eci": 0
          }
        },
        pdn: {
          "apns": [
            ""
          ]
        },
      },
    ],
    entitlements: [
      {
        id: "1",
        apns: [
          "1.1.1.1"
        ],
        imsis: [
          {
            begin: "",
            end: ""
          }
        ]
      }
    ],
  },
];
