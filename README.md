# 小米摄像头视频文件合并

> 由于小米摄像头保存的视频传到NAS都是一分钟一个文件，而且NAS上保存的视频不支持在APP回看，所以回看相当的不方便。所以百度查找到一个合并视频文件的方式，根据[大神的开发的windows脚本](./小米摄像头录像合并脚本.bat)开发成golang版本，跨平台使用，支持docker。

## 项目开源地址及最新教程
项目地址：https://github.com/hslr-s/xiaomi-camera-merge

博客地址：https://blog.enianteam.com/u/sun/content/239

哔哩哔哩视频：https://www.bilibili.com/video/BV1rN4y1k7K3/

## 软件下载
天翼云盘：https://cloud.189.cn/web/share?code=raANbiiA3QRb （访问码：qqy4）

Github：https://github.com/hslr-s/xiaomi-camera-merge/releases

## 实现原理和说明
go程序通过代码，调用开源`ffmpeg`（用于生成处理多媒体数据的各类库和程序）程序来实现视频的合并转存。将一分钟一个视频合并为以小时为一个视频进行保存按天保存

**❗温馨提示：已知部分版本、或新购买的部分摄像头不支持，首次使用，请先拷贝出部分目录测试一下。技术有价,数据无价。若造成数据丢失，作者不负责哦❗**
**❗温馨提示：已知部分版本、或新购买的部分摄像头不支持，首次使用，请先拷贝出部分目录测试一下。技术有价,数据无价。若造成数据丢失，作者不负责哦❗**
**❗温馨提示：已知部分版本、或新购买的部分摄像头不支持，首次使用，请先拷贝出部分目录测试一下。技术有价,数据无价。若造成数据丢失，作者不负责哦❗**


程序会在视频文件夹中生成以年为名的的文件夹和日志文件夹（xiaomi-video-merge-log），每次执行的日志会以日期命名的.log文件。可以使用记事本方式打开查看记录。


视频文件合并完成后是按每小时一个视频储存的`.mov`格式的视频文件

合并前后的对比图

![](./doc/screenshot20231024.png)

## 编译运行

#### 1. 搭建ffmpeg环境（如果你有docker环境，请直接向下看docker运行）
根据自身系统查找教程即可

#### 2. 编译（确保已经搭建完成go开发环境）

```
# linux
go run -o xiaomi_camera_merge main.go

# windows
go run -o xiaomi_camera_merge.exe main.go
```


#### 运行

参数说明（windows双击执行参数将都是默认参数）：

- -delete （默认：不删除）合并完成后删除源文件夹及视频（如果不删除每次运行都会把目录所有的视频再合并一次）
- -path string  （默认:当前所在目录）视频的保存目录 (default "./")
- -max_merge （默认:0，不限制）每次最大合并数量，如果视频量巨多，请使用本参数，分段执行 (default "0")


命令行执行示例：
```
# linux 
./xiaomi_camera_merge -delete -path /mnt/xiaomi/video

# windows
.\xiaomi_camera_merge.exe -delete -path D:\xiaomi_video

```

## docker运行（推荐）


运行参数说明
- `/mnt/xiaomi/camera` 换成你摄像机储存视频的目录即可
- `-e "DELETE_SUCCESS=true"` [可选]合并完成后删除源文件夹及视频（如果不删除每次运行都会把目录所有的视频再合并一次）
- `-e "MAX_MERGE=100"` [可选]默认:0 不限制. 每次最大合并数量，如果视频量巨多，请使用本参数，分段执行 (default "0")

方式1.使用构建好的镜像：
```
# 拉取(最新)镜像
docker pull hslr/xiaomi-camera-merge:latest
# 运行
docker run --rm -it -v /mnt/xiaomi/camera:/app/video -e "DELETE_SUCCESS=true" hslr/xiaomi-camera-merge:latest
```

方式2.自构建镜像:
```
# 1.拉取本仓库代码后，进入项目跟目录运行下方代码
docker build -t xiaomi-camera-merge:latest .

# 2.运行容器
docker run --rm -it -v /mnt/xiaomi/camera:/app/video -e "DELETE_SUCCESS=true" xiaomi-camera-merge
```

