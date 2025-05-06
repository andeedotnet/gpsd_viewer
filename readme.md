# GPSD Viewer

A lightweight GPSD Client that connects to a `gpsd` daemon, processes live GPS data, and serves it as JSON via a web interface.

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
| -e    | use embedded web interface (true / false)  | true |


### Example
The project was developed just for fun, as I was playing around with my 5G router (GL.iNet GL-X3000) and wanted to access the GPS data. Here's how to enable GPSD on the router:

1. Send manual commands on the admin page for the modem:
```
AT+QGPSCFG="autogps",1
AT+QGPS=1
```

2. Install gpsd via the Plugins page

3. Configure gpsd
```
uci set gpsd.core.device='/dev/mhi_LOOPBACK'
uci set gpsd.core.listen_globally='1'
uci set gpsd.core.enabled='1'

/etc/init.d/gpsd enable
/etc/init.d/gpsd start
```

4. Optional: Enable Glonass, Galileo, Beidou NMEA sentence output
```
AT+QGPSCFG="glonassnmeatype",1
AT+QGPSCFG="galileonmeatype",1
AT+QGPSCFG="beidounmeatype",1
```

5. Copy gpsd_viewer_linux_arm64 from releases and run
```
./gpsd_viewer_linux_arm64 -gpsd localhost:2947 -p 9000
```

If you want to use your own web app, you can simply place it under /static and start gpsd_viewer with ```-e false```.