// Copyright 2018 fatedier, fatedier@gmail.com
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

package sub

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/fatedier/frp/models/config"
	"github.com/fatedier/frp/models/consts"
)

func init() {
	stcpCmd.PersistentFlags().StringVarP(&serverAddr, "server_addr", "s", "127.0.0.1:7000", "frp server's address")
	stcpCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "user")
	stcpCmd.PersistentFlags().StringVarP(&protocol, "protocol", "p", "tcp", "tcp or kcp")
	stcpCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "auth token")
	stcpCmd.PersistentFlags().StringVarP(&logLevel, "log_level", "", "info", "log level")
	stcpCmd.PersistentFlags().StringVarP(&logFile, "log_file", "", "console", "console or file path")
	stcpCmd.PersistentFlags().IntVarP(&logMaxDays, "log_max_days", "", 3, "log file reversed days")

	stcpCmd.PersistentFlags().StringVarP(&proxyName, "proxy_name", "n", "", "proxy name")
	stcpCmd.PersistentFlags().StringVarP(&role, "role", "", "server", "role")
	stcpCmd.PersistentFlags().StringVarP(&sk, "sk", "", "", "secret key")
	stcpCmd.PersistentFlags().StringVarP(&serverName, "server_name", "", "", "server name")
	stcpCmd.PersistentFlags().StringVarP(&localIp, "local_ip", "i", "127.0.0.1", "local ip")
	stcpCmd.PersistentFlags().IntVarP(&localPort, "local_port", "l", 0, "local port")
	stcpCmd.PersistentFlags().StringVarP(&bindAddr, "bind_addr", "", "", "bind addr")
	stcpCmd.PersistentFlags().IntVarP(&bindPort, "bind_port", "", 0, "bind port")
	stcpCmd.PersistentFlags().BoolVarP(&useEncryption, "ue", "", false, "use encryption")
	stcpCmd.PersistentFlags().BoolVarP(&useCompression, "uc", "", false, "use compression")

	rootCmd.AddCommand(stcpCmd)
}

var stcpCmd = &cobra.Command{
	Use:   "stcp",
	Short: "Run frpc with a single stcp proxy",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := parseClientCommonCfg(CfgFileTypeCmd, "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cfg := &config.StcpProxyConf{}
		var prefix string
		if user != "" {
			prefix = user + "."
		}
		cfg.ProxyName = prefix + proxyName
		cfg.ProxyType = consts.StcpProxy
		cfg.Role = role
		cfg.Sk = sk
		cfg.ServerName = serverName
		cfg.LocalIp = localIp
		cfg.LocalPort = localPort
		cfg.BindAddr = bindAddr
		cfg.BindPort = bindPort
		cfg.UseEncryption = useEncryption
		cfg.UseCompression = useCompression

		err = cfg.CheckForCli()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if cfg.Role == "server" {
			proxyConfs := map[string]config.ProxyConf{
				cfg.ProxyName: cfg,
			}
			err = startService(proxyConfs, nil)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			visitorConfs := map[string]config.ProxyConf{
				cfg.ProxyName: cfg,
			}
			err = startService(nil, visitorConfs)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		return nil
	},
}
