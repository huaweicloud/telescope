#!/bin/bash
CURRENTDIR=$(cd `dirname $0`;pwd)/..
PACKAGE_DIR="${CURRENTDIR}/packageTmp"
PACKAGE_DIST="${CURRENTDIR}/packageDistTmp"

OS=""
ARCH=""
AGENT_NAME=""
DAEMON_NAME=""
PACKAGE_NAME_PREFIX="telescope"

function assertExitCode()
{
    if [ $1 -ne 0 ]; then
        exit 1
    fi
}

function buildGoProject()
{
    export GOPATH=${WORKSPACE}
    
    OS="linux"
    ARCH_LIST="amd64 arm64"

    OS_LIST="linux windows"
    for i in $OS_LIST;
    do
        for j in $ARCH_LIST;
        do
        OS=$i
        ARCH=$j
        DAEMON_NAME="telescope"
        AGENT_NAME="agent"
        if [ ${OS} = "windows" ]; then
            ARCH="amd64"
            DAEMON_NAME="telescope.exe"
            AGENT_NAME="agent.exe"
        fi
        build
        assertExitCode $?
        package $1
        assertExitCode $?
        done
    done
}

function build()
{
	export GOOS=${OS}
	export GOARCH=${ARCH}
	
	cd $CURRENTDIR/agent
    go build -o ${AGENT_NAME}
    
    cd $CURRENTDIR/daemon
    go build -o ${DAEMON_NAME}
}

function package()
{
    PACKAGE_NAME=${PACKAGE_NAME_PREFIX}_${OS}_${ARCH}
    rm -rf ${PACKAGE_DIR}/${PACKAGE_NAME}
    mkdir ${PACKAGE_DIR}/${PACKAGE_NAME}
    cd ${PACKAGE_DIR}/${PACKAGE_NAME}
    mkdir bin
    
    cd $CURRENTDIR/agent
    mv ${AGENT_NAME} ${PACKAGE_DIR}/${PACKAGE_NAME}/bin
    cp conf.json ${PACKAGE_DIR}/${PACKAGE_NAME}/bin 
    cp conf_ces.json ${PACKAGE_DIR}/${PACKAGE_NAME}/bin
    cp conf_lts.json ${PACKAGE_DIR}/${PACKAGE_NAME}/bin
    cp record.json ${PACKAGE_DIR}/${PACKAGE_NAME}/bin
    if [ ${OS} = "windows" ]; then
        cp windows_os_log_record.json ${PACKAGE_DIR}/${PACKAGE_NAME}/bin
    fi
    cp logs_config.xml ${PACKAGE_DIR}/${PACKAGE_NAME}/bin

    cd $CURRENTDIR/daemon
    mv ${DAEMON_NAME} ${PACKAGE_DIR}/${PACKAGE_NAME}/bin
	
    if [ ${OS} = "linux" ]; then
        cd $CURRENTDIR/deploy/linux
        cp install.sh ${PACKAGE_DIR}/${PACKAGE_NAME}
        cp uninstall.sh ${PACKAGE_DIR}/${PACKAGE_NAME}
        cp telescoped ${PACKAGE_DIR}/${PACKAGE_NAME}
    else
        cd $CURRENTDIR/deploy/windows
        cp install.bat ${PACKAGE_DIR}/${PACKAGE_NAME}
        cp shutdown.bat ${PACKAGE_DIR}/${PACKAGE_NAME}
        cp start.bat ${PACKAGE_DIR}/${PACKAGE_NAME}
        cp uninstall.bat ${PACKAGE_DIR}/${PACKAGE_NAME}
    fi 
    
    cp ${CURRENTDIR}/ReleaseNote/telescope-${AGENT_VERSION}-release.json ${PACKAGE_DIR}/${PACKAGE_NAME}
    
    cd ${PACKAGE_DIR}
    tar -zcf ${PACKAGE_NAME}.tar.gz ${PACKAGE_NAME}
    tar -zcf ${PACKAGE_NAME}_${AGENT_VERSION}.tar.gz ${PACKAGE_NAME}
    md5sum ${PACKAGE_NAME}_${AGENT_VERSION}.tar.gz > ${PACKAGE_NAME}_${AGENT_VERSION}.tar.gz.md5
    zip -r ${PACKAGE_NAME}.zip ${PACKAGE_NAME}
    zip -r ${PACKAGE_NAME}_${AGENT_VERSION}.zip ${PACKAGE_NAME}
    md5sum ${PACKAGE_NAME}_${AGENT_VERSION}.zip > ${PACKAGE_NAME}_${AGENT_VERSION}.zip.md5
    echo "build package ${PACKAGE_NAME} success"
}

function generateSha()
{
    cd $PACKAGE_DIR

    for pkg in `ls **.tar.gz **.zip`
    do
        sha256sum $pkg > $pkg.sha256
    done
    
    zip -r telescope-$1-${AGENT_VERSION}.zip telescope_*.tar.gz telescope_*.tar.gz.sha256 telescope_*.zip telescope_*.zip.sha256 telescope_*.md5
    cp telescope-$1-${AGENT_VERSION}.zip  $PACKAGE_DIST
    assertExitCode $?
        
    sha256sum telescope-$1-${AGENT_VERSION}.zip > telescope-$1-${AGENT_VERSION}.zip.sha256
    cp telescope-$1-${AGENT_VERSION}.zip.sha256 $PACKAGE_DIST
    assertExitCode $?
    
    rm -rf telescope*
    
}

function packageAgentZip()
{
    buildGoProject $1
    generateSha $1
}

function main()
{
    releaseNote=${CURRENTDIR}/ReleaseNote/config.txt
    if [ ! -f ${releaseNote} ];then 
        echo "ERROR: config.txt file does not exist!"
        exit 1 
    else
        AGENT_VERSION=$(grep "releaseVersion=" ${releaseNote} | awk -F '=' {'print $2'}| tr -d '\n\r')
    fi

    cd ${CURRENTDIR}/ReleaseNote

    python GenerateHTML.py ftp://ftp/Agent
    python generate_release_json.py telescope telescope service
    
    #clear dir
    rm -rf $PACKAGE_DIR
    mkdir -p $PACKAGE_DIR
    rm -rf $PACKAGE_DIST
    mkdir -p $PACKAGE_DIST
    if [ ! -f ${CURRENTDIR}/ReleaseNote/*.json ]; then
        echo "ERROR: Json file does not exist!"
        exit 1
    fi
    cp ${CURRENTDIR}/ReleaseNote/*.json $PACKAGE_DIST
    cp ${CURRENTDIR}/ReleaseNote/COPYRIGHT.README $PACKAGE_DIST
    
    packageAgentZip cn-north-1
    
    rm -rf $PACKAGE_DIR
}

main
exit 0
