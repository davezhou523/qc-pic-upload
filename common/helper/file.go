package helper

import (
	"bufio"
	b64 "encoding/base64"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"os"
	"time"
)

func Base64Encode(data string) string {
	strEncode := b64.StdEncoding.EncodeToString([]byte(data))
	return strEncode
}

func Base64Decode(data string) string {
	byteStr, _ := b64.StdEncoding.DecodeString(data)
	s := string(byteStr)
	return s
}

/**
 * 文件转base64编码
 * @param unknown $attachment_path
 * @return string
 */
func FileToBase64(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logx.Errorf("%v 附件不存在:%s", filePath, err)
		return "", fmt.Errorf("%v 附件不存在:%s", filePath, err)
	}
	fileInfo, _ := file.Stat()
	defer file.Close()
	r := bufio.NewReader(file)
	buf := make([]byte, fileInfo.Size())
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("FileToBase64:%s", err)
		}
		if n == 0 {
			break
		}
	}
	return Base64Encode(string(buf)), nil
}

func Base64ToFile(filePath string, fileName string, base64Data string) error {
	fileInfo, _ := os.Stat(filePath)
	if fileInfo == nil {
		err := os.MkdirAll(filePath, 0777)
		if err != nil {
			msg := "创建文件夹失败"
			logx.Errorf("%v,%v:%v", filePath, msg, err)
			return fmt.Errorf("%v,%v:%v", filePath, msg, err)
		}
	}
	byteStr, _ := b64.StdEncoding.DecodeString(base64Data)
	filePtr, err := os.Create(filePath + fileName)
	fmt.Println(err)
	if err != nil {
		msg := "创建文件失败"
		logx.Errorf("%v %v:%v", filePath, msg, err)
		return fmt.Errorf("%v %v:%v", filePath, msg, err)
	}
	defer filePtr.Close()
	_, err = filePtr.Write(byteStr)
	if err != nil {
		msg := "创建文件保存内容失败"
		logx.Errorf("%v %v:%v", filePath, msg, err)
		return fmt.Errorf("%v %v:%v", filePath, msg, err)
	}
	return nil
}

func SaveFile(filePath string, fileName string, content string) error {
	fileInfo, _ := os.Stat(filePath)
	if fileInfo == nil {
		err := os.MkdirAll(filePath, 0777)
		if err != nil {
			msg := "创建文件夹失败"
			logx.Errorf("%v,%v:%v", filePath, msg, err)
			return fmt.Errorf("%v,%v:%v", filePath, msg, err)
		}
	}
	if len(fileName) == 0 {
		fileName = time.Now().Format("2006-01-02") + ".log"
	}
	filePtr, err := os.Create(filePath + "/" + fileName)
	if err != nil {
		logx.Errorf("%v 创建文件失败:%s", filePath, err)
		return fmt.Errorf("%v 创建文件失败:%s", filePath, err)
	}
	defer filePtr.Close()
	_, err = filePtr.Write([]byte(content))
	if err != nil {
		logx.Errorf("%v 创建文件保存内容失败:%s", filePath, err)
		return fmt.Errorf("%v 创建文件保存内容失败:%s", filePath, err)
	}
	return nil
}

func MoveFile(sourcePath string, destPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		msg := "源文件错误"
		logx.Errorf("%v:%v", msg, err)
		return fmt.Errorf("%v:%v", msg, err)
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		msg := "目标文件错误"
		logx.Errorf("%v:%v", msg, err)
		return fmt.Errorf("%v:%v", msg, err)
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, sourceFile)
	sourceFile.Close()
	if err != nil {
		msg := "复制文件错误"
		logx.Errorf("%v:%v", msg, err)
		return fmt.Errorf("%v:%v", msg, err)
	}
	fmt.Printf("移动文件成功,源文件：%v,目标文件:%v\n", sourceFile.Name(), destFile.Name())
	err = os.Remove(sourcePath)
	if err != nil {
		msg := "删除源文件错误"
		logx.Errorf("%v:%v", msg, err)
		return fmt.Errorf("%v:%v", msg, err)
	}
	return nil
}
