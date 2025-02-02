package fileutil

import (
	"os"
	"path/filepath"
	"sync"
)

var currentDir string
var initOnce sync.Once

func init() {
	initOnce.Do(func() {
		currentDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	})
}

func FileWriteAll(filename string, data []byte) error {
	MakeSureDirExists(filename)
	if IsFileExist(filename) {
		DeleteFile(filename)
	}
	return os.WriteFile(filename, data, 0744)
}

// 获取当前运行目录
func GetCurrentDir() (ret string) {
	return currentDir
}

func IsFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func DeleteFile(filename string) {
	_ = os.Remove(filename)
}

func DeleteFileAll(path string) {
	_ = os.RemoveAll(path)
}

// 确保文件目录被创建，如果文件目录不存在则创建
func MakeSureDirExists(path string) {
	dir := filepath.Dir(path)
	_, err := os.Stat(dir)
	if err == nil {
		return
	}
	if os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0644)
	}
}

// 清空目录文件
func RemoveDirAll(dirname string) error {
	f, err := os.Open(dirname)
	if err != nil {
		return err
	}
	defer f.Close()
	names, err := f.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dirname, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// 判断路径是否是文件
func IsFile(path string) bool {
	fi, e := os.Stat(path)
	if e != nil {
		return false
	}
	return !fi.IsDir()
}
