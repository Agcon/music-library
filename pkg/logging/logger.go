package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

func Init() {
	Log = logrus.New()
	Log.Out = os.Stdout
	Log.SetLevel(logrus.DebugLevel)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
