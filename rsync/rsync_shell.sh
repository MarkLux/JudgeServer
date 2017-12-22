#!/bin/bash
while true
do

        rsync -av --password-file=/home/judge/testcases/rsyncd.passwd --include="*.in" --include="*.out"  --include="*/"  --exclude="*" --delete  user@host.cn::testcases /home/judge/testcases
        sleep 300
done
