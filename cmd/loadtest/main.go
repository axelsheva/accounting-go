package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	apiBaseURL = "http://localhost:8081/api"
)

// Parameters from command line
var (
	numUsers        = flag.Int("users", 100, "Number of users to create")
	numTransactions = flag.Int("transactions", 1000, "Number of transactions per user")
	concurrency     = flag.Int("concurrency", 50, "Number of concurrent goroutines")
	verbose         = flag.Bool("verbose", false, "Enable verbose output")
	baseURL         = flag.String("url", apiBaseURL, "Base URL for the API")
)

// Structures for requests and responses
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type CreateUserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type CreateTransactionRequest struct {
	UserID   int     `json:"user_id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Type     string  `json:"type"`
}

type CreateTransactionResponse struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

// Metrics for collecting statistics
type Metrics struct {
	sync.Mutex
	TotalRequests      int
	SuccessfulRequests int
	FailedRequests     int
	TotalResponseTime  time.Duration
	MinResponseTime    time.Duration
	MaxResponseTime    time.Duration
	ResponseTimes      []time.Duration
}

func (m *Metrics) AddResponseTime(duration time.Duration) {
	m.Lock()
	defer m.Unlock()
	m.TotalRequests++
	m.TotalResponseTime += duration
	m.ResponseTimes = append(m.ResponseTimes, duration)

	if m.MinResponseTime == 0 || duration < m.MinResponseTime {
		m.MinResponseTime = duration
	}
	if duration > m.MaxResponseTime {
		m.MaxResponseTime = duration
	}
}

func (m *Metrics) AddSuccess() {
	m.Lock()
	defer m.Unlock()
	m.SuccessfulRequests++
}

func (m *Metrics) AddFailure() {
	m.Lock()
	defer m.Unlock()
	m.FailedRequests++
}

func (m *Metrics) CalculatePercentiles() map[string]time.Duration {
	m.Lock()
	defer m.Unlock()

	if len(m.ResponseTimes) == 0 {
		return map[string]time.Duration{}
	}

	// Sort response times
	sorted := make([]time.Duration, len(m.ResponseTimes))
	copy(sorted, m.ResponseTimes)
	slices.Sort(sorted)

	n := len(sorted)
	p50 := sorted[n*50/100]
	p90 := sorted[n*90/100]
	p95 := sorted[n*95/100]
	p99 := sorted[n*99/100]

	return map[string]time.Duration{
		"p50": p50,
		"p90": p90,
		"p95": p95,
		"p99": p99,
	}
}

func (m *Metrics) AverageResponseTime() time.Duration {
	m.Lock()
	defer m.Unlock()
	if m.TotalRequests == 0 {
		return 0
	}
	return m.TotalResponseTime / time.Duration(m.TotalRequests)
}

// Function for creating a user
func createUser(client *http.Client, metrics *Metrics) (*CreateUserResponse, error) {
	name := fmt.Sprintf("User-%s", uuid.New().String()[:8])
	email := fmt.Sprintf("%s@example.com", strings.ToLower(name))

	userReq := CreateUserRequest{
		Name:  name,
		Email: email,
		Age:   rand.Intn(80) + 18,
	}

	reqBody, err := json.Marshal(userReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user request: %w", err)
	}

	start := time.Now()
	req, err := http.NewRequest("POST", *baseURL+"/users", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	elapsed := time.Since(start)

	metrics.AddResponseTime(elapsed)

	if err != nil {
		metrics.AddFailure()
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		metrics.AddFailure()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var userResp CreateUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		metrics.AddFailure()
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	metrics.AddSuccess()
	return &userResp, nil
}

// Function for creating a transaction
func createTransaction(client *http.Client, userID int, metrics *Metrics) (*CreateTransactionResponse, error) {
	// Generate random transaction data
	amount := rand.Float64() * 1000
	currencies := []string{"USD", "EUR", "GBP", "JPY"}
	currency := currencies[rand.Intn(len(currencies))]
	types := []string{"deposit", "withdrawal"}
	txType := types[0]

	txReq := CreateTransactionRequest{
		UserID:   userID,
		Amount:   amount,
		Currency: currency,
		Type:     txType,
	}

	reqBody, err := json.Marshal(txReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction request: %w", err)
	}

	start := time.Now()
	req, err := http.NewRequest("POST", *baseURL+"/transactions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	elapsed := time.Since(start)

	metrics.AddResponseTime(elapsed)

	if err != nil {
		metrics.AddFailure()
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		metrics.AddFailure()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var txResp CreateTransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&txResp); err != nil {
		metrics.AddFailure()
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	metrics.AddSuccess()
	return &txResp, nil
}

// Worker function for processing users
func worker(id int, userIDs <-chan int, wg *sync.WaitGroup, client *http.Client, metrics *Metrics) {
	defer wg.Done()

	for userID := range userIDs {
		// Create transactions for each user
		for i := 0; i < *numTransactions; i++ {
			tx, err := createTransaction(client, userID, metrics)
			if err != nil {
				log.Printf("Worker %d: error creating transaction for user %d: %v", id, userID, err)
				continue
			}
			if *verbose {
				log.Printf("Worker %d: created transaction %s for user %d, type: %s, amount: %.2f %s",
					id, tx.ID, userID, tx.Type, tx.Amount, tx.Currency)
			}
		}
	}
}

func main() {
	flag.Parse()
	log.Printf("Starting load test with %d users, %d transactions per user, %d parallel goroutines",
		*numUsers, *numTransactions, *concurrency)

	// Create a properly configured HTTP client with connection pooling
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   100, // This is important - default is 2
		MaxConnsPerHost:       0,   // 0 means no limit
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 30,
	}

	// Initialize metrics
	userMetrics := &Metrics{
		MinResponseTime: time.Hour, // Initialize with a large value
	}
	txMetrics := &Metrics{
		MinResponseTime: time.Hour,
	}

	// Create users
	users := make([]*CreateUserResponse, 0, *numUsers)
	log.Printf("Creating %d users...", *numUsers)
	for range make([]struct{}, *numUsers) {
		user, err := createUser(client, userMetrics)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			continue
		}
		users = append(users, user)
		if *verbose {
			log.Printf("Created user: ID=%d, Name=%s, Email=%s", user.ID, user.Name, user.Email)
		}
	}

	log.Printf("Successfully created %d users", len(users))
	if len(users) == 0 {
		log.Fatalf("No users for transaction generation, exiting")
	}

	// Create channel with user IDs
	userIDs := make(chan int, len(users))
	for _, user := range users {
		userIDs <- user.ID
	}
	close(userIDs)

	// Start workers for creating transactions
	log.Printf("Generating transactions for %d users using %d workers...", len(users), *concurrency)
	startTime := time.Now()

	var wg sync.WaitGroup
	wg.Add(*concurrency)
	for i := 0; i < *concurrency; i++ {
		go worker(i+1, userIDs, &wg, client, txMetrics)
	}
	wg.Wait()

	// Calculate and display metrics
	userAvgTime := userMetrics.AverageResponseTime()
	txAvgTime := txMetrics.AverageResponseTime()
	userPercentiles := userMetrics.CalculatePercentiles()
	txPercentiles := txMetrics.CalculatePercentiles()

	totalTime := time.Since(startTime)
	totalRequests := userMetrics.TotalRequests + txMetrics.TotalRequests
	successRate := float64(userMetrics.SuccessfulRequests+txMetrics.SuccessfulRequests) / float64(totalRequests) * 100
	rps := float64(totalRequests) / totalTime.Seconds()

	// Output the results in markdown format
	fmt.Println("\n### General Results")
	fmt.Println()
	fmt.Println("| Metric                    | Value                                     |")
	fmt.Println("| ------------------------- | ----------------------------------------- |")
	fmt.Printf("| Total execution time      | %s                             |\n", totalTime)
	fmt.Printf("| Total requests            | %d (Users: %d, Transactions: %d) |\n",
		totalRequests, userMetrics.TotalRequests, txMetrics.TotalRequests)
	fmt.Printf("| Successful requests       | %d (%.2f%%)                          |\n",
		userMetrics.SuccessfulRequests+txMetrics.SuccessfulRequests, successRate)
	fmt.Printf("| Failed requests           | %d (%.2f%%)                               |\n",
		userMetrics.FailedRequests+txMetrics.FailedRequests, 100-successRate)
	fmt.Printf("| Requests per second (RPS) | %.2f                                   |\n", rps)

	fmt.Println("\n### User Creation Request Metrics")
	fmt.Println()
	fmt.Println("| Metric                    | Value                                                       |")
	fmt.Println("| ------------------------- | ----------------------------------------------------------- |")
	fmt.Printf("| Average response time     | %s                                                     |\n", userAvgTime)
	fmt.Printf("| Minimum response time     | %s                                                    |\n", userMetrics.MinResponseTime)
	fmt.Printf("| Maximum response time     | %s                                                  |\n", userMetrics.MaxResponseTime)
	fmt.Printf("| Response time percentiles | P50=%s, P90=%s, P95=%s, P99=%s |\n",
		userPercentiles["p50"], userPercentiles["p90"], userPercentiles["p95"], userPercentiles["p99"])

	fmt.Println("\n### Transaction Creation Request Metrics")
	fmt.Println()
	fmt.Println("| Metric                    | Value                                                       |")
	fmt.Println("| ------------------------- | ----------------------------------------------------------- |")
	fmt.Printf("| Average response time     | %s                                                     |\n", txAvgTime)
	fmt.Printf("| Minimum response time     | %s                                                    |\n", txMetrics.MinResponseTime)
	fmt.Printf("| Maximum response time     | %s                                                  |\n", txMetrics.MaxResponseTime)
	fmt.Printf("| Response time percentiles | P50=%s, P90=%s, P95=%s, P99=%s |\n",
		txPercentiles["p50"], txPercentiles["p90"], txPercentiles["p95"], txPercentiles["p99"])
}
