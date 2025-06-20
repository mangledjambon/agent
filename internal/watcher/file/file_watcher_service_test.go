// Copyright (c) F5, Inc.
//
// This source code is licensed under the Apache License, Version 2.0 license found in the
// LICENSE file in the root directory of this source tree.

package file

import (
	"bytes"
	"context"
	"os"
	"path"
	"testing"
	"time"

	"github.com/nginx/agent/v3/test/helpers"
	"github.com/nginx/agent/v3/test/stub"

	"github.com/fsnotify/fsnotify"
	"github.com/nginx/agent/v3/test/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	directoryPermissions = 0o700
)

func TestFileWatcherService_NewFileWatcherService(t *testing.T) {
	fileWatcherService := NewFileWatcherService(types.AgentConfig())

	assert.Empty(t, fileWatcherService.directoriesBeingWatched)
	assert.True(t, fileWatcherService.enabled.Load())
	assert.False(t, fileWatcherService.filesChanged.Load())
}

func TestFileWatcherService_SetEnabled(t *testing.T) {
	fileWatcherService := NewFileWatcherService(types.AgentConfig())
	assert.True(t, fileWatcherService.enabled.Load())

	fileWatcherService.SetEnabled(false)
	assert.False(t, fileWatcherService.enabled.Load())

	fileWatcherService.SetEnabled(true)
	assert.True(t, fileWatcherService.enabled.Load())
}

func TestFileWatcherService_addWatcher(t *testing.T) {
	ctx := context.Background()
	fileWatcherService := NewFileWatcherService(types.AgentConfig())
	watcher, err := fsnotify.NewWatcher()
	require.NoError(t, err)
	fileWatcherService.watcher = watcher

	tempDir := os.TempDir()
	testDirectory := path.Join(tempDir, "test_dir")
	err = os.Mkdir(testDirectory, directoryPermissions)
	require.NoError(t, err)
	defer os.Remove(testDirectory)

	info, err := os.Stat(testDirectory)
	require.NoError(t, err)

	fileWatcherService.addWatcher(ctx, testDirectory, info)

	value, ok := fileWatcherService.directoriesBeingWatched.Load(testDirectory)
	assert.True(t, ok)
	boolValue, ok := value.(bool)
	assert.True(t, ok)
	assert.True(t, boolValue)
}

func TestFileWatcherService_addWatcher_Error(t *testing.T) {
	ctx := context.Background()
	fileWatcherService := NewFileWatcherService(types.AgentConfig())
	watcher, err := fsnotify.NewWatcher()
	require.NoError(t, err)
	fileWatcherService.watcher = watcher

	tempDir := os.TempDir()
	testDirectory := path.Join(tempDir, "test_dir")
	err = os.Mkdir(testDirectory, directoryPermissions)
	require.NoError(t, err)
	info, err := os.Stat(testDirectory)
	require.NoError(t, err)

	// Delete directory to cause the addWatcher function to fail
	err = os.Remove(testDirectory)
	require.NoError(t, err)

	fileWatcherService.addWatcher(ctx, testDirectory, info)

	value, ok := fileWatcherService.directoriesBeingWatched.Load(testDirectory)
	assert.True(t, ok)
	boolValue, ok := value.(bool)
	assert.True(t, ok)
	assert.False(t, boolValue)
	assert.True(t, ok)
}

