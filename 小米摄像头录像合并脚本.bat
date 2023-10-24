@echo off & setlocal EnableDelayedExpansion
for /f "delims=" %%a in ('dir /ad/b') do (
set name=%%a
set name=%cd%\!name:~,8!_!name:~-2!
set var=%cd%\%%a
set var=!var:\\=\!
echo !var!
cd "!var!"
for /f %%s in ('dir /b "*.mp4"') do ( 
echo file %%s >> files.txt
)
set /p ms=<"files.txt"
set name=!name!!ms:~5,2!!ms:~8,2!
ffmpeg -f concat -i files.txt -c copy !name!.mov
del files.txt
echo !name!.mov 已生成。
cd ..
)
TIMEOUT /T 600
for /f "delims=" %%z in ('dir /b *.mov') do (
set b=%%z
set c=!b:~,4!
set d=!b:~4,2!
set e=!b:~6,2!

echo %%z" "!c!\!d!\!e!

if not exist "%cd%\!c!\!d!\!e!" md "%cd%\!c!\!d!\!e!"
move "%%z" "!c!\!d!\!e!"
echo %%z 成功移动至 !c!\!d!\!e! 文件夹！
)
pause

代码来自：http://maryd.cn/