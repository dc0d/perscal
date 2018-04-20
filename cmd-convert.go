package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dc0d/errgo"
	"github.com/dc0d/persical"
)

type cmdConvert struct {
	st     *stat
	out    *bytes.Buffer
	logger *log.Logger
	day    Day
}

func newCmdConvert(st *stat) *cmdConvert {
	result := &cmdConvert{
		st:  st,
		out: &bytes.Buffer{},
	}
	result.logger = log.New(
		result.out,
		"",
		0)
	return result
}

func (c *cmdConvert) start() (state, error) {
	if !c.st.input.convert.p2g &&
		!c.st.input.convert.g2p {
		errlog.Fatalln("conversion direction not specified")
		return nil, nil
	}
	if c.st.input.convert.Year <= 0 {
		errlog.Fatalln("year can not be negative")
		return nil, nil
	}
	if c.st.input.convert.Month <= 0 {
		errlog.Fatalln("month can not be negative")
		return nil, nil
	}
	if c.st.input.convert.Day <= 0 {
		errlog.Fatalln("day can not be negative")
		return nil, nil
	}
	return c.convert, nil
}

func (c *cmdConvert) convert() (state, error) {
	if c.st.input.convert.g2p {
		gdate := time.Date(c.st.input.convert.Year, time.Month(c.st.input.convert.Month), c.st.input.convert.Day, 0, 0, 0, 0, time.Local)
		py, pm, pd := persical.GregorianToPersian(gdate.Year(), int(gdate.Month()), gdate.Day())
		var d Day
		d.Weekday = gdate.Weekday()
		d.Persian.Year, d.Persian.Month, d.Persian.Day = py, pm, pd
		d.Gregorian.Year, d.Gregorian.Month, d.Gregorian.Day = gdate.Year(), int(gdate.Month()), gdate.Day()
		c.day = d
	}
	if c.st.input.convert.p2g {
		py, pm, pd := c.st.input.convert.Year, c.st.input.convert.Month, c.st.input.convert.Day
		gy, gm, gd := persical.PersianToGregorian(py, pm, pd)
		gdate := time.Date(gy, time.Month(gm), gd, 0, 0, 0, 0, time.Local)
		var d Day
		d.Weekday = gdate.Weekday()
		d.Persian.Year, d.Persian.Month, d.Persian.Day = py, pm, pd
		d.Gregorian.Year, d.Gregorian.Month, d.Gregorian.Day = gdate.Year(), int(gdate.Month()), gdate.Day()
		c.day = d
	}

	if c.st.piped {
		return c.json, nil
	}
	return c.terminal, nil
}

func (c *cmdConvert) json() (state, error) {
	js, err := json.Marshal(&c.day)
	if err != nil {
		return nil, errgo.Mark(err)
	}
	c.out.Write(js)
	return c.show, nil
}

func (c *cmdConvert) terminal() (state, error) {
	c.out.WriteString(fmt.Sprintf("%v\n", c.day.Weekday))
	c.out.WriteString(fmt.Sprintf("%v %v %v\n",
		c.day.Persian.Year,
		persical.PersianMonth(c.day.Persian.Month),
		c.day.Persian.Day))
	c.out.WriteString(fmt.Sprintf("%v %v %v\n",
		c.day.Gregorian.Year,
		time.Month(c.day.Gregorian.Month),
		c.day.Gregorian.Day))

	return c.show, nil
}

func (c *cmdConvert) show() (state, error) {
	fmt.Print(c.out.String())
	return nil, nil
}
