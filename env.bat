
@set GOROOT=D:\go
@set GOPATH=%cd%\vendor;C:\Users\dyh\go;%GOROOT%
@SETX GOPATH %GOPATH%
@set PATH=%PATH%;%GOROOT%\bin;%cd%\vender\bin;

@cmd.exe