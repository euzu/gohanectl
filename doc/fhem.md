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

### Sample tasmota socket:

```
define MQTT2_DVES_1791A1 MQTT2_DEVICE DVES_1791A1
setuuid MQTT2_DVES_1791A1 5fe0b743-f33f-6f28-4515-985b423a91a5b17e
attr MQTT2_DVES_1791A1 IODev mqttBroker
attr MQTT2_DVES_1791A1 alias SocketPrinter
attr MQTT2_DVES_1791A1 autocreate 0
attr MQTT2_DVES_1791A1 comment NOTE: For on-for-timer SetExtensions are used. You may add on-for-timer option running on the device. The following is limited to 1h max duration, but will not affe>
attr MQTT2_DVES_1791A1 icon hue_filled_outlet
attr MQTT2_DVES_1791A1 jsonMap POWER1:0 POWER2:0 POWER3:0 POWER4:0 Dimmer:0 Channel_0:0 Channel_1:0 Channel_2:0 Channel_3:0 Channel_4:0 HSBColor:0 Color:0
attr MQTT2_DVES_1791A1 model tasmota_basic_state_power1
attr MQTT2_DVES_1791A1 readingList tele/socket_printer/LWT:.* LWT\
  tele/socket_printer/STATE:.* { json2nameValue($EVENT,'',$JSONMAP) }\
  tele/socket_printer/SENSOR:.* { json2nameValue($EVENT,'',$JSONMAP) }\
  tele/socket_printer/INFO.:.* { json2nameValue($EVENT,'',$JSONMAP) }\
  tele/socket_printer/UPTIME:.* { json2nameValue($EVENT,'',$JSONMAP) }\
  stat/socket_printer/POWER1:.* state\
  stat/socket_printer/RESULT:.* { json2nameValue($EVENT,'',$JSONMAP) }
attr MQTT2_DVES_1791A1 room Steckdose
attr MQTT2_DVES_1791A1 setList off:noArg    cmnd/socket_printer/POWER1 0\
  on:noArg     cmnd/socket_printer/POWER1 1\
  toggle:noArg cmnd/socket_printer/POWER1 2\
  setOtaUrl:textField cmnd/socket_printer/OtaUrl $EVTPART1\
  upgrade:noArg   cmnd/socket_printer/upgrade 1
attr MQTT2_DVES_1791A1 setStateList on off toggle
define FileLog_MQTT2_DVES_1791A1 FileLog ./log/MQTT2_DVES_1791A1-%Y.log MQTT2_DVES_1791A1
setuuid FileLog_MQTT2_DVES_1791A1 5fe0b744-f33f-6f28-1407-20e39305f6f2eae5
attr FileLog_MQTT2_DVES_1791A1 logtype text
attr FileLog_MQTT2_DVES_1791A1 room Steckdose
```

## Sample shelly 1
```
define MQTT2_shelly1_68C63AFA0C2B MQTT2_DEVICE shelly1_68C63AFA0C2B
setuuid MQTT2_shelly1_68C63AFA0C2B 5fe0e334-f33f-6f28-6b42-26e532b2909694db
attr MQTT2_shelly1_68C63AFA0C2B IODev mqttBroker
attr MQTT2_shelly1_68C63AFA0C2B alias Licht-Arbeitszimmer
attr MQTT2_shelly1_68C63AFA0C2B devStateIcon {my $onl = ReadingsVal($name,"online","false") eq "false" ? "rot" : ReadingsVal($name,"new_fw","false") eq "true" ? "gelb" : "gruen";; my $light = Rea>
attr MQTT2_shelly1_68C63AFA0C2B model shelly1
attr MQTT2_shelly1_68C63AFA0C2B readingList shellies/shelly1-68C63AFA0C2B/relay/0:.* state\
  shellies/shelly1-68C63AFA0C2B/relay/0:.* relay0\
  shellies/shelly1-68C63AFA0C2B/input/0:.* input0\
  shellies/shelly1-68C63AFA0C2B/online:.* online\
  shellies/shelly1-68C63AFA0C2B/announce:.* { json2nameValue($EVENT) }\
  shellies/announce:.* { $EVENT =~ m,..id...shelly1-68C63AFA0C2B...mac.*, ? json2nameValue($EVENT) : return }\
shelly1_68C63AFA0C2B:shellies/shelly1-68C63AFA0C2B/info:.* { json2nameValue($EVENT) }\
shelly1_68C63AFA0C2B:shellies/shelly1-68C63AFA0C2B/input_event/0:.* { json2nameValue($EVENT) }
attr MQTT2_shelly1_68C63AFA0C2B room Licht
attr MQTT2_shelly1_68C63AFA0C2B setList off:noArg shellies/shelly1-68C63AFA0C2B/relay/0/command off\
  on:noArg shellies/shelly1-68C63AFA0C2B/relay/0/command on\
  x_update:noArg shellies/shelly1-68C63AFA0C2B/command update_fw\
  x_mqttcom shellies/shelly1-68C63AFA0C2B/command $EVTPART1
define FileLog_MQTT2_shelly1_68C63AFA0C2B FileLog ./log/MQTT2_shelly1_68C63AFA0C2B-%Y.log MQTT2_shelly1_68C63AFA0C2B
setuuid FileLog_MQTT2_shelly1_68C63AFA0C2B 5fe0e334-f33f-6f28-ef3c-bc2785d08ac7ce8d
attr FileLog_MQTT2_shelly1_68C63AFA0C2B logtype text
attr FileLog_MQTT2_shelly1_68C63AFA0C2B room Licht
```

### Sample Door/Window sensor
```
define mqtt_Sensor notify (KontaktHastuer):(open|closed) { \
fhem "set mqttBroker publish fhem/$NAME {\"status\":\"$EVENT\"}";;\
}
setuuid mqtt_Sensor 5fe475a0-f33f-6f28-6f3e-ca9d731cb000f4d1
attr mqtt_Sensor room MQTT,Sensor
```

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

