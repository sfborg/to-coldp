package config

import (
	"os"
	"path/filepath"
)

var (
	// repoURL is the URL to the SFGA schema repository.
	repoURL = "https://github.com/sfborg/sfga"

	// tag of the sfga repo to get correct schema version.
	verSFGA = "v0.3.32"

	// jobsNum is the default number of concurrent jobs to run.
	jobsNum = 5
)

type Config struct {
	// MinVersionSFGA sets minimal version of SFGA archive schema
	// that is needed for data extraction.
	MinVersionSFGA string

	// CacheDir keeps temporary directories for extracting and accessing
	// SFGA data.
	CacheDir string

	// CacheSfgaDir is where SFGA database is downloaded.
	CacheSfgaDir string

	// CacheColdpDir is used to place newly generated CoLDP files.
	CacheColdpDir string

	// JobsNum is the number of concurrent jobs to run.
	JobsNum int

	// BatchSize sets the size of batch for insert statements.
	BatchSize int

	// WithCommaDelim sets CSV delimeter as ',' (comma) instead of '\t' (tab).
	WithCommaDelim bool

	// WithNameUsage allows to create CoLDP where name, taxon, synonym
	// files are combined into name-usage file.
	WithNameUsage bool
}

type Option func(*Config)

func OptCacheDir(s string) Option {
	return func(c *Config) {
		c.CacheDir = s
	}
}

func OptWithNameUsage(b bool) Option {
	return func(c *Config) {
		c.WithNameUsage = b
	}
}

func OptWithTabDelim(b bool) Option {
	return func(c *Config) {
		c.WithCommaDelim = b
	}
}

func OptJobsNum(i int) Option {
	return func(c *Config) {
		c.JobsNum = i
	}
}

func OptBatchSize(i int) Option {
	return func(c *Config) {
		c.BatchSize = i
	}
}

func New(opts ...Option) Config {
	tmpDir := os.TempDir()
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = tmpDir
	}

	cacheDir = filepath.Join(cacheDir, "sfborg", "to", "coldp")

	res := Config{
		MinVersionSFGA: verSFGA,
		CacheDir:       cacheDir,
		BatchSize:      50_000,
		JobsNum:        4,
	}
	for _, o := range opts {
		o(&res)
	}

	res.CacheSfgaDir = filepath.Join(cacheDir, "sfga")
	res.CacheColdpDir = filepath.Join(cacheDir, "coldp")
	return res
}
