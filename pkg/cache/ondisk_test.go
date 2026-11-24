package cache

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestOnDiskSetGet(t *testing.T) {
	c := NewOnDiskCache(t.TempDir())

	if _, err := c.Get("missing"); err == nil {
		t.Fatal("expected an error for a missing key")
	}

	want := []byte(`{"type": "object"}`)
	if err := c.Set("key", want); err != nil {
		t.Fatalf("Set: %s", err)
	}
	got, err := c.Get("key")
	if err != nil {
		t.Fatalf("Get: %s", err)
	}
	if !bytes.Equal(got.([]byte), want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

// Set must replace an existing entry, so that a corrupt one (for example left
// behind by a process killed mid-write by an older version) can be repaired.
func TestOnDiskSetReplacesExistingEntry(t *testing.T) {
	c := NewOnDiskCache(t.TempDir())

	if err := c.Set("key", []byte(`{"type": "objec`)); err != nil {
		t.Fatalf("Set: %s", err)
	}
	want := []byte(`{"type": "object"}`)
	if err := c.Set("key", want); err != nil {
		t.Fatalf("Set: %s", err)
	}
	got, err := c.Get("key")
	if err != nil {
		t.Fatalf("Get: %s", err)
	}
	if !bytes.Equal(got.([]byte), want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

// Set must not leave temporary files behind in the cache folder.
func TestOnDiskSetLeavesNoTempFiles(t *testing.T) {
	dir := t.TempDir()
	c := NewOnDiskCache(dir)

	for i := 0; i < 3; i++ {
		if err := c.Set("key", []byte(`{}`)); err != nil {
			t.Fatalf("Set: %s", err)
		}
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		names := []string{}
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("expected exactly one file in the cache folder, got %v", names)
	}
	if entries[0].Name() != filepath.Base(cachePath(dir, "key")) {
		t.Errorf("unexpected file %s in the cache folder", entries[0].Name())
	}
}

// Several kubeconform processes may share one cache folder (for example
// parallel CI jobs passing the same -cache). Each process has its own onDisk,
// so the in-process mutex does not serialise them. A reader must see either no
// entry or the complete entry, never a partially written one.
//
// Separate onDisk instances stand in for separate processes here: they share
// the folder but not the mutex.
func TestOnDiskConcurrentProcessesNeverSeePartialEntries(t *testing.T) {
	dir := t.TempDir()
	writer := NewOnDiskCache(dir)
	reader := NewOnDiskCache(dir)

	// Large enough that writing it takes more than one write syscall.
	payload := bytes.Repeat([]byte("x"), 8*1024*1024)
	const keys = 20

	var wg sync.WaitGroup
	errs := make(chan error, keys)
	for i := 0; i < keys; i++ {
		key := fmt.Sprintf("schema-%d", i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				got, err := reader.Get(key)
				if err != nil {
					continue // not written yet
				}
				if n := len(got.([]byte)); n != len(payload) {
					errs <- fmt.Errorf("%s: read a partial entry of %d/%d bytes", key, n, len(payload))
				}
				return
			}
		}()
		if err := writer.Set(key, payload); err != nil {
			t.Fatalf("Set: %s", err)
		}
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
}
