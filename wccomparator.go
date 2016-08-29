package nagiosrangeparser

import (
	"fmt"
)

func NagiosComparator(warning, critical string, value float64) (level, message string, rc int) {
	if criticalComparator, err := Compile(critical); err != nil {
	        return "UNKNOWN", fmt.Sprintf("error parsing critical pattern %v: %#v", critical, err), 3
	} else {
	        if criticalComparator.Compare(value) {
	                return "CRITICAL", "", 2
	        }
	}
	if warningComparator, err := Compile(warning); err != nil {
	        return "UNKNOWN", fmt.Sprintf("error parsing warning pattern %v: %#v", warning, err), 3
	} else {
	        if warningComparator.Compare(value) {
	                return "WARNING", "", 1
	        }
	}
	return "OK", "", 0
}