func TestFileWatcherService_removeWatcher(t *testing.T) {
	ctx := context.Background()
	fileWatcherService := NewFileWatcherService(types.AgentConfig())
	watcher, err := fsnotify.NewWatcher()
	require.NoError(t, err)
	fileWatcherService.watcher = watcher

	tempDir := os.TempDir()
	testDirectory := path.Join(tempDir, "test_dir")
	err = os.Mkdir(testDirectory, directoryPermissions)
	require.NoError(t, err)
	defer os.Remove(testDirectory)

	err = fileWatcherService.watcher.Add(testDirectory)
	require.NoError(t, err)
	fileWatcherService.directoriesBeingWatched.Store(testDirectory, true)

	fileWatcherService.removeWatcher(ctx, testDirectory)

	value, ok := fileWatcherService.directoriesBeingWatched.Load(testDirectory)
	assert.Nil(t, value)
	assert.False(t, ok)

	logBuf := &bytes.Buffer{}
	stub.StubLoggerWith(logBuf)

	fileWatcherService.directoriesBeingWatched.Store(testDirectory, true)
	fileWatcherService.removeWatcher(ctx, testDirectory)

	helpers.ValidateLog(t, "Failed to remove file watcher", logBuf)

	logBuf.Reset()
}

func TestFileWatcherService_isEventSkippable(t *testing.T) {
	config := types.AgentConfig()
	config.Watchers.FileWatcher.ExcludeFiles = []string{"^/var/log/nginx/.*.log$", "\\.*swp$", "\\.*swx$", ".*~$"}
	fws := NewFileWatcherService(config)

	assert.False(t, fws.isEventSkippable(fsnotify.Event{Name: "test.conf"}))
	assert.True(t, fws.isEventSkippable(fsnotify.Event{Name: "test.swp"}))
	assert.True(t, fws.isEventSkippable(fsnotify.Event{Name: "test.swx"}))
	assert.True(t, fws.isEventSkippable(fsnotify.Event{Name: "test.conf~"}))
	assert.True(t, fws.isEventSkippable(fsnotify.Event{Name: "/var/log/nginx/access.log"}))
}

func TestFileWatcherService_isExcludedFile(t *testing.T) {
	excludeFiles := []string{"/var/log/nginx/access.log", "^.*(\\.log|.swx|~|.swp)$"}

	assert.True(t, isExcludedFile("/var/log/nginx/error.log", excludeFiles))
	assert.True(t, isExcludedFile("/var/log/nginx/error.swx", excludeFiles))
	assert.True(t, isExcludedFile("test.swp", excludeFiles))
	assert.True(t, isExcludedFile("/var/log/nginx/error~", excludeFiles))
	assert.True(t, isExcludedFile("/var/log/nginx/access.log", excludeFiles))
	assert.False(t, isExcludedFile("/etc/nginx/nginx.conf", excludeFiles))
	assert.False(t, isExcludedFile("/var/log/accesslog", excludeFiles))
}

func TestFileWatcherService_Watch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tempDir := os.TempDir()
	testDirectory := path.Join(tempDir, "test_dir")
	os.RemoveAll(testDirectory)
	err := os.Mkdir(testDirectory, directoryPermissions)
	require.NoError(t, err)
	defer os.RemoveAll(testDirectory)

	agentConfig := types.AgentConfig()
	agentConfig.Watchers.FileWatcher.MonitoringFrequency = 100 * time.Millisecond
	agentConfig.AllowedDirectories = []string{testDirectory, "/unknown/directory"}

	channel := make(chan FileUpdateMessage)

	fileWatcherService := NewFileWatcherService(agentConfig)
	go fileWatcherService.Watch(ctx, channel)

	time.Sleep(100 * time.Millisecond)

	file, err := os.CreateTemp(testDirectory, "test.conf")
	require.NoError(t, err)
	defer os.Remove(file.Name())

	t.Run("Test 1: File updated", func(t *testing.T) {
		fileUpdate := <-channel
		assert.NotNil(t, fileUpdate.CorrelationID)
	})

	t.Run("Test 2: Skippable file updated", func(t *testing.T) {
		skippableFile, skippableFileError := os.CreateTemp(testDirectory, "*test.conf.swp")
		require.NoError(t, skippableFileError)
		defer os.Remove(skippableFile.Name())

		select {
		case <-channel:
			t.Fatalf("Expected file to be skipped")
		case <-time.After(150 * time.Millisecond):
			return
		}
	})
}
