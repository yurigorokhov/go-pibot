#! /bin/sh
### BEGIN INIT INFO
# Provides:          pibotservice
# Required-Start:    $all
# Required-Stop:
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: PiBot!
### END INIT INFO

PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/opt/bin

. /lib/init/vars.sh
. /lib/lsb/init-functions
# If you need to source some other scripts, do it here

case "$1" in
    start)
	log_begin_msg "Starting pibot"
	killall go-pibot
	pkill uv4l
	sleep 2
	uv4l -nopreview --auto-video_nr --driver raspicam --encoding mjpeg --width 320 --height 240 --framerate 10 --server-option '--port=9090' --server-option '--max-queued-connections=30' --server-option '--max-streams=25' --server-option '--max-threads=29'
	echo "starting go-pibot"
	/usr/local/bin/pibotstart.sh

	log_end_msg $?
	exit 0
	;;
    stop)
	log_begin_msg "Stopping the coolest service ever unfortunately"
	pkill uv4l
	killall go-pibot

	log_end_msg $?
	exit 0
	;;
    *)
	echo "Usage: /etc/init.d/<your script> {start|stop}"
	exit 1
	;;
esac
