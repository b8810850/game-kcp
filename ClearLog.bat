@echo off

����cd /d "PingTcpClient/log"
	echo ȷ��Ҫɾ�������ļ��е������ļ�? %cd%
	pause
����del /s /q /f *.*

����for /d %%i in (*) do rd /s /q "%%i"

	cd /d "../../PingUdpClient/log"
	echo ȷ��Ҫɾ�������ļ��е������ļ�? %cd%
	pause
����del /s /q /f *.*

����for /d %%i in (*) do rd /s /q "%%i"

	echo ɾ���ɹ�
����pause