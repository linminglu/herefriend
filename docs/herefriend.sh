#!/bin/bash
#
# herefriend        Startup script for the Apache HTTP Server
#
# chkconfig: - 85 15
# description: The Apache HTTP Server is an efficient and extensible  \
#	       server implementing the current HTTP standards.
# processname: herefriend
# pidfile: /var/run/herefriend/herefriend.pid
#
### BEGIN INIT INFO
# Provides: herefriend
# Description: The Apache HTTP Server is an extensible server 
#  implementing the current HTTP standards.
### END INIT INFO

# Source function library.
. /etc/rc.d/init.d/functions

basedir=/root/workspace/hf
herefriend=/root/workspace/hf/herefriend
prog=herefriend
pidfile=${PIDFILE-/root/workspace/hf/herefriend.pid}
lockfile=${LOCKFILE-/var/lock/subsys/herefriend}
RETVAL=0
STOP_TIMEOUT=${STOP_TIMEOUT-10}

# The semantics of these two functions differ from the way apachectl does
# things -- attempting to start while running is a failure, and shutdown
# when not running is also a failure.  So we just do it the way init scripts
# are expected to behave here.
start() {
        echo -n $"Starting $prog: "
        pushd ${basedir} > /dev/null
        daemon --pidfile=${pidfile} "$herefriend > /dev/null 2>&1 &"
        RETVAL=$?
        popd > /dev/null
        
        echo
        [ $RETVAL = 0 ] && touch ${lockfile}
        return $RETVAL
}

# When stopping herefriend, a delay (of default 10 second) is required
# before SIGKILLing the herefriend parent; this gives enough time for the
# herefriend parent to SIGKILL any errant children.
stop() {
	echo -n $"Stopping $prog: "
	killproc -p ${pidfile} -d ${STOP_TIMEOUT} $herefriend
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
        status -p ${pidfile} $herefriend
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
