package main

import (
	"fmt"
	"github.com/fossas/fossa-cli/analyzers"
	"github.com/fossas/fossa-cli/module"
	"github.com/fossas/fossa-cli/pkg"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"regexp"
)

type V3Config struct {
	Targets Targets `yaml:"targets,omitempty" json:"targets,omitempty"`
}

type Targets struct {
	Exclude []Target `yaml:"exclude,omitempty" json:"exclude,omitempty"`
}

type Target struct {
	Type string `yaml:"type,omitempty" json:"type,omitempty"`
	Path string `yaml:"path,omitempty" json:"path,omitempty"`
}

// SuspiciousModules is similar to FilterSuspiciousModules, but it returns only
// the suspicious modules rather than removing them from the result
func SuspiciousModules(modules []module.Module) (targets []Target) {
	suspicious := regexp.MustCompile(`docs?[/\\]|[Tt]est|examples?|vendor[/\\]|site-packages[/\\]|node_modules[/\\]|.srclib-cache[/\\]|spec[/\\]|Godeps[/\\]|.git[/\\]|bower_components[/\\]|third[_-]party[/\\]|tmp[/\\]|Carthage[/\\]Checkouts[/\\]`)
	for _, m := range modules {
		matched := suspicious.MatchString(m.Dir)

		// For Go, filter out non-executable packages.
		if matched || (m.Type == pkg.Go && !m.IsExecutable) {
			targets = append(targets, Target{Path: m.Dir, Type: m.Type.String()})
		}
	}
	return
}

func main() {
	app := &cli.App{
		Name:  "fossa-1to3",
		Usage: "Utilities to help migrate projects integrated with FOSSA CLI v1 to v3",
		Commands: []*cli.Command{
			{
				Name:  "targets",
				Usage: "Generate a v3-compatible config file to exclude targets that v1 would have implicitly excluded from analysis",
				Action: func(c *cli.Context) error {
					targetDir := c.Args().First()
					if targetDir == "" {
						targetDir = "."
					}
					dir, err := os.Stat(targetDir)
					if err != nil {
						return err
					}
					if !dir.IsDir() {
						return fmt.Errorf("path is not a directory: %v", dir.Name())
					}
					modules, err := analyzers.Discover(targetDir, map[string]interface{}{})
					filtered := SuspiciousModules(modules)
					config := V3Config{
						Targets{
							Exclude: filtered,
						},
					}
					if len(config.Targets.Exclude) > 0 {
						result, err := yaml.Marshal(&config)
						if err != nil {
							return err
						}
						fmt.Println(string(result))
					}
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
