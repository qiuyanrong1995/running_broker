package windows

import (
	"bytes"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os/exec"
	"strings"
)

type Charset string

const (
	UTF8	= Charset("UTF-8")
	GB18030	= Charset("GB18030")

)

func ConvertByte2String(b []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(b)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(b)
	}
	return str
}

// 返回进程名列表
func ListProcesses() ([]string, error) {
	var pMap = make(map[string]struct{})
	var pNames []string
	var err error
	processes, _ := process.Processes()
	for _, p := range processes {
		pName, err := p.Name()
		if err != nil {
			return pNames, err
		}
		pMap[pName] = struct{}{}
	}

	for name, _ := range pMap {
		pNames = append(pNames, name)
	}
	return pNames, err
}

// 运行系统命令
func Exec(command string) (string, error) {
	cArray := strings.Split(command, " ")
	c := exec.Command(cArray[0], cArray[1:]...)
	buf := bytes.NewBuffer([]byte{})
	c.Stdout = buf
	err := c.Run()
	return ConvertByte2String(buf.Bytes(), GB18030), err
}

func ListInstallPrograms() error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, "")
}