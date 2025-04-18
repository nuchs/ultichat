package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
)

func main() {
	signalHandler()

	scanner := bufio.NewScanner(os.Stdin)
	if err := run(scanner, os.Stdout, os.Stderr, os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func signalHandler() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, os.Kill)
}

func run(scanner *bufio.Scanner, stdout, stderr io.Writer, args []string) error {
	service := readOptions(stderr, args)
	fmt.Println("Filter:", service)

	for scanner.Scan() {
		m := make(map[string]any)
		err := json.Unmarshal(scanner.Bytes(), &m)
		if err != nil {
			if service == "" {
				fmt.Fprintln(stdout, scanner.Text())
			}
			continue
		}

		if !filter(service, m) {
			continue
		}

		printStructuredLog(stdout, m)
	}

	return scanner.Err()
}

func readOptions(stderr io.Writer, args []string) string {
	var service string

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.StringVar(&service, "service", "", "filter on service name")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "usage: %s [options]\n", args[0])
		fs.PrintDefaults()
		fmt.Fprintf(stderr, "\nThis program accepts json logs piped to it and pretty prints them to the stdout. It can filter on service attribute of structured logs.\n")
	}

	fs.Parse(args[1:])

	return strings.ToLower(service)
}

func filter(service string, m map[string]any) bool {
	if service == "" {
		return true
	}
	serviceStr, ok := m["service"].(string)
	fmt.Printf("filter: %s, service: %s, ok: %t\n", service, serviceStr, ok)
	return ok && strings.ToLower(serviceStr) == service
}

func printStructuredLog(stdout io.Writer, m map[string]any) {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%s | %s | %s | %s\n", m["service"], m["time"], m["level"], m["msg"]))
	for k, v := range m {
		switch k {
		case "service", "time", "level", "msg":
			continue
		}
		b.WriteString(fmt.Sprintf("\t%s = %v\n", k, v))
	}

	fmt.Fprint(stdout, b.String())
}
