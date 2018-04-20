package main

import (
	"flag"
	"os"

	"github.com/dc0d/clarg"
)

func main() {
	st := &stat{}

	fi, _ := os.Stdout.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		st.piped = true
	} else {
		st.piped = false
	}

	top, sub := initCommands(st)
	cmd, err := clarg.Parse(top, sub...)
	if err != nil {
		errlog.Fatalln(err)
		return
	}

	if st.input.json {
		st.piped = true
	}

	if err := dispatch(cmd, st); err != nil {
		errlog.Println(err)
	}
}

func dispatch(cmd *flag.FlagSet, st *stat) error {
	switch cmd.Name() {
	case "ᐸTᐳ":
		errlog.Println("no generics for you.")
	case "conv":
		return activate(newCmdConvert(st).start)
	default:
		return activate(newCmdDefault(st).start)
	}
	return nil
}

func initCommands(st *stat) (top *flag.FlagSet, sub []*flag.FlagSet) {
	var (
		topFlags     = flag.NewFlagSet("", flag.ExitOnError)
		convertFlags = flag.NewFlagSet("conv", flag.ExitOnError)
	)

	topFlags.BoolVar(&st.input.today, "t", st.input.today, "-t show info about today")
	topFlags.BoolVar(&st.input.json, "j", st.input.json, "-j json output")
	topFlags.IntVar(&st.input.month, "m", st.input.month, "-m persian month")
	topFlags.IntVar(&st.input.year, "y", st.input.year, "-m persian year")
	top = topFlags

	convertFlags.IntVar(&st.input.convert.Year, "y", st.input.convert.Year, "-y year")
	convertFlags.IntVar(&st.input.convert.Month, "m", st.input.convert.Month, "-m month")
	convertFlags.IntVar(&st.input.convert.Day, "d", st.input.convert.Day, "-d day")
	convertFlags.BoolVar(&st.input.convert.p2g, "p2g", st.input.convert.p2g, "-p2g persian to gregorian")
	convertFlags.BoolVar(&st.input.convert.g2p, "g2p", st.input.convert.g2p, "-g2p gregorian to persian")
	sub = append(sub, convertFlags)

	return
}
