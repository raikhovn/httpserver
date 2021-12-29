# syntax=docker/dockerfile:1

# Base image
FROM alpine:latest as base

# Build image
FROM golang:1.16-alpine as build

## Run build stage 1##

# Create build directory
RUN mkdir /build
# Copy all source code from folder/subfolders to a /builder folder
ADD . /build
# Switch current folder to a /builder folder
WORKDIR /build
# RUn the go build to build httpfileserver binary
RUN go build -o httpfileserver .

## Copy build stage 2##

# Picking up from the prev sdtage "base"
FROM base as final
# Create and set new /app folder
RUN mkdir /app
WORKDIR /app
# Copy from the build image using local image tag "build"
COPY --from=build /build .



# Expose container port
EXPOSE 8080

CMD [ "/app/httpfileserver" ]