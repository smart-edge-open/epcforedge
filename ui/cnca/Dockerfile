# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019-2020 Intel Corporation

# Build image with online dependencies
FROM node:lts-alpine AS cnca-ui-deps-image

# Create app directory & install dependencies for running the server
WORKDIR /usr/src/app
COPY ./server/package.json ./server/yarn.lock ./
RUN yarn install --pure-lockfile --production
COPY ./server ./

# Create app directory & install dependencies for building the app
WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn install --pure-lockfile --production


# Second stage - building the app with already existing dependencies
FROM cnca-ui-deps-image

# Set the AF API URL variable during build time
ARG REACT_APP_CNCA_AF_API
ENV REACT_APP_CNCA_AF_API=$REACT_APP_CNCA_AF_API
ARG REACT_APP_CNCA_5GOAM_API
ENV REACT_APP_CNCA_5GOAM_API=$REACT_APP_CNCA_5GOAM_API

# Set the Controller UI URL variable during build time
ARG LANDING_UI_URL
ENV REACT_APP_LANDING_UI_URL=$LANDING_UI_URL

# Instruct webpack to only import runtime scripts
ENV INLINE_RUNTIME_CHUNK=false

# Copy UI source & build
COPY . ./
RUN yarn build

WORKDIR /usr/src/app
RUN cp -r /app/build ./build

EXPOSE 80
CMD ["node", "server.js"]
