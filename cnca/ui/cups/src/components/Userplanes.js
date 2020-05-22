// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019-2020 Intel Corporation

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

const baseURL = (process.env.NODE_ENV === 'production') ? process.env.REACT_APP_CUPS_API : '/api';
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

class Userplanes extends Component {
  _isMounted = false;

  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
      hasError: false,
      userplanes: [],
    };
  }

  _cancelIfUnmounted = (action) => {
    if (this._isMounted) {
      action();
    }
  }

  _getUserplanes = async () => {
    try {
      const response = await axios.get(`${baseURL}/userplanes`);

      return response.data.userplanes || [];
    } catch (error) {
      console.error("Unable to get userplanes: " + error.toString());
      throw error;
    }
  }

  componentWillUnmount() {
    // Signal to cancel any pending async requests to prevent setting state
    // on an unmounted component.
    this._isMounted = false;
  }

  async componentDidMount() {
    this._isMounted = true;

    try {
      // Fetch userplanes.
      const userplanes = await this._getUserplanes() || [];

      // Update userplanes iff the component is mounted.
      this._cancelIfUnmounted(() => this.setState({
        loaded: true,
        userplanes: userplanes,
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
      classes,
      match,
      history,
    } = this.props;

    const {
      loaded,
      hasError,
      error,
      userplanes,
    } = this.state;

    if (!loaded) {
      return <Loader />;
    }

    if (hasError) {
      throw error;
    }

    const UserplaneTableRow = ({ match, history, item }) => {
      return (
        <TableRow>
          <TableCell>
            {item.id}
          </TableCell>
          <TableCell>
            {item.uuid}
          </TableCell>
          <TableCell>
            {item.function}
          </TableCell>
          <TableCell>
            <Button
              onClick={() => history.push(`${match.url}/${item.id}`)}
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
                Userplanes
              </Typography>
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
                <TableCell>ID</TableCell>
                <TableCell>UUID</TableCell>
                <TableCell>Function</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {
                userplanes.length === 0
                  ? <div>No userplanes to display.</div>
                  : userplanes.map(item => <UserplaneTableRow
                    key={item.id}
                    item={item}
                    history={history}
                    match={match} />)
              }
            </TableBody>
          </Table>
        </Paper >
      </div>
    );
  }
};

export default withStyles(styles)(Userplanes);
