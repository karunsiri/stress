package main

import (
	"flag"
	"fmt"
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

// Generate workload for all available CPU cores for the given `seconds`
func generateCPULoad(seconds int) *sync.WaitGroup {
	fmt.Printf("Performing CPU work for %d seconds\n", seconds)
	cpuCount := runtime.NumCPU()

	var tasks sync.WaitGroup
	tasks.Add(cpuCount)

	for i := 0; i < cpuCount; i++ {
		go func() {
			defer tasks.Done()
			start := time.Now()
			for time.Since(start).Seconds() < float64(seconds) {
				_ = rand.Float64() * rand.Float64()
			}
		}()
	}

	return &tasks
}

func main() {
	flags := parseFlags()

	var cpuTasks *sync.WaitGroup

	if flags.Help || flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	if flags.CPUSeconds > 0 {
		cpuTasks = generateCPULoad(flags.CPUSeconds)
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
	}

	fmt.Println("Workload generation complete.")
}
