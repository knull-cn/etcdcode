#!/bin/sh
cp example.exe test/dbmgr 
cp example.exe test/example 
cp example.exe test/gateway 
cp example.exe test/logic 
cp example.exe test/room 
cp example.exe test/upload 

test/dbmgr &
test/example &
test/gateway &
test/logic &
test/room &
test/upload &
