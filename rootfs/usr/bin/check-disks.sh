#!/usr/bin/env bash

# Finds all physical disks on the machine and runs a short or long disk check using smartctl.

if [ -z "$1" ]; then
    echo "Error: one argument needed."
    echo "Usage: ./check-disks.sh <short|long>"
    exit 1
fi

# skip if on days 1-2, 28-31, 16-17 to avoid running same time as ZFS scrub (runs on 1st and 16th every month)
DAY=$(date +"%d")
if [[ "$1" -eq "long" ]]; then
  if [[ $DAY -lt 3 ]] || [[ $DAY -gt 27 ]] || [[ $DAY -gt 15 && $DAY -lt 18 ]]; then
    echo "ZFS scrub is scheduled soon... skipping this iteration!"
    exit 0
  fi
fi

DISK_LIST=$(/sbin/fdisk -l | grep "Disk /dev" | grep -v "/dev/loop" | awk '{print $2}' | sed 's/.$//')

for disk in $DISK_LIST; do
  echo "Running $1 disk test on $disk"
  /sbin/smartctl -t "$1" "$disk"
done
