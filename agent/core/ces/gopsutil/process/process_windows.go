package process

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/process"
)

func GetWin32Proc(pids ...int32) ([]process.Win32_Process, error) {
	return GetWin32ProcWithContext(context.Background(), pids...)
}

// Test in PowerShell:  Get-WmiObject -Query "SELECT Name,ProcessId,CreationDate,CommandLine FROM Win32_Process"
func GetWin32ProcWithContext(ctx context.Context, pids ...int32) ([]process.Win32_Process, error) {
	var dst []process.Win32_Process
	query := "SELECT Name,ProcessId,CreationDate,CommandLine FROM Win32_Process"
	if len(pids) > 0 {
		query = " WHERE ProcessId = %d"
		for _, pid := range pids {
			query += strconv.Itoa(int(pid)) + " OR ProcessId = "
		}
		query = query[:len(query)-16]
	}
	err := WMIQueryWithContext(ctx, query, &dst)
	if err != nil {
		return []process.Win32_Process{}, fmt.Errorf("could not get win32Proc: %v", err)
	}

	if len(dst) == 0 {
		return []process.Win32_Process{}, fmt.Errorf("could not get win32Proc: empty")
	}

	return dst, nil
}

// WMIQueryWithContext - wraps wmi.Query with a timed-out context to avoid hanging
func WMIQueryWithContext(ctx context.Context, query string, dst interface{}, connectServerArgs ...interface{}) error {
	if _, ok := ctx.Deadline(); !ok {
		ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		ctx = ctxTimeout
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- wmi.Query(query, dst, connectServerArgs...)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}
