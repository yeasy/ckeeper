// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yeasy/ckeeper/engine"
	"github.com/spf13/viper"
	"time"
	"runtime"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the ckpeeper service",
	Long:  `It will read the config and run accoridingly.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	RootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func run() {
	// TODO: Work your own magic here
	handler := engine.NewHanlder()
	var mem   runtime.MemStats

	for {
		interval := time.Duration(viper.GetInt("check.interval"))
		logger.Infof(">>>Start health check, interval = %d seconds\n", interval)
		taskStart := time.Now()
		handler.Load()
		handler.Process()
		taskEnd := time.Now()
		taskTime := taskEnd.Sub(taskStart)

		//runtime.GC()

		runtime.ReadMemStats(&mem)
		logger.Infof("<<<End health check. Time = %s, interval= %d seconds, memory usage = %d KB.\n\n", taskTime, interval, mem.Alloc/1024)
		time.Sleep(interval * time.Second)
	}
}

func processRule() {

}
