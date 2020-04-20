/*
Copyright © 2020 Egon Firman <egon.firman@gmail.com>

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
	"fmt"

	"github.com/egon12/ghr/app"
	"github.com/egon12/ghr/path"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// commentCmd represents the comment command
var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, line, err := path.ParseFileAndLine(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		cfg := app.Config{}
		_ = viper.Unmarshal(&cfg)

		a := app.InitApp(cfg)
		c := a.CommitSource.GetCurrentCommit()
		a.ReviewProcess.Start(c)
		a.ReviewProcess.AddComment(filePath, line, args[1])
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
