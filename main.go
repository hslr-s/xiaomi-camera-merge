package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	var path string
	var delete bool
	fmt.Println("小米摄像头视频合并工具启动成功")
	flag.StringVar(&path, "path", "./", "视频的保存目录")         // 定义字符串类型的标志
	flag.BoolVar(&delete, "delete", false, "删除已经合并的视频文件夹") // 定义字符串类型的标志
	flag.Parse()                                           // 解析标志

	fmt.Println("视频储存的主目录：", path)

	// DOCKER环境使用 读取环境变量是否删除删除已经合并的视频文件夹
	envDelete := os.Getenv("DELETE_SUCCESS")
	if envDelete == "true" {
		delete = true
	}
	fmt.Println("合并后是否删除源文件夹：", delete)
	MergeHour(path, delete)

}

func MergeHour(path string, deleteSrc bool) {
	// 指定要遍历的目录
	files, _ := os.ReadDir(path)

	for _, v := range files {
		if v.IsDir() && len(v.Name()) == 10 {
			fmt.Println(path + "/" + v.Name())
			workDir := path + "/" + v.Name()
			if MergeMp4ToMovByHour(workDir, path) && deleteSrc {
				if err := os.RemoveAll(workDir); err != nil {
					fmt.Println("删除源文件夹失败:", workDir, err.Error())
				} else {
					fmt.Println("已删除源文件夹:", workDir)
				}
			}
		}
	}
	fmt.Println("合并完成")
}

// 按小时合并视频
func MergeMp4ToMovByHour(dir, outputPath string) bool {
	fileContent := ""
	_, pathName := filepath.Split(dir)
	year := pathName[0:4]
	month := pathName[4:6]
	dateNum := pathName[6:8]
	// fmt.Println(pathName, year, month, dateNum, pathName[8:10])

	outputPath += "/" + year + "/" + month + "/" + dateNum

	// 创建目录
	os.MkdirAll(outputPath, 0777)

	// 使用 filepath.Walk 函数遍历目录
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// 判断文件是否为 .mp4 文件
		if filepath.Ext(path) == ".mp4" {
			// 打印文件路径
			// fmt.Println(path)
			_, fileName := filepath.Split(path)
			fileContent += "file " + fileName + "\n"
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	// 生成文件列表
	filesTxtPath := dir + "/files.txt"
	SaveFileList(filesTxtPath, fileContent)
	defer os.Remove(filesTxtPath)

	mergeFileName := pathName[8:10] + ".mov"
	// ffmpeg -f concat -i files.txt -c copy !name!.mov
	if ExecCommand(dir, "ffmpeg", "-f", "concat", "-i", "files.txt", "-c", "copy", mergeFileName) {
		// 移动文件
		os.Rename(dir+"/"+mergeFileName, outputPath+"/"+mergeFileName)
		fmt.Println("保存到文件", outputPath+"/"+mergeFileName)
		return true
	}
	return false
}

func SaveFileList(filePath, content string) {

	// 创建一个文件
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// 写入文件内容
	fmt.Fprintln(file, content)

	// 保存文件
	err = file.Sync()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ExecCommand(workpath, name string, arg ...string) bool {
	// 创建一个命令
	cmd := exec.Command(name, arg...)
	cmd.Dir = workpath
	fmt.Println("执行命令：", cmd)
	// 执行命令
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("指令执行错误：\n", err)
		return false
	}

	// 打印命令输出
	if string(out) != "" {
		fmt.Println("指令执行结果输出：\n", string(out))
	}
	return true
}
