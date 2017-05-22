set /p num=please enter number of process that you want to start:
set /p time=please enter the time that you want to run process:(s)
cd PingTcpClient
for /l %%i in (1,1,%num%) do @start PingTcpClient.exe

cd../PingUdpClient
for /l %%i in (1,1,%num%) do @start PingUdpClient.exe

TIMEOUT /T %time%

cd ..
@start ShutDown.bat