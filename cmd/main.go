package main

import (
	// Stdlib packages
	"os"

	// 3rd-Party packages for main
	flag "github.com/spf13/pflag"

	// Rays
	"github.com/BeamFoundry/rays"

	// Importing rays/lang
	. "github.com/BeamFoundry/rays/pkg/lang"

	// Initializing rays external modules
	_ "github.com/BeamFoundry/rays-fiber"
	_ "github.com/BeamFoundry/rays-orm"
	_ "github.com/BeamFoundry/rays-web-security"

	// Importing embedded resources
	"github.com/BeamFoundry/rays-sample-todoers/resources"

	// Initializing local packages
	_ "github.com/BeamFoundry/rays-sample-todoers/internal/configurations"
	_ "github.com/BeamFoundry/rays-sample-todoers/internal/repositories"
	_ "github.com/BeamFoundry/rays-sample-todoers/internal/services"
	_ "github.com/BeamFoundry/rays-sample-todoers/pkg/controllers"
	_ "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"
)

var (
	config    = []string{}
	profiles  = []string{}
	overrides = map[string]string{}
)

// This is a sample default main() function for rays.
//
// All flags can be used multiple times and/or values can be comma-separated.
//
// Example :
//
// ./todolist --profile testing,sqlite --profile admin
//
// This will activate the profiles "testing", "sqlite" and "admin".
//
// It let you parse profiles to use, additional configuration io/fs.FS directories and configuration overrides
// (from command-line or environment variables).
//
// NOTE: There is pluggable io/fs.FS implementation for reading from S3 or Hashcorp Vault and it's here
// you can configure and add them to the ApplicationContext.
//
// The configuration engine will search in all io/fs.FS directories for YAML files (ending with .yml or .yaml) starting with a lower letter
// and followed by letters. An optional hyphen and a profile name may be followed limit to this configuration to that profile.
//
// Example :
//
// values.yaml -> loaded by default
// values-test.yml -> loaded to profile "test"
// app/values.yml -> loaded by default
// app/MyValues.yml -> not loaded (file must start with a lowercase letter)
//
// The override flag let you define key=value pairs to replace configuration values.
// It use the same dot separated format for the key hierarchie as ConfigurationProperties
//
// For example, for a yaml configuration :
//
// app:
//
//	module:
//	  value: foo
//
// To replace the inner "value" key, you should define an override like :
//
//	--override app.module.config.value=bar
//
// The "overrides" system is designed to work with string values even if the target value is a bool or numeric value.
func main() {
	// Getting rays ApplicationContext
	ctx := rays.ApplicationContext()

	// Loading embedded configuration
	if err := ctx.AddConfigFS(resources.Config()); err != nil {
		err.LogError()
	}

	// Adding flags
	flag.StringSliceVarP(&config, "config", "c", []string{}, "Additional configuration paths")
	flag.StringSliceVarP(&profiles, "profile", "p", []string{}, "Profile to use, default : none")
	flag.StringToStringVar(&overrides, "override", map[string]string{}, "Overrides configuration key inline")
	flag.Parse()

	// Set profiles
	if len(profiles) > 0 {
		if err := ctx.SetActiveProfiles(profiles); err != nil {
			err.LogError()
		}
	}

	// Add additional local config directory
	for _, dir := range config {
		if stat, err := os.Stat(dir); os.IsExist(err) && stat.IsDir() {
			if err := ctx.AddConfigFS(os.DirFS(dir)); err != nil {
				err.LogError()
			}
		}
	}

	// Load any overrides defined
	for key, value := range overrides {
		if err := ctx.AddOverride(key, value); err != nil {
			err.LogError()
		}
	}

	// Load overrides with environment variables
	//
	// Only adding environment variables keys starting with 'RAYS_' or 'APP_',
	// lowering the case and replacing '_' with '.' .
	for _, e := range os.Environ() {
		env := Strings(e).SplitN("=", 2)
		key := String(env[0]).ToUpperCase()
		if key.StartsWith("RAYS_") || key.StartsWith("APP_") {
			key = key.Replace("_", ".").ToLowerCase()
			if err := ctx.AddOverride(key.String(), env[1]); err != nil {
				err.LogError()
			}
		}
	}

	// Initialize the container
	if err := ctx.Initialize(); err != nil {
		panic(err)
	}

	// Find and run "Runnables"
	ctx.Run()
}
