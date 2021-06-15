# Homematic

## Homematic thermostat

> _ClimRT_tr
```json
{
  "TEMP":"22.7",
  "DESIRED":"21.0",
  "VALVE":"4"
}
```

## Homematic Licht

```json
{ "STATUS":"off"}
```

### Example devices.yml configuration for homematic light
```yaml
- type: light
  device_key: licht-flur-eingang
  caption: Flur Eingang
  listen_topics:
    - topic: fhem/LichtFlurEingang
      template: hm_light
  command_topic: cmnd/fhem
  payload_on: set LichtFlurEingang on
  payload_off: set LichtFlurEingang off
  status_topic: cmnd/fhem
  status_payload: list LichtFlurEingang state
```
