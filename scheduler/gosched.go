package main

import (
	"context"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/reugn/go-quartz/job"
	"github.com/reugn/go-quartz/quartz"
)

func playTime() {
	_l := log.New(log.Default().Writer(), "[ticker] - ", log.Default().Flags())

	_l.Print("> playTime")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	go func() {
		_l.Print("> gorutine")
		for {
			select {
			case <-ctx.Done():
				_l.Print("< gorutine")
				return // returning not to leak the goroutine
			case t := <-t.C:
				_l.Printf("Tik - %+v", t)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	_l.Print(">>>>>")
	time.Sleep(time.Second * 5)
	_l.Print("<<<<<")

	cancel()
	time.Sleep(time.Second)

	_l.Print("< playTime")
}

func playQuartz() {
	_l := log.New(log.Default().Writer(), "[quartz] - ", log.Default().Flags())

	_l.Print("> playQuartz")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sched := quartz.NewStdScheduler()
	sched.Start(ctx)

	triger, err := quartz.NewCronTrigger("* * * * * *")
	if err != nil {
		_l.Panic(err)
	}
	fjob := job.NewFunctionJobWithDesc("Check appiations", func(_ context.Context) (int, error) {
		_l.Print("Tik")
		// panic(errors.New("die!"))
		return 1, nil
	})
	// err = sched.ScheduleJob(fjob, triger)
	djob := quartz.NewJobDetail(fjob, quartz.NewJobKey("functionJob"))
	err = sched.ScheduleJob(djob, triger)
	if err != nil {
		_l.Panic(err)
	}
	_l.Printf("scheduled job: %v - %v", djob.JobKey(), fjob.Description())

	_l.Print(">>>>>")
	time.Sleep(time.Second * 5)
	_l.Print("<<<<<")

	_l.Printf("cancel job: %v", djob.JobKey())
	sched.DeleteJob(djob.JobKey())

	sched.Stop()
	sched.Wait(ctx)
	_l.Print("< playQuartz")
}

func playGocron() {
	_l := log.New(log.Default().Writer(), "[gocron] - ", log.Default().Flags())

	_l.Print("> playGocron")
	s, err := gocron.NewScheduler(
		gocron.WithLocation(time.UTC),
		gocron.WithLogger(gocron.NewLogger(gocron.LogLevelDebug)),
		gocron.WithGlobalJobOptions(
			gocron.WithEventListeners(
				gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
					_l.Printf("> %v, %v", jobID, jobName)
				}),
				gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
					_l.Printf("< %v, %v", jobID, jobName)
				}),
				gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
					_l.Printf("< %v, %v, %v", jobID, jobName, err)
				}),
			),
		),
	)
	if err != nil {
		_l.Panic(err)
	}
	s.Start()
	defer s.Shutdown()

	j, err := s.NewJob(
		gocron.CronJob("* * * * * *", true),
		gocron.NewTask(
			func() {
				_l.Print("Tik")
			},
		),
		gocron.WithName("olala"),
		gocron.WithTags("game", "llm"),
	)
	if err != nil {
		_l.Panic(err)
	}
	_l.Printf("Job: %v", j.ID())

	_l.Print(">>>>>")
	time.Sleep(time.Second * 5)
	_l.Print("<<<<<")

	// s.Shutdown()

	_l.Print("< playGocron")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	log.Print("------------")
	playTime()

	log.Print("------------")
	playQuartz()

	log.Print("------------")
	playGocron()
}
