// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

import React from 'react';
import { Link } from 'react-router-dom';
import { AppBar, Typography } from '@material-ui/core';
import './Header.css';

export default ({ className }) => 
    <AppBar className={className}>
      <Typography variant="h6" component="h2" id="title">
        <Link to="/">
          3GPP CUPS Management
        </Link>
      </Typography>
    </AppBar>
