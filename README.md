# watchdog

Automatically conduct, monitor, and report on health of disks, ZFS clusters, and UPS hardware to maintain your physical server.

## Getting Started

Use the example below in order to download, compile, and run the app on your server.

```
git clone https://github.com/dantheman213/watchdog /tmp/watchdog
cd /tmp/watchdog
make
make install

# edit me before continuing
nano /etc/watchdog/config.json

systemctl enable watchdog
systemctl start watchdog
```

## Features

* Run S.M.A.R.T. `short` test on all your compatible physical disks on your custom schedule.

* Run S.M.A.R.T. `long` test on all your compatible physical disks on your custom schedule.

* Run ZFS scrub on all your pools on your custom schedule.

* Test and review UPS health on your custom schedule.

* Generate reports and have it sent to you via email on your custom schedule.

## Review Logs

```
systemctl status watchdog
tail -f /var/log/watchdog.log
```

## Configuration

### Email Reports

Reports can be sent out via email by modifying `config.json` and modifying the array as you see fit. For day of week, use (0-6) to denote Sunday-Saturday. For time, use military time format (e.g. 23:20 to indicate 11:20 PM or 00:00 to indicate midnight.

#### Enable SMTP access to a Gmail Account

Create and use a free Gmail account to send your email alerts from. Learn more:

* https://support.google.com/mail/answer/185833?hl=en

* https://support.google.com/mail/answer/7126229?p=BadCredentials

### UPS

`watchdog` only supports CyberPower UPS hardware at present. The `pwrstat` utility must be installed and have the ability to detect your UPS.

#### PowerPanel

CyperPower's PowerPanel software also known as `pwrstat` can be installed by [visiting here](https://www.cyberpowersystems.com/product/software/power-panel-personal/powerpanel-for-linux/) and following the directions.
