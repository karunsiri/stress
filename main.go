package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Flag struct {
	CPUSeconds int
	Help       bool
	Memory     string
}

func parseFlags() Flag {
	flags := Flag{}

	flag.StringVar(&flags.Memory, "mem", "", "Amount of memory to allocate. Ex 300, 1K, 5G, 20M")
	flag.IntVar(&flags.CPUSeconds, "time", 0, "Duration of CPU workload in seconds")
	flag.BoolVar(&flags.Help, "help", false, "Show usage")
	flag.Parse()

	return flags
}

func parseBytes(str string) (int64, error) {
	str = strings.TrimSpace(str)
	suffixes := map[string]int64{
		"K": 1024,
		"M": 1024 * 1024,
		"G": 1024 * 1024 * 1024,
	}

	for suffix, multiplier := range suffixes {
		if strings.HasSuffix(str, suffix) {
			valueStr := strings.TrimSuffix(str, suffix)
			value, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse integer: %v", err)
			}
			return value * multiplier, nil
		}
	}

	// Return as bytes directly otherwise
	return strconv.ParseInt(str, 10, 64)
}

func allocateMemory(bytes int64) {
	// Allocate the memory
	memory := make([]byte, bytes)

	// Actually use the allocated memory, do not let OS cache/optimize
	for i := range memory {
		memory[i] = byte(i % 256)
	}
}

// Generate workload for all available CPU cores for the given `seconds`.
// Returns a `slice` containing number of iterations done by each goroutine (cpu
// cores) and goroutine WaitGroup.
func generateCPULoad(seconds int) (*[]int64, *sync.WaitGroup) {
	fmt.Printf("Performing CPU work for %d seconds\n", seconds)
	cpuCount := runtime.NumCPU()

	var tasks sync.WaitGroup
	tasks.Add(cpuCount)

	iterations := make([]int64, cpuCount)

	// Increasing this number will increase calculation complexity
	iterationsPerSecond := 100000
	secondsAsFloat := float64(seconds)

	for i := 0; i < cpuCount; i++ {
		go func(index int) {
			defer tasks.Done()

			start := time.Now()
			elapsed := time.Since(start).Seconds()

			// Ever increasing number to make math.Sqrt more complex
			problemFactor := 0.0

			for {
				if elapsed >= secondsAsFloat {
					break
				}

				for k := 0; k < iterationsPerSecond; k++ {
					// Pi * e
					_ = 3.14159265358979323846 * 2.71828182845904523536
					_ = 1.0
					_ = rand.Float64() * rand.Float64()
					_ = math.Sqrt(problemFactor)
					problemFactor += 1.0
				}
				iterations[index] += 1
				elapsed = time.Since(start).Seconds()
			}
		}(i)
	}

	return &iterations, &tasks
}

func calculateCPUScore(seconds int, iterations *[]int64) float64 {
	sum := float64(0)
	for _, value := range *iterations {
		sum += float64(value)
	}

	return sum / float64(seconds)
}

func main() {
	flags := parseFlags()

	var cpuTasks *sync.WaitGroup
	var cpuIterations *[]int64

	if flags.Help || flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	if flags.CPUSeconds > 0 {
		cpuIterations, cpuTasks = generateCPULoad(flags.CPUSeconds)
	}

	if flags.Memory != "" {
		bytes, err := parseBytes(flags.Memory)
		if err == nil {
			fmt.Printf("Allocating %s (%d bytes) of memory...\n", flags.Memory, bytes)
			allocateMemory(bytes)
		} else {
			fmt.Println("Error: ", err)
		}
	}

	if cpuTasks != nil {
		// Block untill all cpuTasks finish
		cpuTasks.Wait()
		fmt.Printf("CPU score: %.0f\n", calculateCPUScore(flags.CPUSeconds, cpuIterations))
	}

	fmt.Println("Workload generation complete.")
}
