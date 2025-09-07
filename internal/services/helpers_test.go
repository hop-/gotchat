package services

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestAtomicRunningStatus_setRunningStatus(t *testing.T) {
	ars := &AtomicRunningStatus{}

	// Test setting to true
	ars.setRunningStatus(true)
	if !ars.runningStatus {
		t.Error("Expected runningStatus to be true")
	}

	// Test setting to false
	ars.setRunningStatus(false)
	if ars.runningStatus {
		t.Error("Expected runningStatus to be false")
	}
}

func TestAtomicRunningStatus_isRunning(t *testing.T) {
	ars := &AtomicRunningStatus{}

	// Test initial state (should be false)
	if ars.isRunning() {
		t.Error("Expected initial running status to be false")
	}

	// Test after setting to true
	ars.setRunningStatus(true)
	if !ars.isRunning() {
		t.Error("Expected running status to be true after setting to true")
	}

	// Test after setting to false
	ars.setRunningStatus(false)
	if ars.isRunning() {
		t.Error("Expected running status to be false after setting to false")
	}
}

func TestAtomicRunningStatus_ConcurrentAccess(t *testing.T) {
	ars := &AtomicRunningStatus{}
	var wg sync.WaitGroup

	// Number of goroutines to test concurrent access
	numGoroutines := 100
	numOperations := 1000

	// Test concurrent writes
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				// Alternate between true and false
				status := (id+j)%2 == 0
				ars.setRunningStatus(status)
			}
		}(i)
	}

	// Test concurrent reads
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				// Just read the status, don't care about the value
				// as long as it doesn't panic or race
				_ = ars.isRunning()
			}
		}()
	}

	// Wait for all goroutines to complete
	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	// Wait with timeout to avoid hanging tests
	select {
	case <-done:
		// Test passed, no race conditions detected
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out, possible deadlock or race condition")
	}
}

func TestGenerateUuid(t *testing.T) {
	uuid1 := generateUuid()
	uuid2 := generateUuid()

	// Test that UUID is not empty
	if uuid1 == "" {
		t.Error("Expected UUID to not be empty")
	}

	// Test that two UUIDs are different
	if uuid1 == uuid2 {
		t.Error("Expected two generated UUIDs to be different")
	}

	// Test UUID format (basic validation)
	// UUID should be in format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidPattern.MatchString(uuid1) {
		t.Errorf("Expected UUID to match standard format, got %s", uuid1)
	}

	// Test length (UUIDs should be 36 characters)
	if len(uuid1) != 36 {
		t.Errorf("Expected UUID length to be 36, got %d", len(uuid1))
	}

	// Test that multiple UUIDs are unique
	uuids := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		uuid := generateUuid()
		if uuids[uuid] {
			t.Errorf("Generated duplicate UUID: %s", uuid)
		}
		uuids[uuid] = true
	}
}

func TestGenerateRandomInt(t *testing.T) {
	min := 10
	max := 50

	// Test that random int is within range
	for i := 0; i < 100; i++ {
		randomInt := generateRandomInt(min, max)
		if randomInt < min || randomInt > max {
			t.Errorf("Expected random int to be between %d and %d, got %d", min, max, randomInt)
		}
	}

	// Test edge cases
	// When min == max
	sameValue := generateRandomInt(25, 25)
	if sameValue != 25 {
		t.Errorf("Expected random int to be 25 when min==max==25, got %d", sameValue)
	}

	// Test that it can generate both min and max values
	foundMin := false
	foundMax := false
	attempts := 10000

	for i := 0; i < attempts && (!foundMin || !foundMax); i++ {
		val := generateRandomInt(min, max)
		if val == min {
			foundMin = true
		}
		if val == max {
			foundMax = true
		}
	}

	if !foundMin {
		t.Errorf("Expected to find minimum value %d in %d attempts", min, attempts)
	}
	if !foundMax {
		t.Errorf("Expected to find maximum value %d in %d attempts", max, attempts)
	}
}

func TestGenerateRandomInt_Distribution(t *testing.T) {
	min := 1
	max := 10
	iterations := 10000
	counts := make(map[int]int)

	// Generate many random numbers and count occurrences
	for i := 0; i < iterations; i++ {
		val := generateRandomInt(min, max)
		counts[val]++
	}

	// Check that all values in range appeared at least once
	for i := min; i <= max; i++ {
		if counts[i] == 0 {
			t.Errorf("Value %d never appeared in %d iterations", i, iterations)
		}
	}

	// Check that no values outside range appeared
	for val := range counts {
		if val < min || val > max {
			t.Errorf("Value %d outside range [%d, %d] appeared", val, min, max)
		}
	}
}

func TestGenerateRandomString(t *testing.T) {
	str1 := generateRandomString()
	str2 := generateRandomString()

	// Test that string is not empty
	if str1 == "" {
		t.Error("Expected random string to not be empty")
	}

	// Test that two random strings are different (with high probability)
	if str1 == str2 {
		t.Error("Expected two generated random strings to be different")
	}

	// Test string format (should be XXXX-XXXX-XXXX-XXXX where X is digit)
	parts := strings.Split(str1, "-")
	if len(parts) != 4 {
		t.Errorf("Expected random string to have 4 parts separated by '-', got %d parts", len(parts))
	}

	// Each part should be 4 digits
	digitPattern := regexp.MustCompile(`^\d{4}$`)
	for i, part := range parts {
		if !digitPattern.MatchString(part) {
			t.Errorf("Expected part %d to be 4 digits, got '%s'", i, part)
		}
	}

	// Test that each part is within expected range (1000-9999)
	for i, part := range parts {
		// Convert to int and check range
		var num int
		if _, err := fmt.Sscanf(part, "%d", &num); err != nil {
			t.Errorf("Failed to parse part %d as integer: %s", i, part)
			continue
		}
		if num < 1000 || num > 9999 {
			t.Errorf("Expected part %d to be between 1000-9999, got %d", i, num)
		}
	}

	// Test uniqueness over multiple generations
	strings := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		str := generateRandomString()
		if strings[str] {
			t.Errorf("Generated duplicate random string: %s", str)
		}
		strings[str] = true
	}
}

func TestGenerateRandomString_Format(t *testing.T) {
	// Test that the format is consistent
	for i := 0; i < 100; i++ {
		str := generateRandomString()

		// Should be exactly 19 characters (4+1+4+1+4+1+4)
		if len(str) != 19 {
			t.Errorf("Expected random string length to be 19, got %d for string '%s'", len(str), str)
		}

		// Should match the pattern XXXX-XXXX-XXXX-XXXX
		pattern := regexp.MustCompile(`^\d{4}-\d{4}-\d{4}-\d{4}$`)
		if !pattern.MatchString(str) {
			t.Errorf("Expected random string to match pattern XXXX-XXXX-XXXX-XXXX, got '%s'", str)
		}
	}
}

// Benchmark tests
func BenchmarkAtomicRunningStatus_setRunningStatus(b *testing.B) {
	ars := &AtomicRunningStatus{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ars.setRunningStatus(i%2 == 0)
	}
}

func BenchmarkAtomicRunningStatus_isRunning(b *testing.B) {
	ars := &AtomicRunningStatus{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ars.isRunning()
	}
}

func BenchmarkGenerateUuid(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = generateUuid()
	}
}

func BenchmarkGenerateRandomInt(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = generateRandomInt(1000, 9999)
	}
}

func BenchmarkGenerateRandomString(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = generateRandomString()
	}
}
