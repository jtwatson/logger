# logger
logger is an HTTP request logger that implements correlated logging to GCP via Logging REST API. Each HTTP request is logged as the parent with all event logs occurring during the request as child logs. This allows for easy viewing in GCP Log Explorer. The logs will also be correlated to Cloud Trace if you instrument your code with tracing.
