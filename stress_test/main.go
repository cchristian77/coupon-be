package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// RequestBody Request payload structure
type RequestBody struct {
	UserID     string `json:"user_id"`
	CouponName string `json:"coupon_name"`
}

func main() {
	var (
		successCount    int64
		failureCount    int64
		totalDurationNs int64
		url             = "http://localhost:9000"
		totalRequests   = 50
	)

	client := &http.Client{Timeout: 5 * time.Second}

	var wg sync.WaitGroup
	wg.Add(totalRequests)

	for i := 0; i < totalRequests; i++ {
		go func() {
			defer wg.Done()

			startTime := time.Now()

			userID := fmt.Sprintf("user_%d", i)
			couponName := "COUPON_TEST"

			bodyBytes, err := json.Marshal(RequestBody{
				UserID:     userID,
				CouponName: couponName,
			})
			if err != nil {
				fmt.Println("JSON marshal error:", err)
				atomic.AddInt64(&failureCount, 1)
				return
			}

			endpoint := fmt.Sprintf("%s/api/coupons/claim", url)
			req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(bodyBytes))
			if err != nil {
				fmt.Println("Request creation error:", err)
				atomic.AddInt64(&failureCount, 1)
				return
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("HTTP error:", err)
				atomic.AddInt64(&failureCount, 1)
				return
			}

			respBodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("[USER=%s] Failed to read response body: %v\n", userID, err)
				atomic.AddInt64(&failureCount, 1)
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode >= 400 {
				atomic.AddInt64(&failureCount, 1)
				fmt.Printf("[USER=%s] FAILED | Status %d | Response : %s \n", userID, resp.StatusCode, string(respBodyBytes))
			} else if resp.StatusCode >= 200 {
				atomic.AddInt64(&successCount, 1)
				fmt.Printf("[USER=%s] SUCCESS | Status :  %d\n", userID, resp.StatusCode)
			}

			duration := time.Since(startTime)
			fmt.Printf("[USER=%s] Request duration : %s\n", userID, duration.String())
			atomic.AddInt64(&totalDurationNs, duration.Nanoseconds())
		}()

	}

	wg.Wait()

	fmt.Println("====== Results ======")
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Success: %d\n", successCount)
	fmt.Printf("Failed: %d\n", failureCount)
	fmt.Printf("Average Response Time: %s\n", time.Duration(totalDurationNs/int64(totalRequests)).String())
}
