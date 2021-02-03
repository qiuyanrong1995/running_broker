//+build windows

package running_broker


import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"log"
)

func main() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`, registry.ALL_ACCESS)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	s, _, err := k.GetStringValue("test")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q\n", s)
}