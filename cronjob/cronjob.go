package cronjob

import (
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

func StartCronJob(cronString string, job cron.Job, timezone *time.Location) {
	c := cron.New(
		cron.WithLocation(timezone),
		cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron:", log.LstdFlags))),
	)

	// Allow recover from panic job
	c.AddJob(cronString, cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(job))

	c.Start()
}
