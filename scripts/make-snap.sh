#!/bin/bash

fs=$1
nsnap=$2
pool='Snuff'

d=`date +%d.%m.%Y-%H.%M.%S`
echo "make snap-$d"
/usr/sbin/zfs snapshot $pool/$fs@snap-$d

lst=`/usr/sbin/zfs list -t snapshot $pool/$fs | head -n -$nsnap | tail -n +2 | awk '{print $1}'`
for i in $lst
do
   echo "remove $i"
   /usr/sbin/zfs destroy $i
done

