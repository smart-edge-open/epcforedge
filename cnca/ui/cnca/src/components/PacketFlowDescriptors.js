/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2020 Intel Corporation
 */

import React, { Component } from 'react';
import axios from 'axios';
import Loader from './Loader';
import SnackbarWrapper from './SnackbarWrapper';
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
  Snackbar,
} from '@material-ui/core';

const baseURL = (process.env.NODE_ENV === 'production') ? 
                    process.env.REACT_APP_CNCA_AF_API : '/api';
const LANDING_URL = process.env.REACT_APP_LANDING_UI_URL;

const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms))
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

class PacketFlowDescriptors extends Component {
  isMounted = false;

  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      hasError: false,
      transactions: [],
      selectOpts: [],
      transactionId: "0",
      appId: "0",
      optsLength: 0,
    }
  }

  cancelIfUnmounted = (action) => {
    if (this.isMounted) {
      action();
    }
  }

  /* Retrive PFD transactions from AF Server */
  getTransactions = async() => {
    try {
      const response = await axios.get(`${baseURL}/af/v1/pfd/transactions`);

      return response.data || [];
    } catch (error) {
      console.error("Unable to get PFD Transactions: " + error.toString());
      throw error;
    }
  }

  /*
   * This function is invoked on clicking edit button.
   * It checks if transaction is to be edited or application is to be edited
   * and invoke the url to edit transaction/application.
   */
  editPfd = async(elementId, transactionId) => {
    const { history } = this.props;
    var selectedApp = document.getElementById(elementId).value;
    if(selectedApp === 'transaction')
      history.push(`/pfd/transactions/${transactionId}`);
    else
      history.push(
          `/pfd/transactions/${transactionId}/applications/${selectedApp}`);
  }

  /*
   * This function is called on clicking delete button.
   * It checks if transaction is to be deleted or application is to be deleted
   * and send request to delete transaction/application.
   */
  deletePfd = async(elementId, transactionId) => {
    const { history } = this.props;

    try {
      var selectedApp = document.getElementById(elementId).value;
      if(selectedApp === 'transaction')
        await axios.delete(
            `${baseURL}/af/v1/pfd/transactions/${transactionId}`);
      else
        await axios.delete(
            `${baseURL}/af/v1/pfd/transactions/${transactionId}/applications/${selectedApp}`);
 
      this.setState({
        snackbarOpen: true,
        snackbarVariant: "success",
        snackbarMessage: "Successfully deleted pfd transaction"
      });

      await sleep(1000);
      // Redirect back to /pfd to refresh the table of pfd transactions
      window.location.reload();
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
      // GET transactions
      const transactions = await this.getTransactions() || [];

      this.cancelIfUnmounted(() => this.setState({
        loaded: true,
        transactions: transactions,
      }));

      /*
       * Populate the applications in the table select element from the global
       * array {this.state.selectOpts}, generated in TransactionTableRow. 
       */
      var currSelectOpts = this.state.selectOpts;
      for(var i=0; i<currSelectOpts.length; ++i){
        var htmlSelect = document.getElementById(currSelectOpts[i].selectId);
        var appNames = currSelectOpts[i].appNames;
        var appIds = currSelectOpts[i].appIds;
        for(var j=0; j<appNames.length; ++j){
          htmlSelect.options[htmlSelect.options.length] = 
                                          new Option(appIds[j], appIds[j]);
        }
      }
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
      transactions,
      snackbarOpen,
      snackbarVariant,
      snackbarMessage,
    } = this.state;

    if (!loaded) {
      return <Loader />;
    }

    if (hasErrors) {
      throw error;
    }

    const TransactionTableRow = ({ match, item }) => {
      /* Get Transaction Id from .self element in JSON */
      var url = item.self;
      const selfArray = url.split('/');
      const transactionId = selfArray[selfArray.length - 1];

      /*
       * From the current json (application) extract the application Ids and
       * fill them in selectObj array, which are used later in
       * componentDidMount for representing applications.
       */
      var appKeys = Object.keys(item.pfdDatas);
      var appIds = [];
      var selectId = "transaction";
      var buttonId1 = "button1-";
      var buttonId2 = "button2-";
      var selectObj = {};

      for(var i=0; i<appKeys.length; ++i){
        appIds.push(item.pfdDatas[appKeys[i]].externalAppID);
      }

      selectId += transactionId;
      buttonId1 += selectId;
      buttonId2 += selectId;

      selectObj['tid'] = transactionId;
      selectObj['appNames'] = appKeys;
      selectObj['appIds'] = appIds;
      selectObj['selectId'] = selectId;
      this.state.selectOpts.push(selectObj);

      return (
        <TableRow>
          <TableCell>
            {transactionId}
          </TableCell>
          <TableCell>
            <select id={selectId}>
              <option value="transaction">Apps</option>
            </select>
          </TableCell>
          <TableCell>
            <Button
              id={buttonId1}
              onClick={() => this.editPfd(selectId, transactionId)}
              variant="outlined"
            >
              Edit
            </Button>
            <Button
              id={buttonId2}
              onClick={() => this.deletePfd(selectId, transactionId)}
              variant="outlined"
            >
              Delete
            </Button>
          </TableCell>
        </TableRow>
      );
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
                PFD Transactions
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
                onClick={() => history.push(`${match.url}/transactions/create`)}
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
                <TableCell>PFD Transaction ID</TableCell>
                <TableCell>App ID</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {
                transactions.length === 0
                  ? <div>No PFD Registered.</div>
                  : transactions.map(item => <TransactionTableRow
                    key={item.pfdDatas}
                    item={item}
                    match={match} />)
              }
            </TableBody>
          </Table>
        </Paper>
      </div>
      </React.Fragment>
    );
  }
};

export default withStyles(styles)(PacketFlowDescriptors);
