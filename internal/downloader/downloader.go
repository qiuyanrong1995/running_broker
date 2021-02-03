package downloader

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/juju/ratelimit"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// 下载程序到savePath
func Download(u string, tempPath string, savePath string, breakContinue bool, bandwidth string) error {
	bw, err := ParseBandwidth(bandwidth)
	if err != nil {
		return err
	}
	// 如果tempPath和savePath不存在则创建
	if !IsExist(tempPath) {
		if err := os.MkdirAll(tempPath, os.ModePerm); err != nil {
			return err
		}
	}
	if !IsExist(savePath) {
		if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
			return err
		}
	}

	var f *os.File
	tempFile := GetTempFile(u, tempPath)
	client := &http.Client{}
	request := http.Request{}
	request.Method = http.MethodGet
	if breakContinue {
		var size int64
		f, size, err = GetTempFileSeek(tempFile)
		if err != nil {
			return err
		}
		if size != 0 {
			header := http.Header{}
			header.Set("Range", "bytes=" + strconv.FormatInt(size, 10) + "-")
			request.Header = header
		}
	}

	parse, err := url.Parse(u)
	if err != nil {
		return err
	}
	request.URL = parse
	get, err := client.Do(&request)
	if err != nil {
		return err
	}
	defer func() {
		_ = get.Body.Close()
		_ = f.Close()
	}()

	if get.ContentLength == 0 {
		return nil
	}
	body := get.Body
	writer := bufio.NewWriter(f)
	bucket := ratelimit.NewBucket(time.Second, bw)
	_, err = io.Copy(writer, ratelimit.Reader(body, bucket))
	if err != nil {
		return err
	}
	return Save(tempFile, savePath)
}


func Save(fileName string, savePath string) error {
	nameArray := strings.Split(fileName, string(os.PathSeparator))
	realName := nameArray[len(nameArray) - 1]
	realName = realName[:len(realName) - 5]
	if strings.HasSuffix(savePath, string(os.PathSeparator)) {
		realName = savePath + realName
	} else {
		realName = savePath + string(os.PathSeparator) + realName
	}

	return os.Rename(fileName, realName)
}

func ParseBandwidth(bandwidth string) (int64, error) {
	prefix := strings.TrimSuffix(bandwidth, "/s")
	suffix := strings.ToLower(prefix[len(prefix) - 2:])
	bw, err := strconv.Atoi(prefix[:prefix[len(prefix) - 2]])
	if err != nil {
		return 0, err
	}
	switch suffix {
	case "kb":
		return int64(bw * 1024), nil
	case "mb":
		return int64(bw * 1024 * 1024), nil
	case "gb":
		return int64(bw * 1024 * 1024), nil
	default:
		return 0, errors.New(fmt.Sprintf("not support '%s'", suffix))
	}
}

func GetTempFileSeek(tempFile string) (*os.File, int64, error) {
	var err error
	var f *os.File
	var size int64
	if IsExist(tempFile) {
		f, err = os.OpenFile(tempFile, os.O_RDWR, os.ModePerm)
		if err != nil {
			return f, size,err
		}
		stat, err := f.Stat()
		if err != nil {
			return f, size, err
		}
		size = stat.Size()
		sk, err := f.Seek(size, 0)
		if sk != size {
			return f, size, errors.New("seek length not equal file size")
		}
	} else {
		f, err = os.Create(tempFile)
	}
	return f, size, err
}

func GetTempFile(u string, dir string) string {
	urlParts := strings.Split(u, string(filepath.Separator))
	file := urlParts[len(urlParts) - 1]
	if strings.HasSuffix(dir, string(os.PathSeparator)) {
		return dir + file + ".temp"
	}
	return dir + string(os.PathSeparator) + file + ".temp"
}