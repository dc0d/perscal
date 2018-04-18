package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

type stat struct {
	piped bool

	input struct {
		// top level
		JSON  bool
		Today bool

		convert struct {
			Date
			P2G bool
			G2P bool
		}
	}
}

// Month .
type Month struct {
	Days []Day `json:"days,omitempty"`
}

// Day .
type Day struct {
	Weekday   time.Weekday `json:"weekday,omitempty"`
	Persian   Date         `json:"persian,omitempty"`
	Gregorian Date         `json:"gregorian,omitempty"`
}

// Date .
type Date struct {
	Year  int `json:"year,omitempty"`
	Month int `json:"month,omitempty"`
	Day   int `json:"day,omitempty"`
}

var (
	errlog *log.Logger
)

func init() {
	appName := filepath.Base(os.Args[0])
	errlog = log.New(
		os.Stderr,
		appName+color.New(color.FgHiRed).Sprintf(" [ error ] "),
		log.Ltime|log.Lshortfile)
}

func iranTime(source time.Time) time.Time {
	var dest time.Time
	loc, err := time.LoadLocation("Asia/Tehran")
	if err == nil {
		dest = source.In(loc)
	} else {
		dest = source
	}
	return dest
}

func iranNow() time.Time {
	return iranTime(time.Now())
}
