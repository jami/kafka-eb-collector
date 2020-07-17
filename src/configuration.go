package src

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	envPrefix            = "COLLECTOR"
	defaultEventBusTopic = "eventbus-collector"
)

// CLIConfiguration model
type CLIConfiguration struct {
	Broker        string
	ShowHelp      bool
	EventBusTopic string
}

func registerVar(name string, target interface{}, desc string) {
	envName := strings.ToUpper(envPrefix + "_" + strings.Replace(name, "-", "_", -1))
	envValue := os.Getenv(envName)

	if name != "help" {
		desc = fmt.Sprintf("%s - %s", desc, envName)
	}

	switch t := target.(type) {
	case *string:
		flag.StringVar(t, name, envValue, desc)
	case *bool:
		flag.BoolVar(t, name, false, desc)
	case *int:
		i, err := strconv.Atoi(envValue)
		if err != nil {
			i = 0
		}

		flag.IntVar(t, name, i, desc)
	default:
		fmt.Printf("unkniwn argument type %T\n", t)
	}
}

// Parse program arguments
func (ec *CLIConfiguration) Parse() error {
	registerVar("help", &ec.ShowHelp, "show help dialog")
	registerVar("event-topic", &ec.EventBusTopic, "kafka topic for events")
	registerVar("brokerlist", &ec.Broker, "kafka broker list")

	flag.Parse()

	if ec.EventBusTopic == "" {
		ec.EventBusTopic = defaultEventBusTopic
	}

	if ec.Broker == "" {
		return fmt.Errorf("Missing argument - brokerlist")
	}

	return nil
}
