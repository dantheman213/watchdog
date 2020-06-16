# watchdog

Automatically conduct, monitor, and report on health of disks, ZFS clusters, and UPS hardware to maintain your physical server.

## Getting Started

Download, compile, and run on your server:

```
git clone https://github.com/dantheman213/watchdog /opt/watchdog
cd /opt/watchdog
make
nohup /opt/watchdog/bin/watchdog >> /var/log/watchdog.log 2>&1 &
tail -f /var/log/watchdog.log
```
