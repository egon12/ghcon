/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

	"github.com/egon12/ghr/app"
	"github.com/egon12/ghr/testresult"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// consumeTestCmd represents the consumeTest command
var consumeTestCmd = &cobra.Command{
	Use:   "consumeTest outpath errpath",
	Short: "consume stdout and stderr that produced by go test",
	Long:  "",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		content, err := testresult.Read(args[0], args[1])
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		if len(content) == 0 {
			return
		}

		cfg := app.Config{}
		_ = viper.Unmarshal(&cfg)

		a := app.InitApp(cfg)
		err = a.ReviewProcess.Finish(content)
		if err != nil {
			cmd.PrintErr(err)
		}
		// so it will failed in jenkins
		os.Exit(2)
	},
}

func init() {
	rootCmd.AddCommand(consumeTestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consumeTestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consumeTestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
