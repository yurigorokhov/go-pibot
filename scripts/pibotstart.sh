#!/bin/bash

cd /home/pi

echo "[INFO] Updating BatBot"
rm -f /home/pi/package.tar.gz

# wait for wifi to be available for a little while
sleep 15

wget https://s3-us-west-2.amazonaws.com/batbot/package/package.tar.gz && tar zxvf /home/pi/package.tar.gz && cp /home/pi/scripts/pibitservice.sh /etc/init.d/ && cp /home/pi/scripts/pibotstart.sh /usr/local/bin/

# star the process
/home/pi/go-pibot >> /home/pi/pibot.log &
