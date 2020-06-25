# watchdog

Automatically conduct, monitor, and report on health of disks, ZFS clusters, and UPS hardware to maintain your physical server.

## Features

* Run S.M.A.R.T. `short` test on all your compatible physical disks every day.

* Run S.M.A.R.T. `long` test on all your compatible physical disks every Sunday.

* Run ZFS scrub on all your pools every Wednesday.

* [IN PROGRESS] Check and report on UPS activations or failures.

* Generate reports and have it sent to you via email every Saturday.

## Getting Started

Use the example below in order to download, compile, and run the app on your server.

### Download and Install

```
git clone https://github.com/dantheman213/watchdog /tmp/watchdog
cd /tmp/watchdog
make install

# edit me before continuing
nano /etc/watchdog/config.json

systemctl enable watchdog
systemctl start watchdog
```

### Review Logs

```
systemctl status watchdog
tail -f /var/log/watchdog.log
```

## Configuration

### Email Reports

Reports can be sent out via email by modifying `config.json` and modifying the array as you see fit. For day of week, use (0-6) to denote Sunday-Saturday. For time, use military time format (e.g. 23:20 to indicate 11:20 PM or 00:00 to indicate midnight.

## Additional Info

### Enable SMTP access to a Gmail Account

Create and use a free Gmail account to send your email alerts from. Learn more:

* https://support.google.com/mail/answer/185833?hl=en

* https://support.google.com/mail/answer/7126229?p=BadCredentials
