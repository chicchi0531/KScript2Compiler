package ast

import(
	"ks2/compiler/vm"
	"os"
	"fmt"
)

type DumpStatement struct{
	outFileName string
	lineno int
	driver *vm.Driver
}

func MakeDumpStatement(out string, lineno int, driver *vm.Driver) *DumpStatement{
	s := new(DumpStatement)
	s.outFileName = out
	s.lineno = lineno
	s.driver = driver
	return s
}

func (s *DumpStatement) Analyze() {
	if s.outFileName == ""{
		s.driver.Dump(os.Stdout)
	} else {
		dumpFile, err := os.Create(s.outFileName)
		if err != nil{
			fmt.Println("dumpファイルを開けません。")
		}
		defer dumpFile.Close()
		s.driver.Dump(dumpFile)
	}
}