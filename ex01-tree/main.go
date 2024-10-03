package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"sort"
	"strings"
	// "path/filepath"
	// "strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

type Stack struct {
	Items []string
}

func (s *Stack) Push(item string) {
	s.Items = append(s.Items, item)
}

func (s *Stack) Pop() string {
	if len(s.Items) == 0 {
		return ""
	}

	lastIdx := len(s.Items) - 1
	item := s.Items[lastIdx]
	s.Items = s.Items[:lastIdx]
	return item
}

func (s *Stack) ReplaceLast(item string) {
	s.Pop()
	s.Push(item)
}

func (s Stack) String() string {
	return strings.Join(s.Items, "")
}

func filter(ss []fs.DirEntry, test func(fs.DirEntry) bool) (ret []fs.DirEntry) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func bytelnToHumanReadableSize(byteln int64) string {
	if byteln == 0 {
		return "empty"
	}

	return fmt.Sprint(byteln, "b")
}

func entryToString(de fs.DirEntry) (string, error) {
	var sb strings.Builder
	sb.WriteString(de.Name())
	if !de.IsDir() {
		fi, err := de.Info()
		if err != nil {
			return "", err
		}
		sb.WriteString(fmt.Sprint(" (", bytelnToHumanReadableSize(fi.Size()), ")"))
	}
	sb.WriteString("\n")
	return sb.String(), nil
}

func writeDirTree(out io.Writer, path string, printFiles bool, bars Stack) error {
	dir, err := os.Open(path)
	if err != nil {
		return err
	}

	entries, err := dir.ReadDir(0)
	if err != nil {
		return err
	}

	if !printFiles {
		entries = filter(entries, func(de fs.DirEntry) bool { return de.IsDir() })
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	lastIdx := len(entries) - 1
	if lastIdx < 0 {
		return nil
	}

	bars.Push("├───")

	for i := range entries[:lastIdx] {
		entryStr, err := entryToString(entries[i])
		if err != nil {
			return err
		}
		io.WriteString(out, fmt.Sprint(bars, entryStr))
		if entries[i].IsDir() {
			bars.ReplaceLast("│	")
			writeDirTree(out, path+"/"+entries[i].Name(), printFiles, bars)
			bars.ReplaceLast("├───")
		}
	}

	bars.ReplaceLast("└───")
	entryStr, err := entryToString(entries[lastIdx])
	if err != nil {
		return err
	}
	io.WriteString(out, fmt.Sprint(bars, entryStr))
	bars.ReplaceLast("	")

	if entries[lastIdx].IsDir() {
		writeDirTree(out, path+string(os.PathSeparator)+entries[lastIdx].Name(), printFiles, bars)
	}

	return nil
}

// func writeDirTree(out *os.File, path string, printFiles bool, bars Stack) error {
// 	dir, err := os.Open(path)
// 	if err != nil {
// 		return err
// 	}

// 	entries, err := dir.ReadDir(0)
// 	if err != nil {
// 		return err
// 	}

// 	sort.Slice(entries, func(i, j int) bool {
// 		return entries[i].Name() < entries[j].Name()
// 	})

// 	lastIdx := len(entries) - 1

// 	if lastIdx > 0 {
// 		bars.Push("├───")
// 	}

// 	for i := range entries[:lastIdx] {
// 		io.WriteString(out, fmt.Sprint(bars, entries[i].Name(), "\n"))
// 		if entries[i].IsDir() {
// 			bars.Pop()
// 			bars.Push("│   ")
// 			writeDirTree(out, path+"/"+entries[i].Name(), printFiles, bars)
// 			bars.Pop()
// 		}
// 	}

// 	bars.Pop()
// 	bars.Push("└───")
// 	io.WriteString(out, fmt.Sprint(bars, entries[lastIdx].Name(), "\n"))
// 	bars.Pop()
// 	bars.Push("    ")
// 	if entries[lastIdx].IsDir() {
// 		writeDirTree(out, path+"/"+entries[lastIdx].Name(), printFiles, bars)
// 	}
// 	bars.Pop()

// 	return nil
// }

func dirTree(out io.Writer, path string, printFiles bool) error {
	var bars Stack
	return writeDirTree(out, path, printFiles, bars)
}
