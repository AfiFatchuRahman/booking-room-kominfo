package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

func baseURL() string {
	if v := os.Getenv("API_BASE_URL"); v != "" {
		return v
	}
	return "http://localhost:2254/api/v1"
}

func loginUser(t *testing.T, email, password string) string {
	t.Helper()

	payload, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})

	resp, err := http.Post(baseURL()+"/auth/login", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("gagal login %s: %v", email, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 saat login %s, got %d", email, resp.StatusCode)
	}

	var out struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("gagal decode token: %v", err)
	}
	return out.Token
}

func TestBookingOverlapAndConcurrency(t *testing.T) {
	adminToken := loginUser(t, "admin@test.com", "password123")
	user1Token := loginUser(t, "user1@test.com", "password123")
	user2Token := loginUser(t, "user2@test.com", "password123")
	user3Token := loginUser(t, "user3@test.com", "password123")
	user4Token := loginUser(t, "user4@test.com", "password123")
	user5Token := loginUser(t, "user5@test.com", "password123")
	user6Token := loginUser(t, "user6@test.com", "password123")
	user7Token := loginUser(t, "user7@test.com", "password123")
	user8Token := loginUser(t, "user8@test.com", "password123")
	user9Token := loginUser(t, "user9@test.com", "password123")
	user10Token := loginUser(t, "user10@test.com", "password123")

	client := &http.Client{Timeout: 10 * time.Second}

	roomPayload, _ := json.Marshal(map[string]interface{}{
		"name":     fmt.Sprintf("Ruang Test %d", time.Now().UnixNano()),
		"capacity": 12,
	})
	reqRoom, _ := http.NewRequest(http.MethodPost, baseURL()+"/rooms", bytes.NewBuffer(roomPayload))
	reqRoom.Header.Set("Content-Type", "application/json")
	reqRoom.Header.Set("Authorization", "Bearer "+adminToken)

	respRoom, err := client.Do(reqRoom)
	if err != nil {
		t.Fatalf("gagal membuat room: %v", err)
	}
	defer respRoom.Body.Close()

	if respRoom.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 saat buat room, got %d", respRoom.StatusCode)
	}

	var createdRoom struct {
		ID int64 `json:"id"`
	}
	_ = json.NewDecoder(respRoom.Body).Decode(&createdRoom)

	startTime := time.Now().UTC().Add(48 * time.Hour).Truncate(time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	const totalGoroutines = 10
	var wg sync.WaitGroup
	results := make(chan int, totalGoroutines)

	for i := 0; i < totalGoroutines; i++ {
		wg.Add(1)
		var token string

		switch i {
		case 0:
			token = user1Token
		case 1:
			token = user2Token
		case 2:
			token = user3Token
		case 3:
			token = user4Token
		case 4:
			token = user5Token
		case 5:
			token = user6Token
		case 6:
			token = user7Token
		case 7:
			token = user8Token
		case 8:
			token = user9Token
		case 9:
			token = user10Token
		}

		go func(idx int, bearer string) {
			defer wg.Done()

			body, _ := json.Marshal(map[string]interface{}{
				"room_id":    createdRoom.ID,
				"title":      fmt.Sprintf("Concurrent Meeting #%d", idx),
				"start_time": startTime.Format(time.RFC3339),
				"end_time":   endTime.Format(time.RFC3339),
			})

			req, _ := http.NewRequest(http.MethodPost, baseURL()+"/bookings", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+bearer)

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				t.Errorf("request %d error: %v", idx, err)
				return
			}
			defer resp.Body.Close()
			log.Print(resp)

			results <- resp.StatusCode
		}(i, token)
	}

	wg.Wait()
	close(results)

	createdCount := 0
	conflictCount := 0
	for code := range results {
		switch code {
		case http.StatusCreated:
			createdCount++
		case http.StatusConflict:
			conflictCount++
		default:
			t.Errorf("status code tidak terduga: %d", code)
		}
	}

	if createdCount != 1 || conflictCount != 9 {
		t.Fatalf("hasil concurrency tidak sesuai: created=%d (expected 1), conflict=%d (expected 9)", createdCount, conflictCount)
	}

	overlapBody, _ := json.Marshal(map[string]interface{}{
		"room_id":    createdRoom.ID,
		"title":      "Overlap Meeting",
		"start_time": startTime.Add(30 * time.Minute).Format(time.RFC3339),
		"end_time":   endTime.Add(30 * time.Minute).Format(time.RFC3339),
	})
	reqOverlap, _ := http.NewRequest(http.MethodPost, baseURL()+"/bookings", bytes.NewBuffer(overlapBody))
	reqOverlap.Header.Set("Content-Type", "application/json")
	reqOverlap.Header.Set("Authorization", "Bearer "+user1Token)

	respOverlap, err := client.Do(reqOverlap)
	if err != nil {
		t.Fatalf("gagal request overlap: %v", err)
	}
	defer respOverlap.Body.Close()

	if respOverlap.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 Conflict untuk jadwal beririsan, got %d", respOverlap.StatusCode)
	}
}
