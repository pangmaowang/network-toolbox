package cmd

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 扫描命令相关参数
var (
	target    string
	startPort int
	endPort   int
	timeout   time.Duration
)

// 定义扫描命令
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan ports for a given target",
	Long:  "Perform a port scan on a specified target IP address or hostname.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Scanning target %s from port %d to %d with timeout %v\n", target, startPort, endPort, timeout)

		// 执行端口扫描
		results := scanPorts(target, startPort, endPort, timeout)
		sort.Slice(results, func(i, j int) bool {
			return results[i].Port < results[j].Port
		})

		// 输出扫描结果
		for _, result := range results {
			status := "Closed"
			if result.Open {
				status = "Open"
			}
			fmt.Printf("Port %d: %s\n", result.Port, status)
		}
	},
}

// 初始化命令参数
func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&target, "target", "t", "127.0.0.1", "Target IP address or hostname")
	scanCmd.Flags().IntVarP(&startPort, "start", "s", 1, "Start port")
	scanCmd.Flags().IntVarP(&endPort, "end", "e", 1024, "End port")
	scanCmd.Flags().DurationVarP(&timeout, "timeout", "o", time.Second, "Timeout duration")
	viper.BindPFlag("target", scanCmd.Flags().Lookup("target"))
	viper.BindPFlag("start", scanCmd.Flags().Lookup("start"))
	viper.BindPFlag("end", scanCmd.Flags().Lookup("end"))
	viper.BindPFlag("timeout", scanCmd.Flags().Lookup("timeout"))
}

type scanResult struct {
	Port int
	Open bool
}

// 扫描端口

func scanPorts(target string, startPort, endPort int, timeout time.Duration) []scanResult {
	// Validate input parameters
	if net.ParseIP(target) == nil {
		// Try to resolve hostname
		_, err := net.LookupHost(target)
		if err != nil {
			fmt.Printf("Error: Invalid target address or unresolvable hostname: %s\n", target)
			return []scanResult{}
		}
	}

	if startPort < 1 || startPort > 65535 {
		fmt.Printf("Error: Start port must be between 1 and 65535\n")
		return []scanResult{}
	}

	if endPort < 1 || endPort > 65535 {
		fmt.Printf("Error: End port must be between 1 and 65535\n")
		return []scanResult{}
	}

	// Swap ports if start is greater than end
	if startPort > endPort {
		startPort, endPort = endPort, startPort
	}

	var wg sync.WaitGroup
	resultChan := make(chan scanResult, endPort-startPort+1)

	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			address := fmt.Sprintf("%s:%d", target, port)
			conn, err := net.DialTimeout("tcp", address, timeout)
			if err != nil {
				resultChan <- scanResult{Port: port, Open: false}
				return
			}
			conn.Close()
			resultChan <- scanResult{Port: port, Open: true}
		}(port)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var results []scanResult
	for result := range resultChan {
		results = append(results, result)
	}
	return results
}
