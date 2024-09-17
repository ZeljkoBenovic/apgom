# APGOM (AsteriskPrometheusGOMetrics)
`apgom` is an Asterisk Prometheus scraper that can fetch Asterisk information and expose and endpoint 
stable for Prometheus scraper to consume.

## Usage
Just run `apgom` as a service on an Asterisk box,
and it will automagically expose Prometheus metrics on `/metrics` endpoint.     
It scrapes metrics from the Asterisk AMI API directly, so it needs to be enabled on the Asterisk server.

### Configuration flags
* `-host` metrics listen host (default `0.0.0.0`)
* `-port` metrics listen port (default `3000`)
* `-path` metrics listen url path (default `/metrics`)
* `-ami-host` asterisk hostname or ip address (default `localhost`)
* `-ami-user` asterisk AMI username (default `admin`)
* `-ami-pass` asterisk AMI password (required)

## Metrics
* `asterisk_calls_total` - the total number of processed calls (Gauge)
* `asterisk_calls_active` - the number of currently active calls (Gauge)
* `asterisk_extensions_total` - the total number of extensions regardless of availability (Gauge)
* `asterisk_extensions_available` - the number of available extensions (Gauge)
* `asterisk_extensions_unavailable` - the number of unavailable extensions (Gauge)
* `asterisk_registries_total` - the total number of registries/trunks regardless of their status (Gauge)
* `asterisk_registries_active` - the number of active registries/trunks (Gauge)
* `asterisk_registries_down` - the number of inactive registries/trunks (Gauge)
