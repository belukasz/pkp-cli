package cmd

import (
	"github.com/belukasz/pkp-cli/app"
	"github.com/spf13/cobra"
)

var (
	startTime string
	days      int
	trainType string
	from      string
	to        string
	workdays  bool
	limit     int

	rootCmd = &cobra.Command{
		Use:   "pkp-cli",
		Short: "Handy PKP connections lookup",
		Run: func(cmd *cobra.Command, args []string) {
			startTime, err := cmd.Flags().GetString("start-time")
			if err != nil {
				panic(err)
			}
			days, err := cmd.Flags().GetInt("days")
			if err != nil {
				panic(err)
			}

			trainType, err := cmd.Flags().GetString("train-type")
			if err != nil {
				panic(err)
			}

			from, err := cmd.Flags().GetString("from")
			if err != nil {
				panic(err)
			}

			to, err := cmd.Flags().GetString("to")
			if err != nil {
				panic(err)
			}

			workdays, err := cmd.Flags().GetBool("workdays")
			if err != nil {
				panic(err)
			}

			limit, err := cmd.Flags().GetInt("limit")
			if err != nil {
				panic(err)
			}

			// scrapped := scrapper.ScrapeConnections(30, "06:00", "EIP", "krakow", "warszawa")
			// scrapper.PrintTable(scrapped, true, 3)
			scrapped := scrapper.ScrapeConnections(days, startTime, trainType, from, to)
			scrapper.PrintTable(scrapped, workdays, limit)
			// Do Stuff Here
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVar(&startTime, "start-time", "06:00", "Time after which we check for connections, example format: 06:00")
	rootCmd.Flags().IntVar(&days, "days", 30, "How many days to check")
	rootCmd.Flags().StringVar(&trainType, "train-type", "", "Type of train, one of: IC, EIP, EIPIC, ALL")
	rootCmd.Flags().StringVar(&from, "from", "krakow", "From station (example: krakow)")
	rootCmd.Flags().StringVar(&to, "to", "warszawa", "To station (example: warszawa)")
	rootCmd.Flags().BoolVar(&workdays, "workdays", true, "Only display workdays")
	rootCmd.Flags().IntVar(&limit, "limit", 6, "Limit of connections to display per day (not going to display more than 6 per day anyways)")
}
