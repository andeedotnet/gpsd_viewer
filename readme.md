# GPSD Viewer

A lightweight Go server that connects to a `gpsd` daemon, processes live GPS data, and serves it as JSON via a web interface.

## 🌍 Features

- Connects to a local or remote `gpsd` daemon
- Handles `TPV`, `SKY`, and `AIS` GPSD message types
- Provides a JSON endpoint (/json) for external tools or frontend use

## 📦 Requirements

- Tested on Go 1.23
- A running [gpsd](https://gpsd.io/) server (e.g., on localhost or remote device)

## 🚀 Usage

./gpsd_viewer -gpsd localhost:2947 -p 9000


### Launch options
| option| description | default |
|---|---|---|
| -gpsd | address and port of gpsd daemon. | localhost:2947 |
| -p    | port for gpsd_viewer web interface.  | 9000 |
