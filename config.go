package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type BaseConfig struct {
	// Directory containing VMX files
	VMXDir string `flag:"vmx" env:"STEMCELL_VMX_DIR"`

	// Optional temp directory for running VMs, if not specified the default
	// temp directory is used
	TempDir string `flag:"temp" env:"STEMCELL_TEMP_DIR"`

	// Path the packer template
	PackerTemplate string `flag:"template" env:"STEMCELL_PACKER_TEMPLATE"`

	// Administrator password of the Windows VM.
	AdminPassword string `flag:"password" env:"STEMCELL_ADMIN_PASSWORD"`

	// Print lots of debugging information.
	Debug bool `flag:"debug" env:"STEMCELL_DEBUG"`
}

func (c *BaseConfig) SetDefaults() {
	if c.TempDir == "" {
		c.TempDir = os.TempDir()
	}
}

type StemcellConfig struct {
	BaseConfig

	// Directory to save stemcells to, if not specified the current working
	// directory is used.
	StemcellDir string `flag:"stemcell" env:"STEMCELL_STEMCELL_DIR"`

	// Stemcell version must be of the form "[NUMBER].[NUMBER]"
	StemcellVersion string `flag:"version" env:"STEMCELL_STEMCELL_VERSION"`

	// Path the packer template
	PackerTemplate string `flag:"template" env:"STEMCELL_PACKER_TEMPLATE"`
}

func (c *StemcellConfig) fieldTag(field, key string) string {
	f, ok := reflect.TypeOf(*c).FieldByName(field)
	if !ok {
		panic("Config.envTag: invalid field name: " + field)
	}
	s, ok := f.Tag.Lookup(key)
	if !ok {
		panic("Config.envTag: field: " + field + " invalid tag key: " + key)
	}
	return s
}

func (c *StemcellConfig) envTag(field string) string {
	return c.fieldTag(field, "env")
}

func (c *StemcellConfig) flagTag(field string) string {
	return c.fieldTag(field, "flag")
}

func (c *StemcellConfig) ParseEnv() error {
	m := map[string]*string{
		c.envTag("VMXDir"):          &c.VMXDir,
		c.envTag("TempDir"):         &c.TempDir,
		c.envTag("StemcellDir"):     &c.StemcellDir,
		c.envTag("StemcellVersion"): &c.StemcellVersion,
		c.envTag("PackerTemplate"):  &c.PackerTemplate,
		c.envTag("AdminPassword"):   &c.AdminPassword,
	}
	for k, p := range m {
		if *p == "" {
			if s := os.Getenv(k); s != "" {
				*p = s
			}
		}
	}
	if !c.Debug {
		if s := os.Getenv(c.envTag("Debug")); s != "" {
			ok, err := strconv.ParseBool(s)
			if err != nil {
				return err
			}
			c.Debug = ok
		}
	}
	return nil
}

func (c *StemcellConfig) SetFlags(set *flag.FlagSet) {
	set.StringVar(&c.VMXDir, c.flagTag("VMXDir"), c.VMXDir, "Directory containing VMX files")
	set.StringVar(&c.TempDir, c.flagTag("TempDir"), c.TempDir, "Optional temp directory for running VMs")
	set.StringVar(&c.StemcellDir, c.flagTag("StemcellDir"), c.StemcellDir, "Directory where stemcell will be saved")
	set.StringVar(&c.StemcellVersion, c.flagTag("StemcellVersion"), c.StemcellVersion, "StemcellVersion")
	set.StringVar(&c.PackerTemplate, c.flagTag("PackerTemplate"), c.PackerTemplate, "PackerTemplate")
	set.StringVar(&c.AdminPassword, c.flagTag("AdminPassword"), c.AdminPassword, "AdminPassword")
	set.BoolVar(&c.Debug, c.flagTag("Debug"), false, "Debug")
}

func (c *StemcellConfig) Parse(arguments []string) error {
	c.SetDefaults()
	if err := c.ParseEnv(); err != nil {
		return err
	}
	set := flag.NewFlagSet("foo", flag.ExitOnError)
	c.SetFlags(set)
	return set.Parse(arguments)
}

func (c *StemcellConfig) Validate() []error {
	m := map[string]*string{
		"VMXDir":          &c.VMXDir,
		"StemcellDir":     &c.StemcellDir,
		"StemcellVersion": &c.StemcellVersion,
		"PackerTemplate":  &c.PackerTemplate,
		"AdminPassword":   &c.AdminPassword,
	}
	var errs []error
	for k, v := range m {
		if *v == "" {
			e := fmt.Errorf("missing argument: %s (flag: %q environment: %q)",
				k, c.flagTag(k), c.envTag(k))
			errs = append(errs, e)
		}
	}
	return errs
}
