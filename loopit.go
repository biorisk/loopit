package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sync"
	"time"
)

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	cmd := flag.String("c", "", "Command to execute. Use MYFILE1 or MYFILE2 as wildcards")
	file1 := flag.String("f1", "", "First file to loop through")
	file2 := flag.String("f2", "", "Second file to loop through")
	pretend := flag.Bool("p", false, "Pretend. List commands that would have been run.")
	multi := flag.Int("n", 1, "number of simultaneous commands")
	flag.Parse()
	if *multi > numCPU {
		*multi = numCPU
	}
	list1, err := readLines(*file1)
	check(err)
	list2, err := readLines(*file2)
	check(err)
	re1 := regexp.MustCompile("MYFILE1")
	re2 := regexp.MustCompile("MYFILE2")
	cmdChan := make(chan string)
	var wg sync.WaitGroup
	if !*pretend {
		for i := 0; i < *multi; i++ {
			wg.Add(1)
			go executeCmd(cmdChan, &wg)
		}
	}
	start := time.Now()
	for _, f1 := range list1 {
		for _, f2 := range list2 {
			out := re1.ReplaceAllString(re2.ReplaceAllString(*cmd, f2), f1)
			if *pretend {
				fmt.Println(out)
			} else {
				cmdChan <- out
			}
		}
	}
	close(cmdChan)
	wg.Wait()
	fmt.Println("Finished at:", time.Now())
	fmt.Println("Elapsed time:", time.Since(start))
}

func executeCmd(cmdChan chan string, wg *sync.WaitGroup) {
	for cmd1 := range cmdChan {
		fmt.Println("starting command: ", cmd1)
		// args := strings.Split(cmd1, " ")
		cmd := exec.Command("sh", "-c", cmd1)
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
	wg.Done()
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
