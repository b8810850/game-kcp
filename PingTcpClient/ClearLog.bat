@echo off

����cd /d "log"
	echo ȷ��Ҫɾ�������ļ��е������ļ�? %cd%
	pause
����del /s /q /f *.*

����for /d %%i in (*) do rd /s /q "%%i"
	echo ɾ���ɹ�
����pause