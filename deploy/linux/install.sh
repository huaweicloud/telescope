#!/bin/bash
OS_DESC=''
OTHERS_OS='OTHERS'
CURRENT_OS=''

support_os_list=("CENTOS" "SUSE" "DEBIAN" "EULER" "UBUNTU" "FEDORA" "ORACLE LINUX")
chkconfig_os_list=("CENTOS" "SUSE" "EULER" "FEDORA" "ORACLE LINUX")
update_rc_os_list=("DEBIAN" "UBUNTU")

telecoped_service="[Unit]
Description=telescoped service
After=network.target

[Service]
Type=simple
ExecStart=/etc/init.d/telescoped start
RemainAfterExit=yes
ExecStop=/etc/init.d/telescoped stop
KillMode=none

[Install]
WantedBy=multi-user.target
"
CES_AGENT_FLAG="ces"

# array_contains ${item} "${arr[@]}"
function array_contains () {
    local seeking=$1; shift
    local in=1
    for element; do
        if [[ ${element} == ${seeking} ]]; then
            in=0
            break
        fi
    done
    return ${in}
}

# query metadata to check if need to install ces agent
# 0:yes; 1:no
function isCesSelected()
{
    #get metadata
    local meta_json

    which curl
    if [[ $? -eq 0 ]]; then
        meta_json=$(curl -s -X GET --connect-timeout 7 http://169.254.169.254/openstack/latest/meta_data.json)
    else
        meta_json=$(wget http://169.254.169.254/openstack/latest/meta_data.json && cat meta_data.json)
        rm -rf meta_data.json
    fi

    AVAILABILITY_ZONE=$(echo ${meta_json} | grep -P '"availability_zone":\s*"(.*?)"' -o 2>/dev/null | awk -F ":" '{print $2}' | sed 's/\"//g' | tr -d ' \t\n\r\f')
    # Debian 7.5.0 64bit does not have libpcre.so.3
    if [[ -z ${AVAILABILITY_ZONE} ]]; then
        AVAILABILITY_ZONE=$(echo ${meta_json} | python -mjson.tool | grep availability_zone | awk -F '"' {'print $4'})
    fi
    agent_list_string=$(echo ${meta_json} | grep -P '"__support_agent_list":\s*"(.*?)"' -o 2>/dev/null | awk -F ":" '{print $2}' | sed 's/\"//g')
    IFS=', ' read -r -a agent_array <<< "${agent_list_string}"

    array_contains ${CES_AGENT_FLAG} "${agent_array[@]}"
    is_ces_in_list=$?

    if [[ ${is_ces_in_list} -eq 0 ]]; then
        echo "ces flag FOUND in __support_agent_list"
        return 0
    else
        echo "ces flag NOT FOUND in __support_agent_list"
        return 1
    fi
}

# 0:yes; 1:no
function isStartedByTelescopeInstall()
{
    PPID_CMDLINE=$(cat /proc/${PPID}/cmdline)
    echo ${PPID_CMDLINE} | grep "TelescopeInstall"
    if [[ $? -eq 0 ]]; then
        return 0
    else
        return 1
    fi
}

isCesSelected
is_ces_selected=$?
isStartedByTelescopeInstall
is_started_by_ti=$?

# 是被TelescopeInstall启动 && ces没有被打标签
if [[ ${is_started_by_ti} -eq 0 ]] && [[ ${is_ces_selected} -ne 0 ]]; then
    rm -rf /usr/local/telescope_linux_amd64*
    exit 0
fi


if [ "`id -u`" = "0" ] || [ "`id -g`" = "0" ] ; then
    echo "Current user is root."
else
    echo "Current user is not root, please use root user install or command [sudo sh install.sh]."
    exit 0
fi

getStatus()
{
    PARENT_PIDS=($(pgrep -l -f -P 1,0 "telescope$" | awk '{print $1}'))
    CHILD_PIDS=()


    for ppid in ${PARENT_PIDS[*]}
    do
       CHILD_PIDS=(${CHILD_PIDS[*]} $(pgrep -l -f -P "${ppid}" "telescope$" | awk '{print $1}'))
    done
    if [ ${#PARENT_PIDS[*]} == 1 -a ${#CHILD_PIDS[*]} == 1 ]; then
       return 0
    fi
    if [ ${#PARENT_PIDS[*]} == 0 -a ${#CHILD_PIDS[*]} == 0 ]; then
      # all telescope process is stopped
       return 1
    elif [ ${#PARENT_PIDS[*]} == 0 -o ${#CHILD_PIDS[*]} == 0 ]; then 
      # "Daemon parent process or Business child process telescope is not running"
       return 2
    else
       echo "The running parent process: ${PARENT_PIDS[*]}, the running child process: ${CHILD_PIDS[*]}"
       return 3
    fi

}

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

getStatus
status=$?
if [ ${status} == 0 ]; then
    echo "telescope is running, so can't be installed"
    exit 1
fi

getCurrentPath
getOS


chmod 755 ${CURRENT_DIR} -R
chown root ${CURRENT_DIR} -R
chgrp root ${CURRENT_DIR} -R

chmod 755 telescoped -R

INSTALL_DIR=/usr/local/telescope
if [[ "$1" && ! -d "$1" ]]; then
    echo "$1" is not a directory! Install telescope failed.
    exit -1
fi
if [[ "$1" && -d "$1" ]]; then
    INSTALL_DIR="$1"/telescope
fi

old=$(grep '^BIN_DIR=' ${CURRENT_DIR}"/telescoped")
sed -i 's#^'"$old"'#BIN_DIR='''"$INSTALL_DIR"'''#g' ${CURRENT_DIR}"/telescoped"

# get current linux os version
CURRENT_OS=${OTHERS_OS}
for support_os in "${support_os_list[@]}"
do 
    if [ `echo ${OS_DESC} | tr [a-z] [A-Z] | grep "${support_os}" | wc -l` -ge 1 ] ; then
        CURRENT_OS=${support_os}
    fi
done
   
echo "Current linux release version : ${CURRENT_OS}"
echo "Start to install telescope..."

mkdir -p ${INSTALL_DIR}
mkdir -p ${INSTALL_DIR}/log
cp -R ${CURRENT_DIR}/bin ${INSTALL_DIR}
cp ${CURRENT_DIR}/telescoped ${INSTALL_DIR}
cp ${CURRENT_DIR}/uninstall.sh ${INSTALL_DIR}

# add telescoped service and set up autostart
if [[ "${chkconfig_os_list[@]}" =~ $CURRENT_OS ]]; then
    echo "In chkconfig "
    cp ${CURRENT_DIR}"/telescoped" /etc/init.d
    chkconfig --add telescoped
    chkconfig telescoped on
elif [[ "${update_rc_os_list[@]}" =~ $CURRENT_OS ]]; then
    echo "In update-rc.d "
    cp ${CURRENT_DIR}"/telescoped" /etc/init.d
    update-rc.d telescoped defaults
else
    if command -v chkconfig >/dev/null 2>&1; then 
        cp ${CURRENT_DIR}"/telescoped" /etc/init.d
        chkconfig --add telescoped
        chkconfig telescoped on
    elif command -v update-rc.d >/dev/null 2>&1; then 
        cp ${CURRENT_DIR}"/telescoped" /etc/init.d
        update-rc.d telescoped defaults
    elif command -v rc-update >/dev/null 2>&1; then 
        cp ${CURRENT_DIR}"/telescoped" /etc/init.d
        rc-update add telescoped default
        if [ -d /etc/local.d ]; then
            touch /etc/local.d/telescoped.start
            chmod 755 /etc/local.d/telescoped.start
            echo "/etc/init.d/telescoped start" > /etc/local.d/telescoped.start
        fi
    elif command -v systemctl >/dev/null 2>&1; then 
        cp ${CURRENT_DIR}"/telescoped" /etc/init.d
        touch /etc/systemd/system/telescoped.service
        chmod 644 /etc/systemd/system/telescoped.service
        echo "$telecoped_service" > /etc/systemd/system/telescoped.service
        systemctl enable telescoped
    else
        echo "Unsupported register command, autostarts unsupported linux"
        sh telescoped start
        exit 0
    fi
fi

echo "Success to install telescope to dir: ${INSTALL_DIR}."

# start telescope
if command -v systemctl >/dev/null 2>&1; then
    systemctl unmask telescoped.service
fi

if command -v service >/dev/null 2>&1; then
    service telescoped start
elif command -v rc-service >/dev/null 2>&1; then
    rc-service telescoped start
elif command -v systemctl >/dev/null 2>&1; then
    systemctl start telescoped
fi

exit 0
