package cmd

import (
	"errors"
	"flag"
	"log/slog"

	"github.com/AugustineAurelius/yuki/cmd/serve"
	"github.com/AugustineAurelius/yuki/cmd/yuki_test"
)

var (
	command string
	host    string
	port    int

	ErrUndefinedCommand = errors.New("undefined command")
)

func Execute() {
	flag.StringVar(&command, "command", "serve", "Call command to yuki (serve, test)")
	flag.StringVar(&host, "host", "0.0.0.0", "host to run yuki")
	flag.IntVar(&port, "port", 5555, "port to run yuki")

	flag.Parse()

	slog.Info("flags", "command", command, "host", host, "port", port)

	switch command {
	case "serve":
		serve.YukiServe(host, port)
	case "test":
		if err := yuki_test.YukiTest(host, port); err != nil {
			slog.Error(err.Error())
		}
	case "inf_test":
		if err := yuki_test.YukiInfTest(host, port); err != nil {
			slog.Error(err.Error())
		}
	}

}
