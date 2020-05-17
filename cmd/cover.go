package cmd

import (
	"github.com/egon12/ghr/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// coverCmd represents the cover command
var coverCmd = &cobra.Command{
	Use:   "cover",
	Short: "Review the cover",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := app.Config{}
		_ = viper.Unmarshal(&cfg)

		a := app.InitApp(cfg)
		err := a.CoverageReviewer.Read(args[0])
		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(coverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// coverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// coverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
