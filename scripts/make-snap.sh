#!/bin/bash

d=`date +%d.%m.%Y-%H.%M.%S`
echo "make snap-$d"
/usr/sbin/zfs snapshot Oasis/fs1@snap-$d

lst=`/usr/sbin/zfs list -t snapshot Oasis/fs1 | head -n -7 | tail -n +2 | awk '{print $1}'`
for i in $lst
do
   echo "remove $i"
   /usr/sbin/zfs destroy $i
done

