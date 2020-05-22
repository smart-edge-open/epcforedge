/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

import React, { Component } from 'react';
import axios from 'axios';
import * as Schema from './schema/ServiceForm';
import { SchemaForm, utils } from 'react-schema-form';
import Loader from './Loader';
import SnackbarWrapper from './SnackbarWrapper';
import { KeyboardBackspaceOutlined } from '@material-ui/icons';
import {
  Typography,
  Button,
  Grid,
  Snackbar,
} from '@material-ui/core';

const baseURL = (process.env.NODE_ENV === 'production') ? process.env.REACT_APP_CNCA_5GOAM_API : '/api';

const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms))
class Service extends Component {
  isMounted = false;

  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      model: {},
      validationResult: {},
      showErrors: false,
      id: 0,
    };
  }

  getService = async (id) => {
    try {
      const response = await axios.get(`${baseURL}/ngcoam/v1/af/services/${id}`);

      return response.data;
    } catch (error) {
      console.error("Unable to get subscription: " + error.toString());
      throw error;
    }
  }

  postService = async() => {
    const { history } = this.props;
    const { validationResult, model } = this.state;

    // Check model contains no errors before continuing
    if ("valid" in validationResult && !validationResult.valid) {
      alert('Invalid request. Please correct any errors in the form and try again');
      return;
    }

    try {
      await axios.post(`${baseURL}/ngcoam/v1/af/services`, model);

      this.cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "success",
        snackbarMessage: "Successfully created subscription"
      }));

      await sleep(1000);

      // Redirect user back to /services on success.
      history.push('/services');
    } catch (error) {
      this.cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "error",
        snackbarMessage: error.toString(),
      }));
    }
  }

  patchService = async () => {
    const { history } = this.props;
    const { validationResult, model, id } = this.state;

    // Check model contains no errors before continuing
    if ("valid" in validationResult && !validationResult.valid) {
      alert('Invalid request. Please correct any errors in the form and try again');
      return;
    };

    try {
      await axios.patch(`${baseURL}/ngcoam/v1/af/services/${id}`, model);

      this.cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "success",
        snackbarMessage: "Successfully updated subscription"
      }));

      await sleep(1000);

      // Redirect user back to /subscriptions on success
      history.push('/services');
    } catch (error) {
      this._cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "error",
        snackbarMessage: error.toString(),
      }));
    }
  }

  onValidate = () => {
    const { model } = this.state;
    const result = utils.validateBySchema(Schema.serviceFormSchema, model);

    this.setState({
      validationResult: result,
      showErrors: true,
    });
  }

  onModelChange = (key, value) => {
    const { model } = this.state;
    const newModel = model;

    utils.selectOrSet(key, newModel, value);
    this.setState({ model: newModel });

    // Validate the fields after the change to the model
    this.onValidate();
  }

  cancelIfUnmounted = (action) => {
    if (this.isMounted) {
      action();
    }
  }

  handleClose = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }

    this.setState({ snackbarOpen: false });
  }

  componentWillUnmount() {
    // Signal to cancel any pending async requests to prevent setting state on
    // an unmounted component
    this.isMounted = false;
  }

  async componentDidMount() {
    this.isMounted = true;
    const { createMode, match } = this.props;

    // When the user is creating a new service registration, skip the fetch request
    if (createMode) {
      this.setState({ loaded: true });
      return;
    }

    try {
      const service = await this.getService(match.params.id);

      this.cancelIfUnmounted(() => this.setState({
        loaded: true,
        id: match.params.id,
        model: service,
      }));
    } catch (error) {
      this.cancelIfUnmounted(() => this.setState({
        loaded: true,
        hasErrors: true,
        error: error,
      }));
    }
  }

  render() {
    const {
      createMode,
      history,
    } = this.props;

    const {
      model,
      showErrors,
      hasErrors,
      error,
      loaded,
      id,
      snackbarOpen,
      snackbarVariant,
      snackbarMessage,
    } = this.state;

    if (!loaded) {
      return <Loader />
    }

    if (hasErrors) {
      throw error;
    }

    return (
      <React.Fragment>
        <Snackbar
          anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'right',
          }}
          open={snackbarOpen}
          autoHideDuration={6000}
          onClose={this.handleClose}
        >
          <SnackbarWrapper
            onClose={this.handleClose}
            variant={snackbarVariant}
            message={snackbarMessage}
          />
        </Snackbar>

        <Grid
          container
          justify="space-between"
          alignItems="flex-end"
          style={{ marginBottom: "5%" }}
        >
          <Grid item>
            <Button onClick={() => history.push('/services')}>
              <KeyboardBackspaceOutlined /> Back to Services
            </Button>
          </Grid>
          <Grid item>
            <Button
              onClick={(createMode) ? this.postService : this.patchService}
              variant="outlined"
              color="primary"
            >
              Save
            </Button>
          </Grid>
        </Grid>

        <Grid
          container
          direction="column"
          justify="center"
          spacing={40}
        >
          <Grid item>
            <Typography
              variant="h4"
              component="h1">
              {
                createMode
                  ? 'Register Service'
                  : `Service ID ${id}`
              }
            </Typography>
          </Grid>

          <Grid item>
            <SchemaForm
              schema={Schema.serviceFormSchema}
              form={Schema.serviceForm}
              model={model}
              onModelChange={this.onModelChange}
              showErrors={showErrors}
            />
          </Grid>
        </Grid>
      </React.Fragment>
    );
  }
}

export default Service;
