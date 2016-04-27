package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Job struct {
	Start time.Time
	Cmd   string
}

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
	list2 := []string{"null"}
	if *file2 != "" {
		list2, err = readLines(*file2)
		check(err)
	}
	re1 := regexp.MustCompile("MYFILE1")
	re2 := regexp.MustCompile("MYFILE2")
	job := regexp.MustCompile("MYJOB")
	cmdChan := make(chan Job)
	var wg sync.WaitGroup
	if !*pretend {
		for i := 0; i < *multi; i++ {
			wg.Add(1)
			go executeCmd(cmdChan, &wg)
		}
	}
	start := time.Now()
	jobsStarted := 0
	jobsTotal := len(list1) * len(list2)
	for _, f1 := range list1 {
		for _, f2 := range list2 {
			jobsStarted++
			jobNum := strconv.Itoa(jobsStarted)
			out := job.ReplaceAllString(re1.ReplaceAllString(re2.ReplaceAllString(*cmd, f2), f1), jobNum)
			if *pretend {
				fmt.Println(out)
			} else {
				cmd := Job{Start: time.Now(), Cmd: out}
				cmdChan <- cmd
				fmt.Printf("Job %d of %d started\n", jobsStarted, jobsTotal)
				fmt.Println("Time elapsed since start ", time.Since(start))
			}
		}
	}
	close(cmdChan)
	wg.Wait()
	fmt.Println("Finished at:", time.Now())
	fmt.Println("Elapsed time:", time.Since(start))
}

func executeCmd(cmdChan chan Job, wg *sync.WaitGroup) {
	for cmd1 := range cmdChan {
		// fmt.Println("starting command: ", cmd1)
		// args := strings.Split(cmd1, " ")
		cmd := exec.Command("sh", "-c", cmd1.Cmd) //use the shell to interpret the command
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
