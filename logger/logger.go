package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

/*
 Logger levels -
 Panic (0)  → the most severe, app usually crashes
 Fatal (1)  → fatal error, app will exit
 Error (2)  → error but app can continue
 Warn  (3)  → something unexpected, but not breaking
 Info  (4)  → general information
 Debug (5)  → verbose debugging details
 Trace (6)  → the most detailed, extremely verbose
*/

var Log *logrus.Logger

func Init() {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetLevel(logrus.InfoLevel)
}
