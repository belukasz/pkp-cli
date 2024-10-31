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
	weekends  bool
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

			weekends, err := cmd.Flags().GetBool("weekends")
			if err != nil {
				panic(err)
			}

			limit, err := cmd.Flags().GetInt("limit")
			if err != nil {
				panic(err)
			}
			scrapped := scrapper.ScrapeConnections(days, startTime, trainType, from, to)
			scrapper.PrintTable(scrapped, !weekends, limit)
		},

	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVar(&startTime, "start-time", "06:00", "Time after which we check for connections, example format: 06:00")
	rootCmd.Flags().IntVar(&days, "days", 30, "How many days to check")
	rootCmd.Flags().StringVar(&trainType, "train-type", "EIP", "Type of train, one of: IC, EIP, EIPIC, ALL")
	rootCmd.Flags().StringVar(&from, "from", "krakow", "From station")
	rootCmd.Flags().StringVar(&to, "to", "warszawa", "To station")
	rootCmd.Flags().BoolVar(&weekends, "weekends", false, "Only display workdays")
	rootCmd.Flags().IntVar(&limit, "limit", 10, "Limit of connections to display per day (not going to display more than 6 per day anyways)")
}
