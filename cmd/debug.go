package http

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	rr "github.com/spiral/roadrunner/cmd/rr/cmd"
	"github.com/spiral/roadrunner/cmd/util"
	"github.com/spiral/jobs"
)

func init() {
	cobra.OnInitialize(func() {
		if rr.Debug {
			svc, _ := rr.Container.Get(jobs.ID)
			if svc, ok := svc.(*jobs.Service); ok {
				svc.AddListener((&debugger{logger: rr.Logger}).listener)
			}
		}
	})
}

// listener provide debug callback for system events. With colors!
type debugger struct{ logger *logrus.Logger }

// listener listens to http events and generates nice looking output.
func (s *debugger) listener(event int, ctx interface{}) {
	if util.LogEvent(s.logger, event, ctx) {
		// handler by default debug package
		return
	}

	switch event {
	case jobs.EventJobAdded:
		e := ctx.(*jobs.JobEvent)
		s.logger.Info(util.Sprintf(
			"jobs.<magenta+h>%s</reset> <white+hb>%s</reset>",
			e.Job.Job,
			e.ID,
		))
	case jobs.EventJobComplete:
		e := ctx.(*jobs.JobEvent)
		s.logger.Info(util.Sprintf(
			"jobs.<green+h>%s</reset> <white+hb>%s</reset>",
			e.Job.Job,
			e.ID,
		))
	case jobs.EventPushError:
		e := ctx.(*jobs.ErrorEvent)
		s.logger.Error(util.Sprintf(
			"jobs.<red>%s</reset> <red+hb>%s</reset>",
			e.Job.Job,
			e.Error.Error(),
		))

	case jobs.EventJobError:
		e := ctx.(*jobs.ErrorEvent)
		s.logger.Error(util.Sprintf(
			"jobs.<red>%s</reset> <white+hb>%s</reset> <red+hb>%s</reset>",
			e.Job.Job,
			e.ID,
			e.Error.Error(),
		))
	}

	//// http events
	//switch event {
	//case rrhttp.EventResponse:
	//	e := ctx.(*rrhttp.ResponseEvent)
	//	s.logger.Info(util.Sprintf(
	//		"<cyan+h>%s</reset> %s <white+hb>%s</reset> %s",
	//		e.Request.RemoteAddr,
	//		statusColor(e.Response.Status),
	//		e.Request.Method,
	//		e.Request.URI,
	//	))
	//case rrhttp.EventError:
	//	e := ctx.(*rrhttp.ErrorEvent)
	//
	//	if _, ok := e.Error.(roadrunner.JobError); ok {
	//		s.logger.Info(util.Sprintf(
	//			"%s <white+hb>%s</reset> %s",
	//			statusColor(500),
	//			e.Request.Method,
	//			uri(e.Request),
	//		))
	//	} else {
	//		s.logger.Info(util.Sprintf(
	//			"%s <white+hb>%s</reset> %s <red>%s</reset>",
	//			statusColor(500),
	//			e.Request.Method,
	//			uri(e.Request),
	//			e.Error,
	//		))
	//	}
	//}
}

//
//func statusColor(status int) string {
//	if status < 300 {
//		return util.Sprintf("<green>%v</reset>", status)
//	}
//
//	if status < 400 {
//		return util.Sprintf("<cyan>%v</reset>", status)
//	}
//
//	if status < 500 {
//		return util.Sprintf("<yellow>%v</reset>", status)
//	}
//
//	return util.Sprintf("<red>%v</reset>", status)
//}
//
//func uri(r *http.Request) string {
//	if r.TLS != nil {
//		return fmt.Sprintf("https://%s%s", r.Host, r.URL.String())
//	}
//
//	return fmt.Sprintf("http://%s%s", r.Host, r.URL.String())
//}
