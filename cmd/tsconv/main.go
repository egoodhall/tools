package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/alecthomas/kong"
)

type Cli struct {
	Timestamp    string `arg:"" name:"timestamp" required:"" help:"The timestamp to convert"`
	OutputFormat string `name:"output-format" short:"o" default:"Mon, 02 Jan 2006 15:04:05.999 -0700" help:""`
}

func main() {
	ctx := kong.Parse(new(Cli))
	ctx.FatalIfErrorf(ctx.Run())
}

func (cli *Cli) Run() error {
	var ts time.Time
	if timestamp, err := strconv.ParseInt(cli.Timestamp, 10, 64); err == nil {
		// Number of digits in current date's timestamp
		currentSecondsSize := len(strconv.Itoa(int(time.Now().Unix())))

		if len(cli.Timestamp) <= currentSecondsSize {
			ts = time.Unix(timestamp, 0)
		} else if len(cli.Timestamp) <= currentSecondsSize+3 {
			ts = time.UnixMilli(int64(timestamp))
		} else {
			ts = time.UnixMicro(int64(timestamp))
		}
	} else {
		for _, format := range []string{time.RFC3339, time.RFC3339Nano, time.RFC1123, time.RFC1123Z} {
			ts, err = time.Parse(format, cli.Timestamp)
			if err == nil {
				break
			}
		}
	}

	if ts.IsZero() {
		return errors.New("unable to parse timestamp")
	}

	fmt.Println(ts.Format(cli.OutputFormat))
	return nil
}
