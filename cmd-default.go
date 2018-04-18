package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dc0d/errgo"
	"github.com/dc0d/persical"
	"github.com/fatih/color"
)

type cmdDefault struct {
	st     *stat
	out    *bytes.Buffer
	logger *log.Logger
}

func newCmdDefault(st *stat) *cmdDefault {
	result := &cmdDefault{
		st:  st,
		out: &bytes.Buffer{},
	}
	result.logger = log.New(
		result.out,
		"",
		0)
	return result
}

func (c *cmdDefault) start() (state, error) {
	if c.st.input.Today {
		return c.today, nil
	}
	return c.month, nil
}

func (c *cmdDefault) today() (state, error) {
	now := iranNow()
	py, pm, pd := persical.GregorianToPersian(now.Year(), int(now.Month()), now.Day())
	if c.st.piped {
		var d Day
		d.Weekday = now.Weekday()
		d.Persian.Year, d.Persian.Month, d.Persian.Day = py, pm, pd
		d.Gregorian.Year, d.Gregorian.Month, d.Gregorian.Day = now.Year(), int(now.Month()), now.Day()
		js, err := json.Marshal(&d)
		if err != nil {
			return nil, errgo.Mark(err)
		}
		c.out.Write(js)
	} else {
		c.logger.Println(now.Weekday())
		c.logger.Println(py, persical.PersianMonth(pm), pd)
		c.logger.Println(now.Format("2006 Jan 02"))
	}

	return c.show, nil
}

func (c *cmdDefault) month() (state, error) {
	now := iranNow()
	data, _today := monthData(now)
	if c.st.piped {
		js, err := json.Marshal(&data)
		if err != nil {
			return nil, errgo.Mark(err)
		}
		c.out.Write(js)
	} else {
		buf := c.out
		var firstLine bool
		fmt.Fprintln(buf, persical.PersianMonth(_today.Persian.Month), _today.Persian.Year)
		fmt.Fprintln(buf, "---------------------")
		fmt.Fprintln(buf, " Sa Su Mo Tu We Th Fr")
		for _, v := range data.Days {
			wd := int((v.Weekday + 1) % 7)
			if !firstLine {
				firstLine = true
				fmt.Fprint(buf, strings.Repeat("   ", wd))
			}
			var attrs []color.Attribute
			if v.Persian.Day == _today.Persian.Day {
				attrs = append(attrs, color.FgBlack, color.BgWhite)
			}

			fmt.Fprintf(buf, " ")
			if v.Weekday == time.Friday {
				attrs = append(attrs, color.FgBlue)
				color.New(attrs...).Fprintf(buf, "%2d", v.Persian.Day)
			} else {
				color.New(attrs...).Fprintf(buf, "%2d", v.Persian.Day)
			}

			if wd == 6 {
				fmt.Fprint(buf, "\n")
			}
		}

		fmt.Fprint(buf, "\n")
	}
	return c.show, nil
}

func (c *cmdDefault) show() (state, error) {
	fmt.Print(c.out.String())
	return nil, nil
}

func monthData(gdate time.Time) (result Month, today Day) {
	py, pm, pd := persical.GregorianToPersian(gdate.Year(), int(gdate.Month()), gdate.Day())
	today = Day{
		Weekday:   gdate.Weekday(),
		Persian:   Date{Year: py, Month: pm, Day: pd},
		Gregorian: Date{Year: gdate.Year(), Month: int(gdate.Month()), Day: gdate.Day()},
	}
	gy, gm, gd := persical.PersianToGregorian(py, pm, 1)
	gdate = time.Date(gy, time.Month(gm), gd, 0, 0, 0, 0, time.Local)

	currentMonth := pm
	py, pm, pd = persical.GregorianToPersian(gdate.Year(), int(gdate.Month()), gdate.Day())

	for currentMonth == pm {
		result.Days = append(result.Days, Day{
			Weekday:   gdate.Weekday(),
			Persian:   Date{Year: py, Month: pm, Day: pd},
			Gregorian: Date{Year: gdate.Year(), Month: int(gdate.Month()), Day: gdate.Day()},
		})
		gdate = gdate.Add(time.Hour * 24)
		py, pm, pd = persical.GregorianToPersian(gdate.Year(), int(gdate.Month()), gdate.Day())
	}
	return
}
