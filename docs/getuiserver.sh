#!/bin/bash
#
# getuiserver        Startup script for the Apache HTTP Server
#
# chkconfig: - 85 15
# description: The Apache HTTP Server is an efficient and extensible  \
#	       server implementing the current HTTP standards.
# processname: getuiserver
# pidfile: /var/run/getuiserver/getuiserver.pid
#
### BEGIN INIT INFO
# Provides: getuiserver
# Description: The Apache HTTP Server is an extensible server 
#  implementing the current HTTP standards.
### END INIT INFO

# Source function library.
. /etc/rc.d/init.d/functions

basedir=/root/workspace/getui
getuiserver=/root/workspace/getui/getuiserver
prog=getuiserver
pidfile=${PIDFILE-/root/workspace/getui/getuiserver.pid}
lockfile=${LOCKFILE-/var/lock/subsys/getuiserver}
RETVAL=0
STOP_TIMEOUT=${STOP_TIMEOUT-10}

# The semantics of these two functions differ from the way apachectl does
# things -- attempting to start while running is a failure, and shutdown
# when not running is also a failure.  So we just do it the way init scripts
# are expected to behave here.
start() {
        echo -n $"Starting $prog: "
        pushd ${basedir} > /dev/null
        daemon --pidfile=${pidfile} "LD_LIBRARY_PATH=/usr/local/lib $getuiserver &"
        RETVAL=$?
        popd > /dev/null
        
        echo
        [ $RETVAL = 0 ] && touch ${lockfile}
        return $RETVAL
}

# When stopping getuiserver, a delay (of default 10 second) is required
# before SIGKILLing the getuiserver parent; this gives enough time for the
# getuiserver parent to SIGKILL any errant children.
stop() {
	echo -n $"Stopping $prog: "
	killproc -p ${pidfile} -d ${STOP_TIMEOUT} $getuiserver
	RETVAL=$?
	echo
	[ $RETVAL = 0 ] && rm -f ${lockfile} ${pidfile}
}

# See how we were called.
case "$1" in
  start)
	start
	;;
  stop)
	stop
	;;
  status)
        status -p ${pidfile} $getuiserver
	RETVAL=$?
	;;
  restart)
	stop
	start
	;;
  *)
	echo $"Usage: $prog {start|stop|restart|status}"
	RETVAL=2
esac

exit $RETVAL
