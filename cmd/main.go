package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thanethomson/alien-invasion/pkg/aliensim"
)

// Command line flags
var (
	flagAlienCount       int
	flagUseExampleMap    bool
	flagWorldMapFilename string
)

var rootCmd = &cobra.Command{
	Use:   "alien-invasion",
	Short: "Alien invasion simulator",
	Long:  "Alien invasion simulator! See https://github.com/thanethomson/alien-invasion for more details.",
	Run: func(cmd *cobra.Command, args []string) {
		var reader io.Reader
		var err error

		if flagUseExampleMap {
			reader = strings.NewReader(aliensim.ExampleWorld)
			fmt.Println("Using example world for simulation.")
		} else {
			reader, err = os.Open(flagWorldMapFilename)
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
			fmt.Println(fmt.Sprintf("Reading world from file: %s", flagWorldMapFilename))
		}

		fmt.Println(fmt.Sprintf("Executing simulation with %d aliens...", flagAlienCount))

		sim := aliensim.NewSimulation(
			aliensim.NewSimulationConfig(
				reader,
				flagAlienCount,
			),
		)

		res, err := sim.Simulate()
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}
		fmt.Println("")
		fmt.Println("Done. Remaining aliens:")
		for _, alien := range res.FinalAliens {
			fmt.Println(alien)
		}
		fmt.Println("")
		fmt.Println("Final world map:")
		fmt.Println("")
		fmt.Println(res.FinalMap.Render())
	},
}

func initCmd() {
	rootCmd.PersistentFlags().IntVarP(
		&flagAlienCount,
		"alien-count",
		"N",
		2,
		"the number of aliens to simulate",
	)
	rootCmd.PersistentFlags().StringVarP(
		&flagWorldMapFilename,
		"world-map",
		"m",
		"world-map.txt",
		"the file from which to load the world map",
	)
	rootCmd.PersistentFlags().BoolVar(
		&flagUseExampleMap,
		"use-example-map",
		false,
		"use the example world map instead of loading one",
	)
}

func main() {
	initCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
