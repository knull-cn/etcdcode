@echo off
set pwd=%~dp0
set binPath=%pwd%/tools
set srcpath=%pwd%/src/myrpc/protocol
echo 'proto'
cd %srcpath%
%binPath%\protoc3.2 --go_out=plugins=grpc:. *.proto
cd %pwd%

echo 'finished'
::pause