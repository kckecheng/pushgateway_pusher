package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kckecheng/pushgateway_pusher/collector"
	"github.com/kckecheng/pushgateway_pusher/parser"
	"github.com/kckecheng/pushgateway_pusher/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	flag "github.com/spf13/pflag"
)

var debug bool = false

func init() {
	raw, ok := os.LookupEnv("DEBUG")
	if ok {
		switch raw {
		case "TRUE", "True", "true", "1":
			debug = true
		}
	}
}

func getArgs() map[string]string {
	var gateway, application, job, fieldf string
	flag.StringVarP(&gateway, "gateway", "g", "http://localhost:9091", "Pushgateway address")
	flag.StringVarP(&application, "application", "a", "vdbench", "The application generating I/O")
	flag.StringVarP(&job, "job", "j", "pushergateway_pusher_demojob1", "Pushgateway job name")
	flag.StringVarP(&fieldf, "field", "f", "field.yaml", `Definitions of fields and the regular expression pattern to extract fields`)
	flag.Parse()

	ret := map[string]string{
		"gateway":     gateway,
		"application": application,
		"job":         job,
		"fieldf":      fieldf,
	}
	return ret
}

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	params := getArgs()
	p, err := parser.New(params["fieldf"])
	if err != nil {
		panic(err)
	}

	fields := []map[string]string{}
	values := map[string]float64{}
	for _, f := range p.Fields {
		field := map[string]string{}
		field["name"] = f.Name
		field["help"] = f.Description
		fields = append(fields, field)
	}
	pc, err := collector.New(fields, params["application"], values)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	reg := prometheus.NewRegistry()
	reg.MustRegister(pc)
	pusher := push.New(params["gateway"], params["job"]).Gatherer(reg)

	// Delete pushgateway job:
	// - based on signals;
	// - after normal exit
	dfunc := func() {
		if err := pusher.Delete(); err != nil {
			fmt.Fprintf(os.Stderr, "Fail to deleter job %s from %s, please delete it manually", params["job"], params["gateway"])
		}
	}
	go func() {
		<-sigc
		dfunc()
		os.Exit(types.ERREXIT)
	}()
	defer dfunc()

	// Start to get input
	rawInput := make(chan string)
	go parser.Scan(rawInput)
	// Start to calculate error counter after DELAY x seconds
	counter := 0
	timeout := time.After(types.DELAY * time.Second)
	delay := true
	// Process latest data
	for {
		select {
		case line := <-rawInput:
			fmt.Fprintln(os.Stdout, line) // Print out the original line
			rawValues, err := p.Parse(line)
			// Since the regular expression pattern is complicated sometimes, print them for debug purpose
			if debug {
				for i, v := range rawValues {
					fmt.Fprintf(os.Stderr, "%s: %f\n", fields[i]["name"], v)
				}
			}
			// Calculate the num. of continuous errors
			if err != nil {
				if debug {
					fmt.Fprintln(os.Stderr, err)
				}

				// Deplay for specified time before calculating continuous errors
				if delay {
					select {
					case <-timeout:
						delay = false
					default:
						//
					}
				} else {
					counter++
				}

				if counter > types.MAXERROR {
					fmt.Fprintln(os.Stderr, "Too many continuous errors are hit, please press Ctrl + C to kill the program")
					os.Exit(types.ERREXIT)
				} else {
					continue
				}
			} else {
				counter = 0 // reset continuous errors once there is a successful match
			}

			values := map[string]float64{}
			for index, field := range fields {
				name := field["name"]
				value := rawValues[index]
				values[name] = value
			}
			pc.UpdateValues(values)

			if err := pusher.Push(); err != nil {
				panic(err)
			}
		case <-time.After(types.MAXINTERVAL * time.Minute):
			fmt.Fprintf(os.Stderr, "There is no new data within %v minutes, exit", types.MAXINTERVAL)
		}
	}
}
