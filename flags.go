package filters

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	filterFlag string
	chainFlag  string
	chainsFlag string
)

// Flags will generate filter flag, this should be called in your init
// It takes your appname as an argument and will also look in env-vars
// as FILTER_APP, CHAIN_APP, CHAINS_APP
// Not production ready, subject to change.
func Flags(app string) {
	appu := strings.ToUpper(app)
	flag.StringVar(
		&filterFlag,
		"filter",
		os.Getenv(fmt.Sprintf("FILTER_%v", appu)),
		fmt.Sprintf("Filter file for %v", app),
	)
	flag.StringVar(
		&chainFlag,
		"chain",
		os.Getenv(fmt.Sprintf("CHAIN_%v", appu)),
		fmt.Sprintf("Chain file for %v", app),
	)
	flag.StringVar(
		&chainFlag,
		"chains",
		os.Getenv(fmt.Sprintf("CHAINS_%v", appu)),
		fmt.Sprintf("Chains file for %v", app),
	)
}

// FilterF will return the Filter set in the flag
// Not production ready, subject to change.
func FilterF() (Filter, error) {
	parse()
	return FilterFile(filterFlag)
}

// ChainF will return the Chain set in the flag
// Not production ready, subject to change.
func ChainF() (Chain, error) {
	parse()
	return ChainFile(chainFlag)
}

// ChainsF will return the Chains set in the flag
// Not production ready, subject to change.
func ChainsF() (Chains, error) {
	parse()
	return ChainsFile(chainsFlag)
}

func parse() {
	if flag.Parsed() {
		return
	}
	flag.Parse()
}
