package app

import (
	"flag"
	"strings"
)

type multiValue []string

func (m *multiValue) String() string     { return strings.Join(*m, ",") }
func (m *multiValue) Set(v string) error { *m = append(*m, v); return nil }

// commonFlags are accepted by every subcommand. Flag values take precedence
// over environment variables, which take precedence over defaults.
type commonFlags struct {
	catalog    string
	dataDir    string
	stateDir   string
	cacheDir   string
	userHome   string
	json       bool
	allowRepos multiValue
}

func addCommonFlags(fs *flag.FlagSet, c *commonFlags) {
	fs.StringVar(&c.catalog, "catalog", "", "marketplace catalog URL (default "+DefaultCatalogURL+")")
	fs.BoolVar(&c.json, "json", false, "emit structured JSON on stdout; diagnostics go to stderr")
	fs.StringVar(&c.dataDir, "data-dir", "", "data directory override")
	fs.StringVar(&c.stateDir, "state-dir", "", "state directory override")
	fs.StringVar(&c.cacheDir, "cache-dir", "", "cache directory override")
	fs.StringVar(&c.userHome, "user-home", "", "user home override; agent skill roots resolve beneath this")
	fs.Var(&c.allowRepos, "allow-repository", "additionally allowed source repository URL (repeatable)")
}

func newFlagSet(name string, c *commonFlags) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(discardWriter{})
	addCommonFlags(fs, c)
	return fs
}

type discardWriter struct{}

func (discardWriter) Write(p []byte) (int, error) { return len(p), nil }
