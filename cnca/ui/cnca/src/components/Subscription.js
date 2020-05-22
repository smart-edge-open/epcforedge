/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019-2020 Intel Corporation
 */

import React, { Component } from 'react';
import axios from 'axios';
import * as Schema from './schema/SubscriptionForm';
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

const baseURL = (process.env.NODE_ENV === 'production') ? process.env.REACT_APP_CNCA_AF_API : '/api';

const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms))
class Subscription extends Component {
  isMounted = false;

  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      model: {},
      trafficFiltersModel: {},
      ethFiltersModel: {},
      trafficRoutesModel: {},
      tempValidateAndSupportModel: {},
      validationResult: {},
      showErrors: false,
      id: 0,
    };
  }

  getSubscription = async (id) => {
    try {
      const response = await axios.get(`${baseURL}/af/v1/subscriptions/${id}`);

      return response.data;
    } catch (error) {
      console.error("Unable to get subscription: " + error.toString());
      throw error;
    }
  }

  postSubscription = async() => {
    const { history } = this.props;
    const { validationResult, model } = this.state;

    // Check model contains no errors before continuing
    if ("valid" in validationResult && !validationResult.valid) {
      alert('Invalid request. Please correct any errors in the form and try again');
      return;
    }

    try {
      await axios.post(`${baseURL}/af/v1/subscriptions`, model);

      this.cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "success",
        snackbarMessage: "Successfully created subscription"
      }));

      await sleep(1000);

      // Redirect user back to /subscriptions on success.
      history.push('/subscriptions');
    } catch (error) {
      this.cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "error",
        snackbarMessage: error.toString(),
      }));
    }
  }

  putSubscription = async () => {
    const { history } = this.props;
    const { validationResult, model, id } = this.state;

    // Check model contains no errors before continuing
    if ("valid" in validationResult && !validationResult.valid) {
      alert('Invalid request. Please correct any errors in the form and try again');
      return;
    };

    try {
      await axios.put(`${baseURL}/af/v1/subscriptions/${id}`, model);

      this.cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "success",
        snackbarMessage: "Successfully updated subscription"
      }));

      await sleep(1000);

      // Redirect user back to /subscriptions on success
      history.push('/subscriptions');
    } catch (error) {
      this.cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "error",
        snackbarMessage: error.toString(),
      }));
    }
  }

  deleteSubscription = async() => {
    const { history } = this.props;
    const { id } = this.state;

    try {
      await axios.delete(`${baseURL}/af/v1/subscriptions/${id}`);

      this.setState({
        snackbarOpen: true,
        snackbarVariant: "success",
        snackbarMessage: "Successfully deleted subscription"
      });

      await sleep(1000);

      // Redirect user back to /subscriptions on success
      history.push('/subscriptions');
    } catch (error) {
      this.cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "error",
        snackbarMessage: error.toString(),
      }));
    }
  }

  onValidate = () => {
    const { model } = this.state;
    const result = utils.validateBySchema(Schema.subscriptionFormSchema, model);

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

  onTrafficFiltersModelChange = (key, value) => {
    const { trafficFiltersModel } = this.state;
    const newModel = trafficFiltersModel;

    utils.selectOrSet(key, newModel, value);
    this.setState({ trafficFiltersModel: newModel });
    this.mergeSubModels();
  }

  onEthFiltersModelChange = (key, value) => {
    const { ethFiltersModel } = this.state;
    const newModel = ethFiltersModel;

    utils.selectOrSet(key, newModel, value);
    this.setState({ ethFiltersModel: newModel });
    this.mergeSubModels();
  }

  onTrafficRoutesModelChange = (key, value) => {
    const { trafficRoutesModel } = this.state;
    const newModel = trafficRoutesModel;

    utils.selectOrSet(key, newModel, value);
    this.setState({ trafficRoutesModel: newModel });
    this.mergeSubModels();
  }

  onTempValidateAndSupportModelChange = (key, value) => {
    const { tempValidateAndSupportModel } = this.state;
    const newModel = tempValidateAndSupportModel;

    utils.selectOrSet(key, newModel, value);
    this.setState({ tempValidateAndSupportModel: newModel });
    this.mergeSubModels();
  }

  mergeSubModels = () => {
    const {
      model,
      trafficFiltersModel,
      ethFiltersModel,
      trafficRoutesModel,
      tempValidateAndSupportModel,
    } = this.state;

    this.setState({
      model: {
        ...model,
        ...trafficFiltersModel,
        ...ethFiltersModel,
        ...trafficRoutesModel,
        ...tempValidateAndSupportModel,
      },
    });
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
    this.isMounted=false;
  }

  async componentDidMount() {
    this.isMounted = true;
    const { createMode, match } = this.props;

    // When the user is creating a new userplane, skip the fetch request
    if (createMode) {
      this.setState({ loaded: true });
      return;
    }

    try {
      const subscription = await this.getSubscription(match.params.id);

      this.cancelIfUnmounted(() => this.setState({
        loaded: true,
        id: match.params.id,
        model: subscription,
        trafficFiltersModel: { trafficFilters: subscription.trafficFilters },
        ethFiltersModel: { ethTrafficFilters: subscription.ethTrafficFilters },
        trafficRoutesModel: { trafficRoutes: subscription.trafficRoutes },
        tempValidateAndSupportModel: { tempValidities: subscription.tempValidities, validGeoZoneIds: subscription.validGeoZoneIds, suppFeat: subscription.suppFeat },
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
      trafficFiltersModel,
      ethFiltersModel,
      trafficRoutesModel,
      tempValidateAndSupportModel,
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
            <Button onClick={() => history.push('/subscriptions')}>
              <KeyboardBackspaceOutlined /> Back to Subscriptions
            </Button>
          </Grid>
          <Grid item>
            <Button
              onClick={(createMode) ? this.postSubscription : this.putSubscription}
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
                  ? 'Create Subscription'
                  : `Subscription ID ${id}`
              }
            </Typography>
          </Grid>

          <Grid item>
            <SchemaForm
              schema={Schema.subscriptionFormSchema}
              form={Schema.subscriptionForm}
              model={model}
              onModelChange={this.onModelChange}
              showErrors={showErrors}
            />
          </Grid>

          <Grid item xs={4}>
            <SchemaForm
              schema={Schema.trafficFiltersFormSchema}
              form={Schema.trafficFiltersForm}
              model={trafficFiltersModel}
              onModelChange={this.onTrafficFiltersModelChange}
              showErrors={showErrors}
            />
          </Grid>

          <Grid item xs={6}>
            <SchemaForm
              schema={Schema.ethFiltersFormSchema}
              form={Schema.ethFiltersForm}
              model={ethFiltersModel}
              onModelChange={this.onEthFiltersModelChange}
              showErrors={showErrors}
            />
          </Grid>

          <Grid item xs={8}>
            <SchemaForm
              schema={Schema.trafficRoutesFormSchema}
              form={Schema.trafficRoutesForm}
              model={trafficRoutesModel}
              onModelChange={this.onTrafficRoutesModelChange}
              showErrors={showErrors}
            />
          </Grid>

          <Grid item xs={10}>
            <SchemaForm
              schema={Schema.tempValidateAndSupportFormSchema}
              form={Schema.tempValidateAndSupportForm}
              model={tempValidateAndSupportModel}
              onModelChange={this.onTempValidateAndSupportModelChange}
              showErrors={showErrors}
            />
          </Grid>


        </Grid>
        {
          !createMode &&
          <Grid
            container
            justify="center"
            alignItems="center"
          >
            <Grid item style={{ margin: "10%" }}>
              <Button
                onClick={this.deleteSubscription}
                variant="outlined"
                color="secondary"
              >
                Delete Subscription
              </Button>
            </Grid>
          </Grid>
        }
      </React.Fragment>
    );
  }
}

export default Subscription;
