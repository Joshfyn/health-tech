package clog

import (
	"net/http"
	"os"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

const (
	HTTPRequestKey = "http_request"
	ServerNameKey  = "server_name"
)

type Clog struct {
	ec     *logrus.Entry
	logrus *logrus.Entry
	raven  *raven.Client
}

func (c *Clog) WithField(message string, args any) *Clog {
	c.ec = c.ec.WithField(message, args)
	return c
}

func (c *Clog) WithFields(fields logrus.Fields) *Clog {
	c.ec = c.ec.WithFields(fields)
	return c
}

func (c *Clog) Errorf(r *http.Request, format string, args ...any) {
	req := c.ec
	req = req.WithField("resource", r)
	req.Errorf(format, args...)
}

// MustNew creates a new logger context, which will panic on error.
func MustNew(serviceName, dsn string) *Clog {
	l := logrus.New()
	l.Out = os.Stdout

	// Set formatter because default behavior not escape newline characters
	// Refer: https://github.com/sirupsen/logrus/issues/608
	l.SetFormatter(&logrus.TextFormatter{
		DisableQuote:     true,
		DisableSorting:   true,
		QuoteEmptyFields: true,
	})

	clog := &Clog{logrus: logrus.NewEntry(l)}

	if serviceName != "" {
		clog.logrus = l.WithField(ServerNameKey, serviceName)
	}

	if dsn != "" {
		client, err := raven.New(dsn)
		if err != nil {
			panic(err)
		}
		hook, err := logrus_sentry.NewAsyncWithClientSentryHook(client, logrus.AllLevels)
		if err != nil {
			panic(err)
		}
		hook.StacktraceConfiguration.Enable = true
		l.Hooks.Add(hook)
		clog.raven = client
	}
	clog.ec = clog.logrus.Dup()

	return clog
}
