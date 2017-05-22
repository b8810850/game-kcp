@echo off

　　cd /d "PingTcpClient/log"
	echo 确定要删除以下文件夹的所有文件? %cd%
	pause
　　del /s /q /f *.*

　　for /d %%i in (*) do rd /s /q "%%i"

	cd /d "../../PingUdpClient/log"
	echo 确定要删除以下文件夹的所有文件? %cd%
	pause
　　del /s /q /f *.*

　　for /d %%i in (*) do rd /s /q "%%i"

	echo 删除成功
　　pause