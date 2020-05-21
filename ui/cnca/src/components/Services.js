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

const baseURL = (process.env.NODE_ENV === 'production') ? process.env.REACT_APP_CNCA_5GOAM_API : '/api';
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

class Services extends Component {
  isMounted = false;

  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      hasError: false,
      services: [],
    }
  }

  cancelIfUnmounted = (action) => {
    if (this.isMounted) {
      action();
    }
  }

  getServices = async() => {
    try {
      const response = await axios.get(`${baseURL}/ngcoam/v1/af/services`);

      return response.data || [];
    } catch (error) {
      console.error("Unable to get services: " + error.toString());
      throw error;
    }
  }

  deleteService = async(id) => {
    const { history } = this.props;

    try {
      await axios.delete(`${baseURL}/ngcoam/v1/af/services/${id}`);

      // Redirect back to /services to refresh the table of services
      history.push('/');
    } catch (error) {
      this.cancelIfUnmounted(() => this.setState({
        loaded: true,
        hasErrors: error,
        error: error.toString(),
      }));
    }
  }

  componentWillUnmount() {
    // Signal to cancel any pending async requests to prevent setting state
    // on an unmounted component.
    this.isMounted = false;
  }

  async componentDidMount() {
    this.isMounted = true;

    try {
      // GET services
      const services = await this.getServices() || [];

      this.cancelIfUnmounted(() => this.setState({
        loaded: true,
        services: services,
      }));
    } catch (error) {
      this.cancelIfUnmounted(() => this.setState({
        loaded: true,
        hasErrors: true,
        error: error.toString(),
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
      hasErrors,
      error,
      services,
    } = this.state;

    if (!loaded) {
      return <Loader />;
    }

    if (hasErrors) {
      throw error;
    }

    const ServiceTableRow = ({ match, item }) => {
      return (
        <TableRow>
          <TableCell>
            {item.afServiceId}
          </TableCell>
          <TableCell>
            {item.locationService.dnn}
          </TableCell>
          <TableCell>
            <Button
              onClick={() => history.push(`${match.url}/${item.afServiceId}`)}
              variant="outlined"
            >
              Edit
            </Button>
            <Button
              onClick={() => this.deleteService(item.afServiceId)}
              variant="outlined"
            >
              Delete
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
            allignItems="flex-end"
          >
            <Grid item>
              <Typography
                component="h1"
                variant="h5"
                gutterBottom
              >
                Services
              </Typography>
            </Grid>
            <Grid item>
              <Button
                onClick={() => history.push('/pfd')}
                variant="outlined"
                color="primary"
              >
                View PFD Transactions
              </Button>
            </Grid>
 
            <Grid item>
              <Button
                onClick={() => history.push('/subscriptions')}
                variant="outlined"
                color="primary"
              >
                View Subscriptions
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
                <TableCell>AF Service ID</TableCell>
                <TableCell>Data Network Name</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {
                services.length === 0
                  ? <div>No Services Registered.</div>
                  : services.map(item => <ServiceTableRow
                    key={item.afServiceId}
                    item={item}
                    match={match} />)
              }
            </TableBody>
          </Table>
        </Paper>
      </div>
    );
  }
};

export default withStyles(styles)(Services);
