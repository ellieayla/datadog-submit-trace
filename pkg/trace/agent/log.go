package submittrace

import (
	"fmt"
	coreconfig "github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/trace/config"
	"github.com/DataDog/datadog-agent/pkg/util/log"

	"github.com/cihub/seelog"
)

const loggerName coreconfig.LoggerName = "TRACE"

const loggerConfig = `
<seelog minlevel="%[1]s">
  <outputs formatid="%[2]s">
    <filter levels="trace,debug,info,critical">
      <console />
    </filter>
  </outputs>
  <formats>
    <format id="json" format="%[3]s"/>
    <format id="common" format="%[4]s"/>
  </formats>
</seelog>
`

func setupLogger(cfg *config.AgentConfig) error {

	logLevel := "trace"
	minLogLvl, ok := seelog.LogLevelFromString(logLevel)
	if !ok {
		minLogLvl = seelog.InfoLvl
	}

	format := "common"
	if coreconfig.Datadog.GetBool("log_format_json") {
		format = "json"
	}

	logConfig := fmt.Sprintf(
		loggerConfig,
		minLogLvl,
		format,
		coreconfig.BuildJSONFormat(loggerName),
		coreconfig.BuildCommonFormat(loggerName),
	)
	logger, err := seelog.LoggerFromConfigAsString(logConfig)
	if err != nil {
		return err
	}

	seelog.ReplaceLogger(logger)

	log.SetupDatadogLogger(logger, minLogLvl.String())

	return nil
}
