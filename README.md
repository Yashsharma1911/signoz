## Introduction
This application shows an implementation of metrics instrumentation in a Go application using OpenTelemetry (Otel sdk) and uses SigNoz as an observability tool to show the visualization of generated metrics. It is a simple Go server with sample endpoints designed to test and generate Otel data, which is then sent to SigNoz Cloud. Application generates three types of metrics Counter, Histogram, and Gauge for error requests, request duration and dynamic data.
SigNoz supports a wide variety of telemetry signal which makes this application extensible and can be customised to test various metrics like DB calls, latencies and more. 

### Application Dependencies
- OpenTelemetry Metrics SDK
- SigNoz Cloud for visualization
- OTLP Exporter for metrics transmission

Find all details of setup at - https://docs.google.com/document/d/19By5JTaqUk5bQv5cu1Wb76yChbHx8CIbEzuCXY5Os9Y/edit?tab=t.0
