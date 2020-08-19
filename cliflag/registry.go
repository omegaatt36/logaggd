package cliflag

import "github.com/urfave/cli"

// CliFlager is interface to describe a struct
// which is a set of options to singleton in package.
// This struct has method CliFlags to returns the options this package needed,
// and InitFromCliFlags that reads value from cli.Context and validate it.
type CliFlager interface {
	CliFlags() []cli.Flag
	InitFromCliFlags(c *cli.Context)
	Finalize()
}

var cliFlagers []CliFlager

// Register registers CliFlager, so we won't use a package without
// init it.
func Register(f CliFlager) {
	cliFlagers = append(cliFlagers, f)
}

// Globals returns flags from all registered packages.
func Globals() []cli.Flag {
	var res []cli.Flag
	for _, f := range cliFlagers {
		res = append(res, f.CliFlags()...)
	}

	return res
}

// Initialize inits all registered packages.
func Initialize(c *cli.Context) {
	for _, f := range cliFlagers {
		f.InitFromCliFlags(c)
	}
}

// Finalize finalizes registered packages, its execution order is reversed.
func Finalize() {
	for _, f := range cliFlagers {
		defer f.Finalize()
	}
}
