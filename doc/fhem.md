# FHEM example config for mqtt

### MQTT Server
```
define mqttBroker MQTT2_SERVER 1883 global
setuuid mqttBroker 5fe0b68c-f33f-6f28-4723-a4c71f5000fa20a5
```

### Command execution for fhem devices
to execute commands for fhem devices we listen to the topic 
> fhem/cmnd

commands can be executed by sending the command as payload to this topic
> set LichtFlur on
> set HeizungWohnzimmer_RT_tr desired-temp 22.5

```
define SYS_MQTT MQTT2_DEVICE
setuuid SYS_MQTT 5fe48f1b-f33f-6f28-8572-00ae68caec121dda
attr SYS_MQTT IODev mqttBroker
attr SYS_MQTT alias MQTT_COMMAND
attr SYS_MQTT readingList fhem/cmnd:.* cmnd
attr SYS_MQTT room MQTT
define FileLog_SYS_MQTT FileLog ./log/SYS_MQTT-%Y.log SYS_MQTT
setuuid FileLog_SYS_MQTT 5fe48f1b-f33f-6f28-1888-cf38fbe5f6f98513
attr FileLog_SYS_MQTT logtype text
attr FileLog_SYS_MQTT room MQTT
define n_SYS_MQTT_cmnd notify SYS_MQTT:cmnd:.* {\
   if ($EVENT =~ qr/.*?: (.*)/p) {\
      my $cmnd = $1;;\
      Log3($NAME, 5, "executed mqtt command: " . $cmnd);;\
      fhem($cmnd);;\
    }\
}
setuuid n_SYS_MQTT_cmnd 5fe46d70-f33f-6f28-d777-ad0b5d773d5121f3
```

### Mqtt device definition
fhem can autodiscover mqtt devices with MQTT2

### Sample thermostat
```
define mqtt_Heizung notify (HeizungArbeitszimmer_ClimRT_tr|HeizungKinderzimmer_ClimRT_tr|HeizungSchlafzimmer_ClimRT_tr):* {\
  if ($EVENT =~ /T:\s+([0-9\.]+)\s+desired:\s+([0-9\.]+)\s+valve:\s+([0-9\.]+)/) {\
      fhem "set mqttBroker publish fhem/$NAME {\"temp\":\"$1\", \"desired\":\"$2\", \"valve\":\"$3\"}";;\
   }\
}
setuuid mqtt_Heizung 5fe475a0-f33f-6f28-6f3e-ca9d731cb000f4c1
attr mqtt_Heizung room Heizung,MQTT
```

