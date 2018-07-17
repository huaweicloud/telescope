# One Cloud, One Agent.
Telescope is an agent of Huawei clouds. Telescope collect the servers' metrics info and send to cloudeye, users can monitor on console.  
  
# Useful links
*   [Reference documentation](https://godoc.org/github.com/huaweicloud/golangsdk)  
*   [Effective Go](https://golang.org/doc/effective_go.html)  
  
# Requirements
* go 1.9+  
  
# How to get telescope
Before getting, you need to ensure that your GOPATH environment variable is pointing to an appropriate directory where you want to install Telescope  
  
go get github.com/huaweicloud/telescope  
  
# How to build the source code 
You can build it into binary tar/zip and install it as a Windows and Linux Service. Package command:  
. cd $HOME/go/src/github.com/huaweicloud/telescope/CI  
. export WORKSPACE=${YOUR GOPATH}  
. sh package.sh  
  
 Note that you can build and package on linux or windows, but windows maybe need other dependency packages like bashzip.
