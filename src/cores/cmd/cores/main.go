package main

import (
	"context"
	"flag"
	"os"
	"time"

	"cos-backend-com/src/cores/cmd/cores/app"
	"cos-backend-com/src/libs/models/bountymodels"

	"github.com/qiniu/x/log"
	s "github.com/wujiu2020/strip"
	"github.com/wujiu2020/strip/utils/helpers"
)

func main() {
	strip := s.New()

	var (
		confPath string
		files    helpers.SliceValue
	)
	flagSet := flag.NewFlagSet("app", flag.ExitOnError)
	flagSet.StringVar(&confPath, "conf", "", "")
	flagSet.Var(&files, "f", "")
	flagSet.Parse(os.Args[1:])

	app.AppInit(strip, confPath, files...).Start()
	go closeBounty()
	strip.Run()
}

func closeBounty() {
	t := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-t.C:
			if err := bountymodels.Bounties.BatchClosedBounty(context.Background()); err != nil {
				log.Warn(err)
			}
		}
	}
}
