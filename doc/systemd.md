## Systemd

### create this file under /lib/systemd/system/hanectl.service

```
[Unit]
Description=Hane-Ctl
After=network.target

[Service]
Type=simple
User=polaris 
Group=polaris 

#RuntimeDirectory=hanectl
#LogsDirectory=hanectl
#PIDFile=/run/hanectl/hanectl.pid
ExecStart=/opt/hanectl/hanectl -config /opt/hanectl/config/config.yml 

[Install]
WantedBy=multi-user.target

```

### activate service
> sudo systemctl daemon-reload

> sudo systemctl enable hanectl

> sudo systemctl start hanectl
