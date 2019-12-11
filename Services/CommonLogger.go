package Services

import (
	kitlog "github.com/go-kit/kit/log"
	"os"
)

var logger kitlog.Logger

func init() {

	{
		logger = kitlog.NewLogfmtLogger(os.Stdout)
		logger = kitlog.WithPrefix(logger, "SpaceManagementSystem", "1.0")
		logger = kitlog.With(logger, "time", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}
}
