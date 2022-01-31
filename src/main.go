package main

import (
	"fmt"
	"io"
	"mohazit/lang/new"
	"mohazit/lib"
	"os"
)

const (
	eArgs int = 1 + iota
	eFile
	eRead
	eInterpreter
	eScript
	eCleanup
)

func main() {
	lib.Load()
	if len(os.Args) < 2 {
		fmt.Println("need input file")
		os.Exit(eArgs)
	} else {
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(eFile)
		}
		s, err := io.ReadAll(f)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(eRead)
		}
		new.Source(string(s))
		for {
			cont, err := new.DoAll()
			if !cont {
				break
			}
			if err != nil {
				if perr, ok := err.(*new.ParseError); ok {
					fmt.Printf("%s:%d:%d [ERROR] %s",
						os.Args[1], perr.Where.Line, perr.Where.Col, perr.Error())
				} else {
					fmt.Println(err.Error())
				}
				os.Exit(eScript)
			}
		}
	}
	if err := lib.Cleanup(); err != nil {
		fmt.Println("-- CLEANUP ERROR --")
		fmt.Println("(this usually isn't a serious problem, but should be avoided!")
		fmt.Println(err.Error())
		os.Exit(eCleanup)
	}
}
