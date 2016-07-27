#!/bin/sh
chmod +x /opt/arl/arl
sudo setcap cap_net_bind_service+ep /opt/arl/arl
