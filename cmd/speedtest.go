package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/spf13/cobra"
)

var uploadServerPoolCredentials = map[string]Credentials{
	"ftp://ftp.dlptest.com": {
		userName: "dlpuser",
		password: "rNrKYTX9g7z3RgJRmxWuGHbeu",
	},
}

// 定义扫描命令
var speedTestCmd = &cobra.Command{
	Use:   "speedtest",
	Short: "Perform an internet speed test",
	Long:  `Measure your internet connection's download and upload speeds using speedtest.`,
	Run: func(cmd *cobra.Command, args []string) {

		speedTest()
	},
}

func init() {
	rootCmd.AddCommand(speedTestCmd)

}
func speedTest() {
	var downloadServerPool = []string{
		"https://speed.cloudflare.com/__down?bytes=104857600",
		"https://sgp.proof.ovh.net/files/100Mb.dat",
	}

	var uploadServerPool = []string{
		"ftp://ftp.dlptest.com",
	}
	fmt.Println("Starting download speed test...")
	downloadSpeed, err := testDownloadSpeed(downloadServerPool[0])
	if err != nil {
		fmt.Printf("Error during download test: %v\n", err)
		return
	}
	fmt.Println("Starting upload speed test...")
	uploadSpeed, err := testUploadSpeed(uploadServerPool[0], uploadServerPoolCredentials[uploadServerPool[0]])

	if err != nil {
		fmt.Printf("Error during upload test: %v\n", err)
		return
	}

	fmt.Printf("Download Speed: %.2f Mbps\n", downloadSpeed)
	fmt.Printf("Upload Speed: %.2f Mbps\n", uploadSpeed)
}

func isValid(urlStr string) bool {
	_, err := url.Parse(urlStr)
	return err == nil
}

func testDownloadSpeed(serverURL string) (float64, error) {
	// 1. 验证URL合法性
	if !isValid(serverURL) {
		return 0, fmt.Errorf("invalid URL: %s", serverURL)
	}

	// 2. 初始化HTTP客户端
	client := &http.Client{Timeout: 60 * time.Second}

	// 3. 记录开始时间
	start := time.Now()
	response, err := client.Get(serverURL)
	if err != nil {
		return 0, fmt.Errorf("download failed: %w", err)
	}
	defer response.Body.Close()

	// 4. 确保响应状态码为200
	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("server error: status code %d", response.StatusCode)
	}

	// 5. 计算接收字节数并丢弃数据
	totalBytes, err := io.Copy(io.Discard, response.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	// 6. 计算时间和速度
	duration := time.Since(start).Seconds()
	if duration == 0 || totalBytes == 0 {
		return 0, fmt.Errorf("invalid measurement: duration=%f, bytes=%d", duration, totalBytes)
	}

	// 7. 返回Mbps速度
	speed := (float64(totalBytes) * 8) / (duration * 1024 * 1024)
	return speed, nil
}

func testUploadSpeed(serverUrl string, credentials Credentials) (float64, error) {
	fmt.Println("Starting upload speed test to:", serverUrl)

	if !isValid(serverUrl) {
		return 0, fmt.Errorf("invalid URL: %s", serverUrl)
	}

	ftpServer := strings.TrimPrefix(serverUrl, "ftp://") + ":21"
	fmt.Printf("Connecting to FTP server: %s\n", ftpServer)

	ftpClient, err := ftp.Dial(ftpServer, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return 0, fmt.Errorf("failed to connect to FTP server: %w", err)
	}
	defer ftpClient.Quit()

	fmt.Printf("Attempting login with username: %s\n", credentials.userName)
	err = ftpClient.Login(credentials.userName, credentials.password)
	if err != nil {
		return 0, fmt.Errorf("failed to login to FTP server: %w", err)
	}
	fmt.Println("Successfully logged in to FTP server")

	// Create test data (100MB)
	dataSize := 100 * 1024 * 1024 // 100MB
	fmt.Printf("Preparing %d MB of test data\n", dataSize/1024/1024)
	data := make([]byte, dataSize)
	for i := range data {
		data[i] = 'A'
	}

	testFileName := fmt.Sprintf("speedtest_%d.dat", time.Now().Unix())
	fmt.Printf("Starting upload of test file: %s\n", testFileName)

	start := time.Now()
	err = ftpClient.Stor(testFileName, bytes.NewReader(data))
	if err != nil {
		return 0, fmt.Errorf("failed to upload file: %w", err)
	}

	duration := time.Since(start).Seconds()
	uploadSpeed := (float64(len(data)) * 8 / duration) / (1024 * 1024)
	fmt.Printf("Upload completed in %.2f seconds\n", duration)

	fmt.Printf("Cleaning up - deleting test file: %s\n", testFileName)
	err = ftpClient.Delete(testFileName)
	if err != nil {
		return 0, fmt.Errorf("failed to delete test file: %w", err)
	}

	fmt.Printf("Raw upload speed calculated: %.2f Mbps\n", uploadSpeed)
	return uploadSpeed, nil
}

type Credentials struct {
	userName string
	password string
}
