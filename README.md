# watchdog

Automatically monitor and report on health of disks, ZFS clusters, and UPS hardware.

## Getting Started

Run locally in the background

```
nohup /opt/watchdog/bin/watchdog >> /var/log/watchdog.log 2>&1 &
```

### Enable SMTP access to a Gmail Account

Create and use a free Gmail account to send your email alerts from. Learn more:

https://support.google.com/mail/answer/7126229?p=BadCredentials
