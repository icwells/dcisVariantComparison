#!/bin/bash

##############################################################################
# This script will install scripts for the dcisVariantComparison package.
# 
# Required programs:	Go 1.11+
##############################################################################

IO="github.com/icwells/go-tools/iotools"
KP="gopkg.in/alecthomas/kingpin.v2"
MAIN="dcisVariantComparison"
SA="github.com/icwells/go-tools/strarray"

# Get install location
SYS=$(ls $GOPATH/pkg | head -1)
PDIR=$GOPATH/pkg/$SYS

installPackage () {
	# Installs go package if it is not present in src directory
	echo "Installing $1..."
	go get -u $1
	echo ""
}

installDependencies () {
# Get dependencies
	for I in $IO $KP $SA ; do
		if [ ! -e "$PDIR/$1.a" ]; then
			installPackage $I
		fi
	done
}

installMain () {
	# compOncDB 
	echo "Building $MAIN..."
	go build -i -o $MAIN src/*.go
	echo ""
}

echo ""
echo "Preparing dcisVariantComparison package..."
echo "GOPATH identified as $GOPATH"
echo ""

if [ $# -eq 0 ]; then
	installMain
elif [ $1 = "all" ]; then
	installDependencies
	installMain
elif [ $1 = "help" ]; then
	echo "Installs Go scripts for compOnDB"
	echo ""
	echo "all	Installs scripts and all depenencies."
	echo "test	Installs depenencies only (for white box testing)."
	echo "db	Installs scripts and dbIO only."
	echo "help	Prints help text and exits."
	echo ""
fi

echo "Finished"
echo ""
