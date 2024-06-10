package backend

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	emailsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "email_processed_total",
		Help: "The total number of processed emails",
	})

	emailsFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "email_failed_total",
		Help: "The total number of emails that has failed",
	})

	emailsSucceeded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "email_succeeded_total",
		Help: "The total number of emails that has been successfully sent",
	})

	authenticationFailures = promauto.NewCounter(prometheus.CounterOpts{
		Name: "authentication_failures_total",
		Help: "The total number of SMTP authentication failures",
	})

	authenticationSuccesses = promauto.NewCounter(prometheus.CounterOpts{
		Name: "authentication_ok_total",
		Help: "The total number of successful SMTP authentications",
	})

	emailParseFailures = promauto.NewCounter(prometheus.CounterOpts{
		Name: "email_parse_failures_total",
		Help: "The total number of unparseable emails",
	})
)
