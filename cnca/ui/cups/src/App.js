// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019-2020 Intel Corporation

import React, { Component } from 'react';
import { Route, Redirect, Switch } from 'react-router-dom';
import { BrowserRouter as Router } from "react-router-dom";
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Userplanes from './components/Userplanes';
import Userplane from './components/Userplane';
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
                  render={() => <Redirect to="/userplanes" />}
                />
                <Route
                  exact
                  path="/userplanes"
                  component={Userplanes}
                />
                <Route
                  exact
                  path="/userplanes/create"
                  render={(props) => <Userplane {...props} createMode={true} />}
                />
                <Route
                  exact
                  path="/userplanes/:id"
                  component={Userplane}
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
