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

```
git clone https://github.com/dantheman213/watchdog /opt/watchdog
cd /opt/watchdog
make
nohup /opt/watchdog/bin/watchdog >> /var/log/watchdog.log 2>&1 &
tail -f /var/log/watchdog.log
```

### Run At Startup

Add or edit your `/etc/rc.local` and add the `nohup` command above to it. If not running on a Debian-based distro please check your distro's man pages for best practicies. 

### Enable SMTP access to a Gmail Account

Create and use a free Gmail account to send your email alerts from. Learn more:

* https://support.google.com/mail/answer/185833?hl=en

* https://support.google.com/mail/answer/7126229?p=BadCredentials
