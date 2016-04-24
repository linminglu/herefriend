package lib

import (
	log "github.com/cihub/seelog"
)

const gLogConfig = `
<seelog>
	<outputs formatid="common">
		<rollingfile type="size" filename="./log/roll.log" maxsize="1000000" maxrolls="5000000" />
	</outputs>
	<formats>
		<format id="common" format="%Date/%Time [%LEV] %Msg%n"/>
	</formats>
</seelog>
`

func init() {
	logger, err := log.LoggerFromConfigAsBytes([]byte(gLogConfig))
	if nil != err {
		panic(err.Error())
	}

	log.ReplaceLogger(logger)
}
