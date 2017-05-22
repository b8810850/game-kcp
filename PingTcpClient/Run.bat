set /p num=please enter number of process that you want to start:
set /p time=please enter the time that you want to run process:(s)

for /l %%j in (1,1,%num%) do (
@echo %cd%
@echo %%j

for /l %%i in (1,1,%%j) do @start PingTcpClient.exe

TIMEOUT /T %time%


taskkill /f /im PingTcpClient.exe



cd PingStatistic
echo now *****************%cd%
PingStatistic.exe %%j

cd ..
	
echo %cd%
cd Log
echo %cd%

echo %cd%

del /s /q /f *.txt

for /d %%i in (*) do rd /s /q "%%i"

cd ..
)