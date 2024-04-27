package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractIndex(t *testing.T) {
	testCases := []struct {
		filename string
		expected int
	}{
		{"archive.7z.001", 1},
		{"archive.7z.002", 2},
		{"archive.7z.010", 10},
		{"archive.7z", -1},
	}

	for _, tc := range testCases {
		result := extractIndex(tc.filename)
		if result != tc.expected {
			t.Errorf("extractIndex(%s) = %d; want %d", tc.filename, result, tc.expected)
		}
	}
}

func TestHasAllParts(t *testing.T) {
	testDir := t.TempDir()

	// 创建分卷文件，前两个大小一致，最后一个不一致
	createTestFile(t, filepath.Join(testDir, "archive.7z.001"), 1024)
	createTestFile(t, filepath.Join(testDir, "archive.7z.002"), 1024)
	createTestFile(t, filepath.Join(testDir, "archive.7z.003"), 512) // 最后一个文件大小不同

	result := hasAllParts(filepath.Join(testDir, "archive.7z.001"))
	if !result {
		t.Errorf("hasAllParts() should be true when last file size is different but sequence is correct")
	}

	os.Remove(filepath.Join(testDir, "archive.7z.002"))
	result = hasAllParts(filepath.Join(testDir, "archive.7z.001"))
	if result {
		t.Errorf("hasAllParts() should be false when file sequence is not continuous")
	}
}

func TestHasAllParts2(t *testing.T) {
	testDir := t.TempDir()

	createTestFile(t, filepath.Join(testDir, "archive.7z.001"), 1024)
	createTestFile(t, filepath.Join(testDir, "archive.7z.002"), 1024)

	result := hasAllParts(filepath.Join(testDir, "archive.7z.001"))
	if result {
		t.Errorf("hasAllParts() should be true when last file size is different but sequence is correct")
	}
}

func createTestFile(t *testing.T, path string, size int64) {
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.Write(make([]byte, size)); err != nil {
		t.Fatal(err)
	}
	file.Close()
}
