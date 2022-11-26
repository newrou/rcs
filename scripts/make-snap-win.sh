#!/bin/bash

pool=$1
fs=$2
nsnap=$3
user=$4
pswd=$5

# mount win
umount /mnt/$fs
mkdir -p /mnt/$fs
mount -t cifs -o username=$user,password=$pswd //192.168.0.10/Share /mnt/$fs

# rsync res
opt="-alrvP -X -A -M --fake-super --stats --bwlimit=70M"
/usr/bin/rsync $opt /mnt/$fs/* /$pool/$fs
umount /mnt/win

# make new snap
d=`date +%d.%m.%Y-%H.%M.%S`
echo "make snap-$d"
/usr/sbin/zfs create $pool/$fs
/usr/sbin/zfs snapshot $pool/$fs@snap-$d

#remove old snaps
lst=`/usr/sbin/zfs list -t snapshot $pool/$fs | head -n -$nsnap | tail -n +2 | awk '{print $1}'`
for i in $lst
do
   echo "remove $i"
   /usr/sbin/zfs destroy $i
done

ln -s /$pool/$fs/.zfs/snapshot /var/lib/rcs/snapshot/$fs
