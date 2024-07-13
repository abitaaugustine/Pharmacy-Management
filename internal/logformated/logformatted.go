package logformated

import (
	"os"
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

const (
	ComponentMain           = "Main"
	ComponentDB             = "DB"
	ComponentAuthentication = "Authentication"
	ComponentUser           = "User"
)

func init() {
	logrus.SetFormatter(&formatter.Formatter{
		HideKeys:        true,
		NoColors:        true,
		TimestampFormat: time.RFC3339,
		FieldsOrder:     []string{"component"},
	})
}

func GetLogger(component string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"component": component,
	})
}

func SetOutput() {
	logrus.SetOutput(os.Stdout)
}
