#!/bin/bash
OS_DESC=''
OTHERS_OS='OTHERS'
CURRENT_OS=''

support_os_list=("CENTOS" "SUSE" "DEBIAN" "EULER" "UBUNTU" "FEDORA" "ORACLE LINUX")
chkconfig_os_list=("CENTOS" "SUSE" "EULER" "FEDORA" "ORACLE LINUX")
update_rc_os_list=("DEBIAN" "UBUNTU")

if [ "`id -u`" = "0" ] || [ "`id -g`" = "0" ] ; then
    echo "Current user is root."
else
    echo "Current user is not root, please use root user uninstall or command [sudo sh uninstall.sh]."
    exit 0
fi

getCurrentPath()
{
    if [ "` dirname "$0" `" = "" ] || [ "` dirname "$0" `" = "." ] ; then
        CURRENT_DIR="`pwd`"
    else
        cd ` dirname "$0" `
        CURRENT_DIR="`pwd`"
        cd - > /dev/null 2>&1
    fi
}

#get linux os version description 
getOS()
{
    if [ -f /usr/bin/lsb_release ]; then
        OS_DESC=$(/usr/bin/lsb_release -a |grep Description |awk -F : '{print $2}' |sed 's/^[ \t]*//g')
    elif [ -f /etc/system-release ]; then
        OS_DESC=$(cat /etc/system-release | sed -n '1p')
    else
        OS_DESC=$(cat /etc/issue | sed -n '1p')
    fi
}


getCurrentPath
getOS

# get current linux os version
CURRENT_OS=${OTHERS_OS}
for support_os in "${support_os_list[@]}"
do 
    if [ `echo ${OS_DESC} | tr [a-z] [A-Z] | grep "${support_os}" | wc -l` -ge 1 ] ; then
        CURRENT_OS=${support_os}
    fi
done

echo "Current linux release version : ${CURRENT_OS}"
echo "Start to uninstall telescope..."
#stop telescoped service
if command -v service >/dev/null 2>&1; then
    service telescoped stop
elif command -v rc-service >/dev/null 2>&1; then
    rc-service telescoped stop
elif command -v systemctl >/dev/null 2>&1; then
    systemctl stop telescoped
fi

#remove telescoped service
if [[ "${chkconfig_os_list[@]}" =~ $CURRENT_OS ]]; then
    chkconfig telescoped off
    chkconfig --del telescoped
    rm -f /etc/init.d/telescoped
elif [[ "${update_rc_os_list[@]}" =~ $CURRENT_OS ]]; then
    update-rc.d -f telescoped remove
    rm -f /etc/init.d/telescoped
else
    if command -v chkconfig >/dev/null 2>&1; then 
        chkconfig telescoped off
        chkconfig --del telescoped
    elif command -v update-rc.d >/dev/null 2>&1; then 
        update-rc.d -f telescoped remove
    elif command -v rc-update >/dev/null 2>&1; then 
        rc-update del telescoped default
        if [ -d /etc/local.d ]; then
            rm -f /etc/local.d/telescoped.start
        fi
    elif command -v systemctl >/dev/null 2>&1; then 
        systemctl disable telescoped
        rm -f /etc/systemd/system/telescoped.service
    else
        echo "Unsupported unregister command"
        exit 0
    fi
    rm -f /etc/init.d/telescoped
fi

INSTALL_DIR=$CURRENT_DIR"/../telescope"
#delete install directory
rm -rf $INSTALL_DIR


echo "Success to uninstall telescope."
exit 0