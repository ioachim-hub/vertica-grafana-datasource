//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"

	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
)

func init() {
	_ = build.SetBeforeBuildCallback(func(cfg build.Config) (build.Config, error) {
		buildInfoJSON := os.Getenv("GRAFANA_BUILD_INFO_JSON")
		if buildInfoJSON == "" {
			return cfg, nil
		}
		if cfg.CustomVars == nil {
			cfg.CustomVars = make(map[string]string)
		}
		cfg.CustomVars["github.com/grafana/grafana-plugin-sdk-go/build/buildinfo.buildInfoJSON"] = buildInfoJSON
		return cfg, nil
	})
}

// Hello prints a message (shows that you can define custom Mage targets).
func Hello() {
	fmt.Println("hello plugin developer!")
}

// Default configures the default target.
var Default = build.Build.Linux
