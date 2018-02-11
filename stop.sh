#!/bin/sh

sevlist=('dbmgr' 'example' 'gateway' 'logic' 'room' 'upload')

for sevr in ${sevlist[@]};
do
	sleep 1
	ID=`ps -elf | grep "$sevr" | grep -v "grep" | awk '{print $2}'` 
	if [ -n "$ID" ]; then 
		echo "$ID->$sevr stopping..."
		kill -9 $ID
	fi
done
