/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019-2020 Intel Corporation
 */

import React, { Component } from 'react';
import axios from 'axios';
import Loader from './Loader';
import { withStyles } from '@material-ui/core/styles';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Button,
  Typography,
  Grid,
  Paper,
} from '@material-ui/core';

const baseURL = (process.env.NODE_ENV === 'production') ? process.env.REACT_APP_CNCA_AF_API : '/api';
const LANDING_URL = process.env.REACT_APP_LANDING_UI_URL;

const styles = theme => ({
  paper: {
    marginTop: theme.spacing(3),
    marginBottom: theme.spacing(3),
    padding: theme.spacing(2),
    [theme.breakpoints.up(600 + theme.spacing(3) * 2)]: {
      marginTop: theme.spacing(6),
      marginBottom: theme.spacing(6),
      padding: theme.spacing(3),
    },
  },
});

class Subscriptions extends Component {
  _isMounted = false;

  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      hasError: false,
      subscriptions: [],
    }
  }

  _cancelIfUnmounted = (action) => {
    if (this._isMounted) {
      action();
    }
  }

  _getSubscriptions = async () => {
    try {
      const response = await axios.get(`${baseURL}/af/v1/subscriptions`);

      return response.data || [];
    } catch (error) {
      console.error("Unable to get subscriptions: " + error.toString());
      throw error;
    }
  }

  parseUUID(selfURL) {
    var split = selfURL.split("/"); 
    if (split.length === 1){
        return "";
    }
    var UUID = split[split.length-1];

    return UUID;
  }

  componentWillUnmount() {
    // Signal to cancel any pending async requests to prevent setting state
    // on an unmounted component.
    this._isMounted = false;
  }

  async componentDidMount() {
    this._isMounted = true;

    try {
      // Fetch subscriptions.
      const subscriptions = await this._getSubscriptions() || [];

      // Update subscriptions if the component is mounted.
      this._cancelIfUnmounted(() => this.setState({
        loaded: true,
        subscriptions: subscriptions,
      }));
    } catch (error) {
      // Update error if the component is mounted.
      this._cancelIfUnmounted(() => this.setState({
        loaded: true,
        hasError: true,
        error: error,
      }));
    }
  }

  render() {
    const {
      classes,
      match,
      history,
    } = this.props;

    const {
      loaded,
      hasError,
      error,
      subscriptions,
    } = this.state;

    if (!loaded) {
      return <Loader />;
    }

    if (hasError) {
      throw error;
    }

    const SubscriptionTableRow = ({ match, history, item }) => {
      return (
        <TableRow>
          <TableCell>
            {this.parseUUID(item.self)}
          </TableCell>
          <TableCell>
            {item.afServiceId}
          </TableCell>
          <TableCell>
            {item.afAppId}
          </TableCell>
          <TableCell>
            <Button
              onClick={() => history.push(`${match.url}/patch/${this.parseUUID(item.self)}`)}
              variant="outlined"
            >
              Patch
            </Button>
            <Button
              onClick={() => history.push(`${match.url}/edit/${this.parseUUID(item.self)}`)}
              variant="outlined"
            >
              Edit
            </Button>
          </TableCell>
        </TableRow>
      );
    }

    return (
      <div>
        <Grid
          allignItems="flex-start"
        >
          <Grid item>
            <Button
              onClick={() => window.location.assign(`${LANDING_URL}/`)}
            >
              Back to Home Page 
            </Button>
          </Grid>
        </Grid>

        <Paper className={classes.paper}>
          <Grid
            container
            direction="row"
            justify="space-between"
            alignItems="flex-end"
          >
            <Grid item>
              <Typography
                component="h1"
                variant="h5"
                gutterBottom
              >
                Subscriptions
              </Typography>
            </Grid>

            <Grid item>
              <Button
                onClick={() => history.push('/services')}
                variant="outlined"
                color="primary"
              >
                View Services
              </Button>
            </Grid>

            <Grid item>
              <Button
                onClick={() => history.push(`${match.url}/create`)}
                variant="outlined"
                color="primary"
              >
                Create
              </Button>
            </Grid>
          </Grid>

          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Subscription ID</TableCell>
                <TableCell>Service ID</TableCell>
                <TableCell>App ID</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {
                subscriptions.length === 0
                  ? <div>No subscriptions to display.</div>
                  : subscriptions.map(item => <SubscriptionTableRow
                    key={item.afServiceId}
                    item={item}
                    history={history}
                    match={match} />)
              }
            </TableBody>
          </Table>
        </Paper>
      </div>
    );
  }
};

export default withStyles(styles)(Subscriptions);

