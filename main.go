package main

import (
	"errors"
	"flag"
	"fmt"
	"mi_camera_merge/lib/log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

const (
	VERSION     = "v1.1"
	AUTHOR      = "开发者：红烧猎人(联系作者:https://blog.enianteam.com/u/sun/content/11)"
	ADDRESS     = "工具开源地址和使用教程:https://github.com/hslr-s/xiaomi-camera-merge"
	ERROR_EMPTY = "empty folder"
)

var Logger *log.Log

func main() {
	var path string
	var delete bool
	var maxMerge string
	flag.StringVar(&path, "path", "./", "视频的保存目录")                                        // 定义字符串类型的标志
	flag.BoolVar(&delete, "delete", false, "删除已经合并的视频文件夹")                                // 定义字符串类型的标志
	flag.StringVar(&maxMerge, "max_merge", "0", "默认:0 不限制. 每次最大合并数量，如果视频量巨多，请使用本参数，分段执行") // 定义字符串类型的标志
	flag.Parse()                                                                          // 解析标志

	// DOCKER环境使用 读取环境变量是否删除删除已经合并的视频文件夹
	envDelete := os.Getenv("DELETE_SUCCESS")
	if envDelete == "true" {
		delete = true
	}

	envMaxMerge := os.Getenv("MAX_MERGE")
	if envDelete != "" {
		maxMerge = envMaxMerge
	}

	PrintAppInfo()
	fmt.Println("正在初始化，请稍后")
	fmt.Println("===========================")

	// 初始化日志
	currentTime := time.Now()
	formattedTime := currentTime.Format("20060102_150405")
	logPath := path + "/xiaomi-video-merge-log"
	if err := os.MkdirAll(logPath, 0777); err != nil {
		fmt.Println("错误中断执行：日志文件夹创建失败:", err.Error())
		return
	}
	Logger = log.NewLog(logPath + "/" + formattedTime + ".log")

	LoggerAppInfo()

	EchoLog("=== 本次运行配置 ===")
	EchoLog("视频储存的主目录：", path)
	EchoLog("合并后是否删除源文件夹：", delete)

	EchoLog("合并的时间取决于视频数量和尺寸，请尽量不要在合并过程中，修改操作目录的权限、删除目录等操作")
	EchoLog("=== 合并开始 ===")
	StartMergeHour(path, delete, maxMerge)
	EchoLog("=== 合并结束 ===")
	EchoLog("请注意！结束不代表全部成功了，具体请向上查看详情")
}

func StartMergeHour(path string, deleteSrc bool, maxMerge string) {
	maxMergeInt := 0
	i := 0
	if v, err := strconv.Atoi(maxMerge); err == nil {
		maxMergeInt = v
	}

	// 指定要遍历的目录
	files, _ := os.ReadDir(path)

	for _, v := range files {
		// 合并 有效路径
		if v.IsDir() && len(v.Name()) == 10 {
			workDir := path + "/" + v.Name()
			if savePath, err := MergeMp4ToMovByHour(workDir, path); err != nil {
				EchoLog("错误：目录：", workDir, "视频合并失败。错误原因：", err)
			} else {
				if deleteSrc {
					if err := os.RemoveAll(workDir); err != nil {
						EchoLog("目录:", workDir, "视频合并成功，已保存在：", savePath, "，未能成功删除源目录，请手动删除，否则下次运行程序，会再次合并")
					} else {
						EchoLog("目录:", workDir, "视频合并成功，已保存在：", savePath, "，已删除源目录")
					}
				} else {
					EchoLog("目录:", workDir, "视频合并成功，已保存在：", savePath)
				}
			}

			i++
			EchoLog("==========", i)

			// 最大执行数量
			if maxMergeInt != 0 && maxMergeInt == i {
				return
			}

		}
	}
}

// 按小时合并视频
func MergeMp4ToMovByHour(dir, outputPath string) (string, error) {
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
	i := 0 // 待合并文件数量统计
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// 判断文件是否为 .mp4 文件
		if filepath.Ext(path) == ".mp4" {
			// 打印文件路径
			// fmt.Println(path)
			_, fileName := filepath.Split(path)
			fileContent += "file " + fileName + "\n"
			i++
		}
		return nil
	})
	if err != nil {
		// 目录遍历失败，终止合并
		return "", errors.New("遍历目录失败：" + err.Error())
	}
	if i == 0 {
		return "", errors.New(ERROR_EMPTY)
	}

	// 生成文件列表
	filesTxtPath := dir + "/files.txt"
	if err := SaveFileList(filesTxtPath, fileContent); err != nil {
		return "", errors.New("生成临时文件出错" + err.Error())
	}
	defer os.Remove(filesTxtPath)

	mergeFileName := pathName[8:10] + ".mov"
	EchoLog("执行 ffmpeg 指令")
	// 执行指令 ffmpeg -f concat -i files.txt -c copy !name!.mov
	if err := ExecCommand(dir, "ffmpeg", "-f", "concat", "-i", "files.txt", "-c", "copy", mergeFileName); err != nil {
		return "", err
	}

	// 移动文件
	if err := os.Rename(dir+"/"+mergeFileName, outputPath+"/"+mergeFileName); err != nil {
		os.Remove(dir + "/" + mergeFileName)
		return "", errors.New("转移合并后的视频出错(为节省空间，已删除临时合并文件，记得重新合并)错误原因:" + err.Error())
	}
	// fmt.Println("保存到文件", outputPath+"/"+mergeFileName)
	return outputPath + "/" + mergeFileName, nil
}

func SaveFileList(filePath, content string) error {

	// 创建一个文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入文件内容
	fmt.Fprintln(file, content)

	// 保存文件
	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}

func ExecCommand(workpath, name string, arg ...string) error {
	// 创建一个命令
	cmd := exec.Command(name, arg...)
	cmd.Dir = workpath
	EchoLog("执行命令：", cmd)
	// 执行命令
	out, err := cmd.Output()
	if err != nil {
		EchoLog("指令执行错误：\n", err)
		return err
	}

	// 打印命令输出
	if string(out) != "" {
		EchoLog("指令执行结果输出：\n", string(out))
	}
	return nil
}

func PrintAppInfo() {
	fmt.Println("===========================")
	fmt.Println("小米摄像头视频合并工具", VERSION)
	fmt.Println("===========================")
	fmt.Println(AUTHOR)
	fmt.Println(ADDRESS)
	fmt.Println("===========================")
}

func LoggerAppInfo() {
	Logger.WriteContent("===========================")
	Logger.WriteContent("小米摄像头视频合并工具", VERSION)
	Logger.WriteContent("===========================")
	Logger.WriteContent(AUTHOR)
	Logger.WriteContent(ADDRESS)
	Logger.WriteContent("===========================")
}

func EchoLog(a ...any) {
	fmt.Println(a...)
	if err := Logger.WriteContent(a...); err != nil {
		fmt.Println("日志生成错误:", err.Error())
	}
}
