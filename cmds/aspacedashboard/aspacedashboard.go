package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	description = `
 USAGE: aspacedashboard [OPTIONS]

`

	configuration = `
 CONFIGURATION

 aspacedashboard is configured through environment variables. The following
 variables are supported

   ASPACE_TEMPLATES

   ASPACE_BLEVE_INDEX

   ASPACE_DASHBOARD_URL

`
	help bool
)

func usage() {
	fmt.Println(description)
	flag.PrintDefaults()
	fmt.Println(configuration)
	os.Exit(0)
}

func init() {
	flag.BoolVar(&help, "h", true, "This help message")
	flag.BoolVar(&help, "help", true, "This help message")
}

func main() {
	flag.Parse()
	if help == true {
		usage()
	}

	fmt.Println("aspacedashboard not implemented.")
}
