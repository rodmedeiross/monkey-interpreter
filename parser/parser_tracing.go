package parser

import (
	"fmt"
	"strings"
)

var traceLevel int = 0

const traceScape string = "\t"

func incTrace() { traceLevel++ }
func decTrace() { traceLevel-- }

func escapeTraceLevel() string {
	return strings.Repeat(traceScape, traceLevel-1)
}

func printTrace(tc string) {
	fmt.Printf("%s%s\n", escapeTraceLevel(), tc)
}

func trace(tc string) string {
	incTrace()
	printTrace("BEGIN " + tc)
	return tc
}

func untrace(tc string) {
	printTrace("END " + tc)
	decTrace()
}
