#!/bin/sh
chmod +x /opt/lookup-portal/lookup-portal
sudo setcap cap_net_bind_service+ep /opt/lookup-portal/lookup-portal
