// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

import React, { Component } from 'react';
import axios from 'axios';
import * as Schema from './schema/UserplaneForm';
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

const baseURL = (process.env.NODE_ENV === 'production') ? process.env.REACT_APP_CUPS_API : '/api';

const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms))
class Userplane extends Component {
  _isMounted = false;

  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      model: {},
      configModel: {},
      selectorsModel: {},
      entitlementsModel: {},
      validationResult: {},
      showErrors: false,
    };
  }

  _getUserplane = async (id) => {
    try {
      const response = await axios.get(`${baseURL}/userplanes/${id}`);

      return response.data;
    } catch (error) {
      console.error("Unable to get userplane: " + error.toString());
      throw error;
    }
  }

  _postUserplane = async () => {
    const { history } = this.props;
    const { validationResult, model } = this.state;

    // Ensure the model is valid and contains no errors before proceeding.
    if ("valid" in validationResult && !validationResult.valid) {
      alert('Invalid request. Please correct any errors in the form and try again');
      return;
    }

    try {
      await axios.post(`${baseURL}/userplanes`, model);

      this._cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "success",
        snackbarMessage: "Successfully created userplane."
      }));

      await sleep(1000);

      // Redirect user back to /userplanes on success.
      history.push('/userplanes');
    } catch (error) {
      this._cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "error",
        snackbarMessage: error.toString(),
      }));
    }
  }

  _patchUserplane = async () => {
    const { history } = this.props;
    const { validationResult, model } = this.state;

    // Ensure the model is valid and contains no errors before proceeding.
    if ("valid" in validationResult && !validationResult.valid) {
      alert('Invalid request. Please correct any errors in the form and try again');
      return;
    }

    try {
      await axios.patch(`${baseURL}/userplanes/${model.id}`, model);

      this._cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "success",
        snackbarMessage: "Successfully updated userplane."
      }));

      await sleep(1000);

      // Redirect user back to /userplanes on success.
      history.push('/userplanes');
    } catch (error) {
      this._cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "error",
        snackbarMessage: error.toString(),
      }));
    }
  }

  _deleteUserplane = async () => {
    const { history } = this.props;
    const { model } = this.state;

    try {
      await axios.delete(`${baseURL}/userplanes/${model.id}`);

      this.setState({
        snackbarOpen: true,
        snackbarVariant: "success",
        snackbarMessage: "Successfully deleted userplane."
      });

      await sleep(1000);

      // Redirect user back to /userplanes on success.
      history.push('/userplanes');
    } catch (error) {
      // Update error iff the component is mounted.
      this._cancelIfUnmounted(() => this.setState({
        snackbarOpen: true,
        snackbarVariant: "error",
        snackbarMessage: error.toString(),
      }));
    }
  }

  _onValidate = () => {
    const { model } = this.state;

    const result = utils.validateBySchema(Schema.userplaneFormSchema, model);

    this.setState({
      validationResult: result,
      showErrors: true,
    });
  };

  _onModelChange = (key, val) => {
    const { model } = this.state;

    const newModel = model;

    utils.selectOrSet(key, newModel, val);

    this.setState({ model: newModel });

    // Validate fields on each model change.
    this._onValidate();
  }

  _onConfigModelChange = (key, val) => {
    const { configModel } = this.state;

    const newModel = configModel;

    utils.selectOrSet(key, newModel, val);

    this.setState({ configModel: newModel });

    this._mergeSubModels();
  }

  _onSelectorsModelChange = (key, val) => {
    const { selectorsModel } = this.state;

    const newModel = selectorsModel;

    utils.selectOrSet(key, newModel, val);

    this.setState({ selectorsModel: newModel });

    this._mergeSubModels();
  }

  _onEntitlementsModelChange = (key, val) => {
    const { entitlementsModel } = this.state;

    const newModel = entitlementsModel;

    utils.selectOrSet(key, newModel, val);

    this.setState({ entitlementsModel: newModel });

    this._mergeSubModels();
  }

  _mergeSubModels = () => {
    const {
      model,
      configModel,
      selectorsModel,
      entitlementsModel,
    } = this.state;

    this.setState({
      model: {
        ...model,
        ...configModel,
        ...selectorsModel,
        ...entitlementsModel,
      },
    });
  }

  _cancelIfUnmounted = (action) => {
    if (this._isMounted) {
      action();
    }
  }

  handleClose = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }

    this.setState({ snackbarOpen: false });
  };

  componentWillUnmount() {
    // Signal to cancel any pending async requests to prevent setting state on
    // an unmounted component.
    this._isMounted = false;
  }

  async componentDidMount() {
    this._isMounted = true;

    const { createMode, match } = this.props;

    // When the user is creating a new userplane, skip the fetch request.
    if (createMode) {
      this.setState({ loaded: true });
      return;
    }

    try {
      const userplane = await this._getUserplane(match.params.id);

      this._cancelIfUnmounted(() => this.setState({
        loaded: true,
        model: userplane,
        configModel: { config: userplane.config },
        selectorsModel: { selectors: userplane.selectors },
        entitlementsModel: { entitlements: userplane.entitlements },
      }));
    } catch (error) {
      // Update error iff the component is mounted.
      this._cancelIfUnmounted(() => this.setState({
        loaded: true,
        hasError: true,
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
      configModel,
      selectorsModel,
      entitlementsModel,
      showErrors,
      hasError,
      error,
      loaded,
      snackbarOpen,
      snackbarVariant,
      snackbarMessage,
    } = this.state;

    if (!loaded) {
      return <Loader />
    }

    if (hasError) {
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
            <Button onClick={() => history.push('/userplanes')}>
              <KeyboardBackspaceOutlined /> Back to Userplanes
            </Button>
          </Grid>
          <Grid item>
            <Button
              onClick={(createMode) ? this._postUserplane : this._patchUserplane}
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
                  ? 'Create Userplane'
                  : `Userplane ID ${model.id}`
              }
            </Typography>
          </Grid>

          <Grid item>
            <SchemaForm
              schema={Schema.userplaneFormSchema}
              form={Schema.userplaneForm}
              model={model}
              onModelChange={this._onModelChange}
              showErrors={showErrors}
            />
          </Grid>

          <Grid item xs={4}>
            <SchemaForm
              schema={Schema.configFormSchema}
              form={Schema.configForm}
              model={configModel}
              onModelChange={this._onConfigModelChange}
              showErrors={showErrors}
            />
          </Grid>

          <Grid item xs={6}>
            <SchemaForm
              schema={Schema.selectorsFormSchema}
              form={Schema.selectorsForm}
              model={selectorsModel}
              onModelChange={this._onSelectorsModelChange}
              showErrors={showErrors}
            />
          </Grid>

          <Grid item xs={6} m={100}>
            <SchemaForm
              schema={Schema.entitlementsFormSchema}
              form={Schema.entitlementsForm}
              model={entitlementsModel}
              onModelChange={this._onEntitlementsModelChange}
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
                onClick={this._deleteUserplane}
                variant="outlined"
                color="secondary"
              >
                Delete Userplane
            </Button>
            </Grid>
          </Grid>
        }
      </React.Fragment>
    );
  }
}

export default Userplane;
