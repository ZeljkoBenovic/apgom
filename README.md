# APGOM (AsteriskPrometheusGOMetrics)
`apgom` is an Asterisk Prometheus scraper that can fetch Asterisk information and expose and endpoint 
stable for Prometheus scraper to consume.

## Usage
Just run `apgom` as a service on an Asterisk box, and it will automagically fetch Asterisk metrics.

## Metrics
* `asterisk_sip_channels` - active channels (Gauge)