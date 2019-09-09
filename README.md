# dcisVariantComparison is a Go program for comparing variants for the Maley lab DCIS project.  


Copyright 2019 by Shawn Rupp

1. [Dependencies](#Dependencies)  
2. [Installation](#Installation)  
3. [Usage](#Usage)  

## Dependencies:  

### Installing Go and Setting Paths  
[Go version 1.11 or higher](https://golang.org/doc/install)  

Go requires a GOPATH environment variable to set to install packages, an compOncDB requires the GOBIN variable to be set as well.  
Follow the directions [here](https://github.com/golang/go/wiki/SettingGOPATH) to set your GOPATH. Before you close your .bashrc or 
similar file, add the following lines after you deifne you GOPATH:  

	export GOBIN=$GOPATH/bin  
	export PATH=$PATH:$GOBIN  

## Installation  

	git clone https://github.com/icwells/dcisVariantComparison.git  
	cd dcisVariantComparison/
	./install.sh

## Usage  

