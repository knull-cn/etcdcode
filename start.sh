#!/bin/sh
basedir=./bin
if [ ! -d "test" ];then
	mkdir test
fi 

function Copy(){
	cp $basedir/example test/dbmgr 
	cp $basedir/example test/example 
	cp $basedir/example test/gateway 
	cp $basedir/example test/logic 
	cp $basedir/example test/room 
	cp $basedir/example test/upload 
}

function Start(){
	Copy
	test/dbmgr &
	test/example &
	test/gateway &
	test/logic &
	test/room &
	test/upload &
}

case $1 in
	"start") 
		Start;;
	"copy") 
		Copy;;
	*) 
		Copy
		;;
esac
