# Shelly 1

## Messages

> relay/0
```
  off | on
```
> input/0
```
  0
```
> input_event/0
```json
{
  "event":"S",
  "event_cnt":73
}
```

## devices.yml config

```yaml
- type: light
  device_key: licht-flur
  caption: Flur
  listen_topics:
    - topic: shellies/shelly1-483FDA91BCEF
      template: shelly_1
  command_topic: shellies/shelly1-483FDA91BCEF/relay/0/command
  payload_on: 'on'
  payload_off: 'off'
  status_topic: shellies/shelly1-483FDA91BCEF/relay/0
  status_payload: ''
```
