//+build windows

package main

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
)

func main() {
	//k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`, registry.ALL_ACCESS)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer k.Close()
	//
	//s, _, err := k.GetStringValue("test")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("%q\n", s)

	//var rootProcess *process.Process
	processes, _ := process.Processes()
	for _, p := range processes {
		fmt.Println(p.String())
	}

}