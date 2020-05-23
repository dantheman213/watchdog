#!/usr/bin/env bash

apt-get update
apt-get install -y nut smartmontools

echo "NOTE: if ZFS is not installed or configured then crontab should be modified after install"

sudo cp -Rav rootfs/* /\

cat << EOF >> /var/spool/cron/crontabs/$(whoami)
# Run SMART - short disk check every 12 hours
0 */12 * * * bash /usr/bin/check-disks.sh short > /var/log/check-disks.log 2>&1

# Run SMART - long disk once a week on Sunday
0 0 * * 6 bash /usr/bin/check-disks.sh long > /var/log/check-disks.log 2>&1

# Run a complete data check on the pool approx every 2 weeks
0 0 */15 * * /usr/sbin/zpool scrub > /var/log/zpool-scrub.log 2>&1
EOF

# TODO: check zfs pool status

# TODO: send issues via email, slack, sms, etc
