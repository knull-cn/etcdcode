#!/bin/sh
basedir=./bin
cp $basedir/example test/dbmgr 
cp $basedir/example test/example 
cp $basedir/example test/gateway 
cp $basedir/example test/logic 
cp $basedir/example test/room 
cp $basedir/example test/upload 

test/dbmgr &
test/example &
test/gateway &
test/logic &
test/room &
test/upload &
