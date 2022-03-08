/*
 Copyright 2022 The KubeSphere Authors.
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

package ctl

import (
	"flag"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/kubesphere/kubeeye/pkg/audit"
)

var KubeConfig string
var additionalregoruleputh string
var output string

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "扫描kubernetes集群中的风险项",
	Run: func(cmd *cobra.Command, args []string) {
		err := audit.Cluster(cmd.Context(), KubeConfig, additionalregoruleputh, output)
		if err != nil {
			glog.Fatalf("run audit failed with error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	rootCmd.PersistentFlags().StringVarP(&KubeConfig, "config", "f", "", "可指定 kubeconfig 路径，默认$HOME/.kube/config")
	// auditCmd.PersistentFlags().StringVarP(&additionalregoruleputh, "additional-rego-rule-path", "p", "", "Specify the path of additional rego rule files directory.")
	auditCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "可选 JSON 和 CSV为输出格式")
}
