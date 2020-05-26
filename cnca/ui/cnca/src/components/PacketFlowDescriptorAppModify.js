/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2020 Intel Corporation
 */

import React, { Component } from 'react';
import axios from 'axios';
import * as Schema from './schema/PFDForm';
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

const baseURL = (process.env.NODE_ENV === 'production') ? 
                    process.env.REACT_APP_CNCA_AF_API : '/api';

const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms))
class PacketFlowDescriptorAppModify extends Component {
  isMounted = false;

  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      model: {},
      validationResult: {},
      showErrors: false,
      id: 0,
      tId: 0,
      appId: 0,
    };
  }

  /* Retrieve PFD Application data from AF Server. */
  getPfdApplication = async (tId, appId) => {
    try {
      const response = await axios.get(
          `${baseURL}/af/v1/pfd/transactions/${tId}/applications/${appId}`);

      return response.data;
    } catch (error) {
      console.error("Unable to get pfd transactions: " + error.toString());
      throw error;
    }
  }

  /* 
   * Convert uiSchema (uiModel) to afSchema (afModel). 
   * Variable with prefix "af" corresponds to json data representation 
   * according to af schema.
   * Variable with prefix "ui" corresponds to json data representation 
   * according to ui schema (pfdAppFormSchema).
   */
  transformToAfAppSchema = (uiModel) => {
      var afModel = {};
      var uiAppKeys  = Object.keys(uiModel);
      var uiPfds = uiModel['pfds'];
      var afPfds = {};
      var j=0;

      for (j=0; j<uiAppKeys.length; ++j){
        if(uiAppKeys[j] === 'pfds')
          continue;
        afModel[uiAppKeys[j]] = uiModel[uiAppKeys[j]];
      }
      for(j=0; j<uiPfds.length; ++j){
        var uiPfd = uiPfds[j].pfd;
        var pfdId = uiPfd.pfdID;
        afPfds[pfdId] = uiPfd;
      }
      afModel['pfds'] = afPfds;
    return afModel;
  }

  /* 
   * Convert afSchema (afModel) to uiSchema (uiModel). 
   * Variable with prefix "af" corresponds to json data representation 
   * according to af schema.
   * Variable with prefix "ui" corresponds to json data representation 
   * according to ui schema (pfdAppFormSchema).
   */
  transformToUiAppSchema = (afModel) => {
    var uiModel = {};
    var afAppKeys = Object.keys(afModel);
    var afPfds = afModel.pfds;
    var uiPfds = [];
    var afPfdNames = Object.keys(afPfds);
    var j=0;

    for(j=0; j<afAppKeys.length; ++j){
      if(afAppKeys[j] === 'pfds')
        continue;
      uiModel[afAppKeys[j]] = afModel[afAppKeys[j]];
    }

    for(j=0; j<afPfdNames.length; ++j){
      var afPfdName = afPfdNames[j];
      var afPfd = afPfds[afPfdName];
      var temp = {};
      temp['pfd'] = afPfd;
      uiPfds.push(temp);
    }
    uiModel['pfds'] = uiPfds;
    return uiModel;
  }


  putPfdApplication = async () => {
    const { history } = this.props;
    const { validationResult, model, tId, appId } = this.state;

    /* Check model contains no errors before continuing. */
    if ("valid" in validationResult && !validationResult.valid) {
      alert('Invalid request. Please correct any errors in the form and try again');
      return;
    };

    try {
      var afModel = this.transformToAfAppSchema(model);
      await axios.put(
          `${baseURL}/af/v1/pfd/transactions/${tId}/applications/${appId}`, 
          afModel);

      this.cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "success",
        snackbarMessage: "Successfully updated PFD Application data"
      }));

      await sleep(1000);

      /* Redirect user back to /pfd on success. */
      history.push('/pfd');
    } catch (error) {
      this.cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "error",
        snackbarMessage: error.toString(),
      }));
    }
  }

  deletePfdApplication = async() => {
    const { history } = this.props;
    const { tId, appId } = this.state;

    try {
      await axios.delete(
          `${baseURL}/af/v1/pfd/transactions/${tId}/applications/${appId}`);

      this.setState({
        snackbarOpen: true,
        snackbarVariant: "success",
        snackbarMessage: "Successfully deleted Pfd Application"
      });

      await sleep(1000);

      /* Redirect user back to /pfd on success. */
      history.push('/pfd');
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
    const result = utils.validateBySchema(Schema.pfdFormSchema, model);

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

    /* Validate the fields after the change to the model. */
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
    /*
     * Signal to cancel any pending async requests to prevent setting state on
     * an unmounted component
     */
    this.isMounted=false;
  }

  async componentDidMount() {
    this.isMounted = true;
    const { match } = this.props;

    try {
      const app = await this.getPfdApplication(match.params.tId, match.params.appId);

      const model = this.transformToUiAppSchema(app);
      this.cancelIfUnmounted(() => this.setState({
        loaded: true,
        tId: match.params.tId,
        appId: match.params.appId,
        model: model,
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
      history,
    } = this.props;

    const {
      model,
      showErrors,
      hasErrors,
      error,
      loaded,
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
            <Button onClick={() => history.push('/pfd')}>
              <KeyboardBackspaceOutlined /> Back to PFD Transactions
            </Button>
          </Grid>
          <Grid item>
            <Button
              onClick={this.putPfdApplication}
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
                Edit PFD Transaction: {this.state.tId}
                Application: {this.state.appId}
            </Typography>
          </Grid>

          <Grid item>
            <SchemaForm
              schema={Schema.pfdAppFormSchema}
              form={Schema.pfdAppForm}
              model={model}
              onModelChange={this.onModelChange}
              showErrors={showErrors}
            />
          </Grid>
        </Grid>
          <Grid
            container
            justify="center"
            alignItems="center"
          >
            <Grid item style={{ margin: "10%" }}>
              <Button
                onClick={this.deletePfdApplication}
                variant="outlined"
                color="secondary"
              >
                Delete Application {this.state.appId}
              </Button>
            </Grid>
          </Grid>
      </React.Fragment>
    );
  }
}

export default PacketFlowDescriptorAppModify;
