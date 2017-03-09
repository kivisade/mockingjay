package main

import (
	"fmt"
	"flag"
	"log"
	"path/filepath"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"runtime"
	"github.com/fatih/color"
)

type configT struct {
	BindAddress string `yaml:"BindAddress"`
	ForwardTo   string `yaml:"ForwardTo"`
	NoColors    *bool   `yaml:"NoColors"`
}

var config configT

func (a *configT) mergeConfigs(b *configT, overwrite bool) {
	if overwrite && b.BindAddress != "" || !overwrite && a.BindAddress == "" {
		a.BindAddress = b.BindAddress
	}
	if overwrite && b.ForwardTo != "" || !overwrite && a.ForwardTo == "" {
		a.ForwardTo = b.ForwardTo
	}
	if overwrite && b.NoColors != nil || !overwrite && a.NoColors == nil {
		*a.NoColors = *b.NoColors
	}
}

func (a *configT) dumpConfig() {
	fmt.Println()
	fmt.Printf(`                    BindAddress = "%s"` + "\n", a.BindAddress)
	fmt.Printf(`                    ForwardTo   = "%s"` + "\n", a.ForwardTo)
	fmt.Printf(`                    NoColors    = %v` + "\n\n", *a.NoColors)
}

func parseFlagsAndLoadConfig() {
	configFile := flag.String("config", "", "Configuration file")
	noConfig := flag.Bool("no-config", false, "Disable autoloading default configuration file ('config.yml')")

	bindingAddress := flag.String("listen", ":8080", "Listen address:port (defaults to :8080)")
	forwardTo := flag.String("forward-to", "", "Forward all incoming HTTP requests here (reverse HTTP proxy mode)")
	flagNoColors := flag.Bool("no-colors", false, "Disable color output")

	flag.Parse()

	config = configT{
		BindAddress: *bindingAddress,
		ForwardTo: *forwardTo,
		NoColors: flagNoColors,
	}

	if *noConfig {
		log.Println("Configuration file will not be loaded as per user request.")
		return
	}

	if *configFile == "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		*configFile = dir + string(os.PathSeparator) + "config.yml"
	}

	if data, err := ioutil.ReadFile(*configFile); err != nil {
		log.Println("No configuration file found.")
	} else {
		configFromFile := new(configT)
		err := yaml.Unmarshal(data, configFromFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Loaded configuration file:", *configFile)
		config.mergeConfigs(configFromFile, false)
		config.dumpConfig()
	}

	if *config.NoColors || runtime.GOOS == "windows" {
		color.NoColor = true
	}
}
