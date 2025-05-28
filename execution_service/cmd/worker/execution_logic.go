package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func getMemoryUsage(pid int) int {
	statusPath := fmt.Sprintf("/proc/%d/status", pid)
	file, err := os.Open(statusPath)
	if err != nil {
		return -1
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "VmRSS:") {
			fields := strings.Fields(scanner.Text())
			memKB, _ := strconv.Atoi(fields[1])
			return memKB
		}
	}
	return -1
}

func ExecutePython(payload ExecuteCodePayload) (response ExecuteCodeResponse) {
	os.MkdirAll("./out", os.ModePerm)
	sourcePath := filepath.Join("out", "test.py")

	err := os.WriteFile(sourcePath, []byte(payload.Code), 0644)
	if err != nil {
		log.Println("Failed to write source file:", err)
		return ExecuteCodeResponse{
			ID:            payload.ID,
			Status:        "Error writing source file",
			ExecutionType: payload.ExecutionType,
		}
	}

	var results []TestCaseResult
	finalStatus := "Accepted"
	for _, tc := range payload.TestCases {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(payload.RuntimeLimitMS)*time.Millisecond)
		defer cancel()

		cmd := exec.CommandContext(ctx, "python3", sourcePath)
		cmd.Stdin = bytes.NewBufferString(tc.Input)

		var outBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &outBuf // Capture stderr too

		start := time.Now()
		if err := cmd.Start(); err != nil {
			results = append(results, TestCaseResult{
				ID:             tc.ID,
				Input:          tc.Input,
				Output:         "",
				ExpectedOutput: tc.ExpectedOutput,
				RuntimeMS:      0,
				MemoryKB:       0,
				Status:         "Error: " + err.Error(),
			})
			finalStatus = "Error"
			continue
		}

		pid := cmd.Process.Pid
		peakMem := getMemoryUsage(pid)

		done := make(chan error)
		go func() { done <- cmd.Wait() }()

		var status string
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()

	loop:
		for {
			select {
			case <-ticker.C:
				mem := getMemoryUsage(pid)
				if mem > peakMem {
					peakMem = mem
				}
				if payload.MemoryLimitKB > 0 && mem > payload.MemoryLimitKB {
					_ = cmd.Process.Kill()
					status = "MLE"
					break loop
				}
			case <-ctx.Done():
				_ = cmd.Process.Kill()
				status = "TLE"
				break loop
			case err := <-done:
				if err != nil {
					status = "Error"
				} else {
					status = "Accepted"
				}
				break loop
			}
		}

		end := time.Since(start)
		output := strings.TrimSpace(outBuf.String())
		expected := strings.TrimSpace(tc.ExpectedOutput)

		if status == "Accepted" && output != expected {
			status = "Wrong Answer"
		}

		if status != "Accepted" {
			finalStatus = status
		}

		results = append(results, TestCaseResult{
			ID:             tc.ID,
			Input:          tc.Input,
			Output:         output,
			ExpectedOutput: expected,
			RuntimeMS:      int(end.Milliseconds()),
			MemoryKB:       peakMem,
			Status:         status,
		})
	}

	response = ExecuteCodeResponse{
		ID:              payload.ID,
		Status:          finalStatus,
		ExecutionType:   payload.ExecutionType,
		TestCaseResults: results,
	}
	return response
}
