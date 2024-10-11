package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// сюда писать код

const (
	MultiHashNum = 6
	Verbose      = true
	BuffSize     = 100
)

func ExecutePipeline(freeFlowJobs ...job) {
	out := make(chan interface{}, BuffSize)
	wg := &sync.WaitGroup{}
	for _, jb := range freeFlowJobs {
		in := out
		out = make(chan interface{}, BuffSize)
		wg.Add(1)
		go func(jb job) {
			defer wg.Done()
			jb(in, out)
		}(jb)
	}
	wg.Wait()
}

// global mutex
var Md5Mu = &sync.Mutex{}

func SingleHash(in, out chan interface{}) {
	data := <-in

	hashes := [2]string{}

	wg := &sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		hashes[0] = DataSignerCrc32(fmt.Sprint(data))
		if Verbose {
			fmt.Println(data, "SingleHash crc32(data)", hashes[0])
		}
	}()

	go func() {
		defer wg.Done()
		Md5Mu.Lock()
		md5hash := DataSignerMd5(fmt.Sprint(data))
		if Verbose {
			fmt.Println(data, "SingleHash md5(data)", md5hash)
		}
		Md5Mu.Unlock()
		hashes[1] = DataSignerCrc32(md5hash)
		if Verbose {
			fmt.Println(data, "SingleHash crc32(md5(data))", hashes[1])
		}
	}()
	wg.Wait()

	res := hashes[0] + "~" + hashes[1]
	if Verbose {
		fmt.Println(data, "SingleHash result", res)
	}
	out <- res
	close(out)
}

func MultiHash(in, out chan interface{}) {
	data := <-in

	hashes := [MultiHashNum]string{}

	wg := &sync.WaitGroup{}

	// calc hashes
	wg.Add(MultiHashNum)
	for th := 0; th < MultiHashNum; th++ {
		go func(th int) {
			defer wg.Done()
			hashes[th] = DataSignerCrc32(fmt.Sprint(th) + fmt.Sprint(data))
			if Verbose {
				fmt.Println(data, "MultiHash crc32(th+data)", th, hashes[th])
			}
		}(th)
	}
	wg.Wait()

	// hashes concatenation
	var sb strings.Builder
	for th := 0; th < MultiHashNum; th++ {
		sb.WriteString(hashes[th])
	}

	res := sb.String()
	if Verbose {
		fmt.Println(data, "MultiHash result:")
		fmt.Println(res)
	}

	out <- res
	close(out)
}

func CombineResults(in, out chan interface{}) {
	var results []string

	for res := range in {
		results = append(results, fmt.Sprint(res))
	}

	sort.Strings(results)

	var sb strings.Builder
	sb.WriteString(results[0])

	for i := 1; i < len(results); i++ {
		sb.WriteString("_")
		sb.WriteString(results[i])
	}

	out <- sb.String()
	close(out)
}
