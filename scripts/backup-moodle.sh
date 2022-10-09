#!/bin/bash

opt="-alrvP -X -A -M --fake-super --stats --bwlimit=70M"
optd="-alrvP -X -A -M --fake-super --stats --bwlimit=70M --delete"

mkdir -p /mnt/backup/Moodle

/usr/bin/rsync $opt  root@10.10.10.54:/hv4-zfs/moodle* /mnt/backup/Moodle 
#> /dev/null


d=`/usr/bin/date +%d.%m.%Y-%H.%M.%S`
echo "make snap-$d"
#/usr/sbin/zfs snapshot BAG/teh@snap-$d
/usr/bin/mkdir /mnt/backup/Moodle/.snap/snap-$d

lst=`/usr/bin/find /mnt/backup/Moodle/.snap -maxdepth 1 -type d -printf "%Ts\t%P\n" | sort -n | cut -f2 | head -n -3 | tail -n +1`
#lst=`/usr/sbin/zfs list -t snapshot | head -n -12 | tail -n +2 | awk '{print $1}'`

for i in $lst
do
   echo "remove /mnt/backup/Moodle/.snap/$i"
#   /usr/bin/rmdir /mnt/backup/Moodle/.snap/$i
done

