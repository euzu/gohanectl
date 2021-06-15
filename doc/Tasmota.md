#Tasmota
## Tasmota messages 
> LWT
```
 Online | Offline
```
> SATE
```json
   {
     "Time":"2020-12-27T12:36:09",
     "Uptime":"5T14:15:34",
     "UptimeSec":483334,
     "Heap":26,
     "SleepMode":"Dynamic",
     "Sleep":50,
     "LoadAvg":19,
     "MqttCount":34,
     "POWER":"OFF",
     "Wifi":{"AP":1,"SSId":"sarcofwan","BSSId":"E0:28:6D:14:AB:94","Channel":1,"RSSI":100,"Signal":-44,"LinkCount":9,"Downtime":"0T00:00:19"}
   }
```
> SENSOR
```json
  {
     "Time":"2020-12-27T13:11:09",
     "ENERGY": {
        "TotalStartTime":"2020-12-21T14:57:55",
        "Total":0.001,
        "Yesterday":0.000,
        "Today":0.001,
        "Period":1,
        "Power":11,
        "ApparentPower":20,
        "ReactivePower":17,
        "Factor":0.54,
        "Voltage":294,
        "Current":0.068
       }
    }
```

> STATE
```json
{
    "Time":"2020-12-27T15:46:16",
    "Uptime":"0T00:05:09",
    "UptimeSec":309,
    "Heap":25,
    "SleepMode":"Dynamic",
    "Sleep":50,
    "LoadAvg":19,
    "MqttCount":1,
    "POWER":"ON",
    "Wifi":{"AP":1,"SSId":"sarcofwan", "BSSId":"E0:28:6D:14:AB:94","Channel":1,"RSSI":94,"Signal":-53,"LinkCount":1,"Downtime":"0T00:00:02"}
}
```

## devices.yml config

```yaml
  - type: light
    device_key: licht-arbeitszimmer
    caption: Arbeitszimmer
    listen_topics:
      - topic: tele/light_az2
        template: tasmota_tele
      - topic: stat/light_az2
        template: tasmota_stat
    command_topic: cmnd/light_az2/Power
    payload_on: 'on'
    payload_off: 'off'
    status_topic: cmnd/light_az2/Power
    status_payload: ''
```
