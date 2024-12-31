package cmd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// 定义扫描命令
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Perform a ping on a specified target IP address or hostname.",
	Long:  `Perform a ping on a specified target IP address or hostname and display the results.`,
	Run: func(cmd *cobra.Command, args []string) {
		target, err := cmd.Flags().GetString("target")
		if err != nil {
			fmt.Printf("Can not parse the target: %v\n", err)
		}
		timeout, err := cmd.Flags().GetDuration("timeout")
		if err != nil {
			fmt.Printf("Can not parse the timeout: %v\n", err)
		}
		count, err := cmd.Flags().GetInt("count")
		if err != nil {
			fmt.Printf("Can not parse the count: %v\n", err)
		}

		advanced, _ := cmd.Flags().GetBool("advanced")
		if advanced {
			maxHops, _ := cmd.Flags().GetInt("max-hops")
			traceRoute(target, maxHops, timeout)
		} else {
			runPingCommand(target, count)
		}
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().StringP("target", "t", "127.0.0.1", "Target IP address or hostname")
	pingCmd.Flags().DurationP("timeout", "o", 1*time.Second, "Timeout in seconds")
	pingCmd.Flags().IntP("count", "c", 5, "Number of ping attempts")
	pingCmd.Flags().BoolP("advanced", "a", false, "Enable advanced mode for traceroute analysis")
	pingCmd.Flags().IntP("max-hops", "m", 30, "Maximum number of hops (used in advanced mode)")
}

func traceRoute(target string, maxHops int, timeout time.Duration) {
	ipAddr, err := net.ResolveIPAddr("ip4", target)
	if err != nil {
		fmt.Printf("Failed to resolve target %s: %v\n", target, err)
		return
	}

	fmt.Printf("Traceroute to %s (%s), %d hops max\n", target, ipAddr.IP.String(), maxHops)

	for ttl := 1; ttl <= maxHops; ttl++ {
		conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
		if err != nil {
			if errors.Is(err, syscall.EPERM) {
				fmt.Println("Permission denied. Please run the program as administrator.")
			} else {
				fmt.Printf("Failed to listen to ICMP: %v\n", err)
			}
			return
		}
		defer conn.Close()

		// 设置 TTL
		if err := conn.IPv4PacketConn().SetTTL(ttl); err != nil {
			fmt.Printf("Failed to set TTL: %v\n", err)
			return
		}

		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  ttl,
				Data: []byte("traceroute"),
			},
		}

		msgBytes, err := msg.Marshal(nil)
		if err != nil {
			fmt.Printf("Failed to encode ICMP message: %v\n", err)
			continue
		}

		start := time.Now()
		if _, err := conn.WriteTo(msgBytes, ipAddr); err != nil {
			fmt.Printf("Hop %d: Failed to send ICMP packet\n", ttl)
			continue
		}

		conn.SetReadDeadline(time.Now().Add(timeout))
		reply := make([]byte, 1500)
		n, addr, err := conn.ReadFrom(reply)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("Hop %d: * (Timeout)\n", ttl)
		} else {
			response, err := icmp.ParseMessage(1, reply[:n])
			if err == nil && response.Type == ipv4.ICMPTypeTimeExceeded {
				fmt.Printf("Hop %d: %v, Time: %v\n", ttl, addr, duration)
			} else if err == nil && response.Type == ipv4.ICMPTypeEchoReply {
				fmt.Printf("Hop %d: %v, Time: %v (Reached Target)\n", ttl, addr, duration)
				break
			} else {
				fmt.Printf("Hop %d: Unexpected response\n", ttl)
			}
		}
	}
	fmt.Println("Traceroute complete.")
}

func runPingCommand(host string, count int) error {
	if ip := net.ParseIP(host); ip == nil {
		ips, err := net.LookupIP(host)
		if err != nil {
			fmt.Printf("Unable to resolve hostname or invalid IP: %s\n", host)
			fmt.Printf("Status: Failed\n")
			return fmt.Errorf("invalid host or hostname: %s", host)
		}
		host = ips[0].String()
		fmt.Printf("Resolved %s to %s\n", host, ips[0].String())
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", "-n", fmt.Sprintf("%d", count), host)
	case "darwin", "linux":
		cmd = exec.Command("ping", "-c", fmt.Sprintf("%d", count), host)
	default:
		fmt.Printf("Unsupported operating system: %s\n", runtime.GOOS)
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to execute system ping command: %v\n", err)
		fmt.Printf("Status: Failed\n")
	} else {
		fmt.Printf("System ping command result:\n%s\n", string(output))
	}

	return nil
}
