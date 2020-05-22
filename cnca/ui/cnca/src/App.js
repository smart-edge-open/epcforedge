/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019-2020 Intel Corporation
 */

import React, { Component } from 'react';
import { Route, Redirect, Switch } from 'react-router-dom';
import { BrowserRouter as Router } from "react-router-dom";
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Subscriptions from './components/Subscriptions';
import Subscription from './components/Subscription';
import SubscriptionModify from './components/SubscriptionModify';
import Services from './components/Services.js';
import Service from './components/Service.js';
import PacketFlowDescriptors from './components/PacketFlowDescriptors.js';
import PacketFlowDescriptor from './components/PacketFlowDescriptor.js';
import PacketFlowDescriptorAppModify from './components/PacketFlowDescriptorAppModify.js';
import Header from './components/Header';
import ErrorBoundary from './components/ErrorBoundary';

const useStyles = theme => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    minHeight: '100vh',
  },
  header: {
    padding: 20,
    flexGrow: 1,
  },
  main: {
    width: 'auto',
    marginTop: theme.spacing(14),
    marginBottom: theme.spacing(2),
    marginLeft: theme.spacing(2),
    marginRight: theme.spacing(2),
    [theme.breakpoints.up(600 + theme.spacing(4))]: {
      width: 800,
      marginLeft: 'auto',
      marginRight: 'auto',
    },
  },
});

class App extends Component {
  render() {
    const { classes } = this.props;

    return (
      <Router>
        <div>
          <header>
            <Header className={classes.header} />
          </header>
          <main className={classes.main}>
            <ErrorBoundary>
              <Switch>
                <Route
                  exact
                  path="/"
                  render={() => <Redirect to="/services" />}
                />
                <Route
                  exact
                  path="/services"
                  component={Services}
                />
                <Route
                  exact
                  path="/services/create"
                  render={(props) => <Service {...props} createMode={true}/>}
                />
                <Route
                  exact
                  path="/services/:id"
                  component={Service}
                />
                <Route
                  exact
                  path="/subscriptions"
                  component={Subscriptions}
                />
                <Route
                  exact
                  path="/subscriptions/create"
                  render={(props) => <Subscription {...props} createMode={true} />}
                />
                <Route
                  exact
                  path="/subscriptions/edit/:id"
                  component={Subscription}
                />
                <Route
                  exact
                  path="/subscriptions/patch/:id"
                  component={SubscriptionModify}
                />
                <Route
                  exact
                  path="/pfd"
                  component={PacketFlowDescriptors}
                />
                <Route
                  exact
                  path="/pfd/transactions/create"
                  render={(props) => <PacketFlowDescriptor {...props} createMode={true} />}
                />

                 <Route
                  exact
                  path="/pfd/transactions/:tId"
                  component={PacketFlowDescriptor}
                />
                <Route
                  exact
                  path="/pfd/transactions/:tId/applications/:appId"
                  render={(props) => <PacketFlowDescriptorAppModify {...props} />}
                />
 
 
                <Route
                  render={() => <span>404 Not Found</span>}
                />
              </Switch>
            </ErrorBoundary>
          </main>
        </div>
      </Router>
    )
  }
}

App.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(useStyles)(App);
