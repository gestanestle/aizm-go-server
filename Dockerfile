# syntax=docker/dockerfile:1

FROM golang:1.22

# Set destination for COPY
WORKDIR /app

ARG DATABASE_URL
ARG MQTT_HOST
ARG MQTT_PORT
ARG MQTT_USER
ARG MQTT_PASS
ARG MQTT_TOPIC

ENV DATABASE_URL ${DATABASE_URL}
ENV MQTT_HOST ${MQTT_HOST}
ENV MQTT_PORT ${MQTT_PORT}
ENV MQTT_USER ${MQTT_USER}
ENV MQTT_PASS ${MQTT_PASS}
ENV MQTT_TOPIC ${MQTT_TOPIC}

# Download Go modules
COPY . ./
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

# Run
CMD ["/docker-gs-ping"]
