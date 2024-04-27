package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	doneExtension = ".done"
)

type conf struct {
	watchDir  string
	passwords []string

	delComplete bool
}

var gConf conf

func loadConf() {
	dir := os.Getenv("AUTOEX_DIR")
	if dir == "" {
		log.Fatal("AUTOEX_DIR environment variable is not set")
	}

	pwList := os.Getenv("AUTOEX_PW_LIST")
	if pwList == "" {
		log.Fatal("AUTOEX_PW_LIST environment variable is not set")
	}
	passwords := strings.Split(pwList, "|")

	del := os.Getenv("AUTOEX_DEL_COMPLETE")

	gConf = conf{
		watchDir:    dir,
		passwords:   passwords,
		delComplete: del == "true",
	}
}

func main() {
	loadConf()

	for {
		do()
		time.Sleep(10 * time.Second)
	}
}

func do() {
	err := filepath.Walk(gConf.watchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("file walk err: %v\n", err)
			return err
		}

		if strings.HasSuffix(info.Name(), ".done") {
			return nil
		}

		if (strings.HasSuffix(info.Name(), ".7z") || strings.HasSuffix(info.Name(), ".tar")) && !isDone(path) {
			log.Printf("scan %s\n", path)
			extract(path, gConf.passwords)
			return nil
		}

		if strings.Contains(info.Name(), ".7z.") && strings.HasSuffix(info.Name(), ".001") && !isDone(path) {
			log.Printf("scan %s\n", path)
			if !hasAllParts(path) {
				return nil
			}
			log.Printf("hasAllParts %s\n", path)
			extract(path, gConf.passwords)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	removeFlush()
}

func getAllParts(path string) []string {
	dir := filepath.Dir(path)
	baseName := filepath.Base(path)
	prefix := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	pattern := fmt.Sprintf("%s.*", prefix)
	files, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		log.Println("Failed to glob parts:", err)
		return []string{}
	}

	var files2 []string
	re := regexp.MustCompile(`^.*\.\d+$`)
	for _, f := range files {
		if re.MatchString(f) {
			files2 = append(files2, f)
		}
	}
	return files2
}

func hasAllParts(path string) bool {
	files := getAllParts(path)
	if len(files) == 0 {
		return false
	}

	sort.Slice(files, func(i, j int) bool {
		ni, _ := strconv.Atoi(filepath.Ext(files[i])[1:])
		nj, _ := strconv.Atoi(filepath.Ext(files[j])[1:])
		return ni < nj
	})

	var commonSize int64 = -1
	var lastSize int64 = -1
	var previousIndex int = -1

	for i, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			return false
		}

		currentIndex := extractIndex(file)
		if previousIndex != -1 && currentIndex != previousIndex+1 {
			return false
		}
		previousIndex = currentIndex

		if i == len(files)-1 {
			lastSize = info.Size()
		} else {
			if commonSize == -1 {
				commonSize = info.Size()
			} else if info.Size() != commonSize {
				return false
			}
		}
	}

	return lastSize != commonSize
}

func extractIndex(file string) int {
	ext := filepath.Ext(file)
	indexPart := ext[1:]

	index, err := strconv.Atoi(indexPart)
	if err != nil {
		log.Printf("Error extracting index from filename %s: %v\n", file, err)
		return -1
	}
	return index
}

func extract(path string, passwords []string) {
	for _, password := range passwords {
		var cmd *exec.Cmd
		if password == "" {
			cmd = exec.Command("7z", "x", path, "-o"+filepath.Dir(path)+"/out", "-y")
		} else {
			cmd = exec.Command("7z", "x", path, "-p"+password, "-o"+filepath.Dir(path)+"/out", "-y")
		}
		err := cmd.Run()
		log.Printf("Start extracted: %s with password: %s. err: %v. cmd: %s\n", path, password, err, cmd.String())
		if err == nil {
			createDoneFileOrDel(path)
			return
		}
	}
	createDoneFile(path)
}

func isDone(path string) bool {
	donePath := path + doneExtension
	_, err := os.Stat(donePath)
	return !os.IsNotExist(err)
}

func createDoneFile(path string) {
	donePath := path + doneExtension
	_, err := os.Create(donePath)
	if err != nil {
		log.Printf("Failed to create done file for %s: %v\n", path, err)
	}
}

func createDoneFileOrDel(path string) {
	if !gConf.delComplete {
		createDoneFile(path)
		return
	}

	if strings.HasSuffix(path, ".001") {
		files := getAllParts(path)
		for _, v := range files {
			removeDelay(v)
		}
	} else {
		removeDelay(path)
	}
}

var removePaths []string

func removeDelay(path string) {
	removePaths = append(removePaths, path)
}
func removeFlush() {
	for _, v := range removePaths {
		log.Println("remove file: ", v)
		_ = os.Remove(v)
	}
	removePaths = []string{}
}
