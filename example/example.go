package main

import (
	"github.com/droundy/goopt"
	"github.com/jamesandariese/naglevelparse"
	"log"
	"os"
	"strconv"
)

var warning = goopt.String([]string{"--warning", "-w"}, ":9", "Warning level")
var critical = goopt.String([]string{"--critical", "-c"}, ":10", "Critical level")

func main() {
	goopt.Version = "1.0"
	goopt.Summary = "Sample of nagios level parsing"
	goopt.Author = "James Andariese"
	goopt.Parse(nil)

	if len(goopt.Args) > 1 {
		log.Println("UNKNOWN: too many arguments.  last argument should be the number to compare to")
		os.Exit(3)
	}
	if len(goopt.Args) < 1 {
		log.Println("UNKNOWN: last argument should be the number to compare to")
		os.Exit(2)
	}
	if s, err := strconv.ParseFloat(goopt.Args[0], 64); err == nil {
		if criticalComparator, err := naglevelparse.Compile(*critical); err != nil {
			log.Printf("UNKNOWN: error parsing critical pattern %v: %#v", *critical, err)
			os.Exit(3)
		} else {
			if criticalComparator.Compare(s) {
				log.Printf("CRITICAL: %v not in range", s)
				os.Exit(2)
			}
		}
		if warningComparator, err := naglevelparse.Compile(*warning); err != nil {
			log.Printf("UNKNOWN: error parsing warning pattern %v: %#v", *warning, err)
			os.Exit(3)
		} else {
			if warningComparator.Compare(s) {
				log.Printf("WARNING: %v not in range", s)
				os.Exit(1)
			}
		}
		log.Printf("OK: %v within range", s)
	} else {
		log.Println("Error parsing float:", goopt.Args[0], err)
	}
}
