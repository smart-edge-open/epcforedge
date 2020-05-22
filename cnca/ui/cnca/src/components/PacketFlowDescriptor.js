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
class PacketFlowDescriptor extends Component {
  isMounted = false;

  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      model: {},
      validationResult: {},
      showErrors: false,
      tId: 0,
    };
  }

  /* Retrive PFD transactions from AF Server */
  getTransaction = async (tId) => {
    try {
      const response = await axios.get(
          `${baseURL}/af/v1/pfd/transactions/${tId}`);

      return response.data;
    } catch (error) {
      console.error("Unable to get pfd transactions: " + error.toString());
      throw error;
    }
  }

  /* 
   * Convert uiSchema (uiModel) to afSchema (afModel). 
   * Variables with prefix "af" corresponds to json data representation 
   * according to af schema.
   * Variables with prefix "ui" corresponds to json data representation 
   * according to ui schema (pfdAppFormSchema).
   */
  transformToAfSchema = (uiModel) => {
    var afModel = {};
    var pfdDatas = uiModel.pfdDatas;
    var newPfdDatas = {};
    var uiModelKeys = Object.keys(uiModel);
    var i=0, j=0;

    for(i=0; i<uiModelKeys.length; ++i){
      if(uiModelKeys[i] === 'pfdDatas')
        continue;
      afModel[uiModelKeys[i]] = uiModel[uiModelKeys[i]];
    }


    for (i=0; i<pfdDatas.length; ++i){
      var origApp = pfdDatas[i].apps;
      var newApp = {};
      var appName = origApp.externalAppID;
      var appKeys  = Object.keys(origApp);
      var pfds = origApp['pfds'];
      var newPfds = {};

      for (j=0; j<appKeys.length; ++j){
        if(appKeys[j] === 'pfds')
          continue;
        newApp[appKeys[j]] = origApp[appKeys[j]];
      }
      for(j=0; j<pfds.length; ++j){
        var pfd = pfds[j].pfd;
        var pfdId = pfd.pfdID;
        newPfds[pfdId] = pfd;
      }
      newApp['pfds'] = newPfds;
      newPfdDatas[appName] = newApp;
    }
    afModel.pfdDatas = newPfdDatas;
    return afModel;
  }

  /* 
   * Convert afSchema (afModel) to uiSchema (uiModel). 
   * Variable with prefix "af" corresponds to json data representation 
   * according to af schema.
   * Variable with prefix "ui" corresponds to json data representation 
   * according to ui schema (pfdAppFormSchema).
   */
  transformToUiSchema = (afModel) => {
    var uiModel = {};
    var uiApps = [];
    var afModelKeys = Object.keys(afModel);
    var i=0, j=0;

    for(i=0; i<afModelKeys.length; ++i){
      if(afModelKeys[i] === 'pfdDatas')
        continue;
      uiModel[afModelKeys[i]] = afModel[afModelKeys[i]];
    }

    var afAppNames = Object.keys(afModel.pfdDatas);
    for(i=0; i<afAppNames.length; ++i){
      var afAppName = afAppNames[i];
      var afApp = afModel.pfdDatas[afAppName];
      var afAppKeys = Object.keys(afApp);
      var afPfds = afApp.pfds;
      var uiApp = {};
      var uiPfds = [];
      var afPfdNames = Object.keys(afPfds);
      var tempApp = {};

      for(j=0; j<afAppKeys.length; ++j){
        if(afAppKeys[j] === 'pfds')
          continue;
        uiApp[afAppKeys[j]] = afApp[afAppKeys[j]];
      }

      for(j=0; j<afPfdNames.length; ++j){
        var afPfdName = afPfdNames[j];
        var afPfd = afPfds[afPfdName];
        var temp = {};
        temp['pfd'] = afPfd;
        uiPfds.push(temp);
      }
      uiApp['pfds'] = uiPfds;
      tempApp['apps'] = uiApp;
      uiApps.push(tempApp);
    }
    uiModel['pfdDatas'] = uiApps;
    return uiModel;
  }

  postPfdTransactions = async() => {
    const { history } = this.props;
    const { validationResult, model } = this.state;

    /* Check model contains no errors before continuing. */
    if ("valid" in validationResult && !validationResult.valid) {
      alert('Invalid request. Please correct any errors in the form and try again');
      return;
    }

    try {
      /* Convert current model to afModel and send post request. */
      const afModel = this.transformToAfSchema(model);
      await axios.post(`${baseURL}/af/v1/pfd/transactions`, afModel);

      this.cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "success",
        snackbarMessage: "Successfully created transaction"
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

  putPfdTransaction = async () => {
    const { history } = this.props;
    const { validationResult, model, tId } = this.state;

    /* Check model contains no errors before continuing */
    if ("valid" in validationResult && !validationResult.valid) {
      alert('Invalid request. Please correct any errors in the form and try again');
      return;
    };

    try {
      /* Convert current model to afModel and send put request. */
      const afModel = this.transformToAfSchema(model);
      await axios.put(`${baseURL}/af/v1/pfd/transactions/${tId}`, afModel);

      this.cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "success",
        snackbarMessage: "Successfully updated transaction"
      }));

      await sleep(1000);

      /* Redirect user back to /pfd on success */
      history.push('/pfd');
    } catch (error) {
      this.cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "error",
        snackbarMessage: error.toString(),
      }));
    }
  }

  deleteTransaction = async() => {
    const { history } = this.props;
    const { tId } = this.state;

    try {
      await axios.delete(`${baseURL}/af/v1/pfd/transactions/${tId}`);

      this.setState({
        snackbarOpen: true,
        snackbarVariant: "success",
        snackbarMessage: "Successfully deleted transaction"
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

    /* Validate the fields after the change to the model.*/
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
    const { createMode, match } = this.props;

    /* When the user is creating a new userplane, skip the fetch request */
    if (createMode) {
      this.setState({ loaded: true });
      return;
    }

    try {
      /*
       * Retreive Transaction from AF server and convert recieved json to 
       * ui json model "pfdFormSchema".
       */
      const transaction = await this.getTransaction(match.params.tId);
      const uiModel = this.transformToUiSchema(transaction);

      this.cancelIfUnmounted(() => this.setState({
        loaded: true,
        tId: match.params.tId,
        model: uiModel,
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
      tId,
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
              onClick={(createMode) ? 
                this.postPfdTransactions : this.putPfdTransaction}
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
                  ? 'Create PFD Transaction'
                  : `Transaction ID ${tId}`
              }
            </Typography>
          </Grid>

          <Grid item>
            <SchemaForm
              schema={Schema.pfdFormSchema}
              form={Schema.pfdForm}
              model={model}
              onModelChange={this.onModelChange}
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
                onClick={this.deleteTransaction}
                variant="outlined"
                color="secondary"
              >
                Delete PFD
              </Button>
            </Grid>
          </Grid>
        }
      </React.Fragment>
    );
  }
}

export default PacketFlowDescriptor;
