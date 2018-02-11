@echo off
::set linux env;
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
::
set input=%1
set pwd=%~dp0
set bindir=%pwd%\bin\
set GOPATH=F:/github/etcdcode

if [%1] == [] goto:all

if %1==clean (
call:clean
) else if %1==proto (
call:proto %2
) else (
call:build %1
)
goto:exit

:all
if not "%input%"=="" (
call:build %input%
exit /b 0
)

call:build etcdop
call:build example


exit /b 0

:build
go build -o %bindir%/%1  -i %1 
if %errorlevel%==0 (
echo build LINUX success : %1 
) else (
echo " "
echo ****ERROR****  : %1 
echo " "
)
exit /b 0

:proto

exit /b 0

:clean
rm pkg/* -rf
rm bin/*.exe -rf
echo clean ok!
exit /b 0

:exit
