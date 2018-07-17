# -*- coding: gbk -*-

########################################################
# Author:  m00316272
# Date:    2015-04-02
# Purpose: Create ReleaseNotes for Public Cloud Project.
########################################################

import os
import sys
import time
import fileinput
from ConfigParser import ConfigParser
import shutil


class MyConfigParser(ConfigParser):
    def optionxform(self, optionstr):
        return optionstr


def checkFile(fileName):
    if os.path.isfile(fileName):
        pass
    else:
        print "[ERROR]: No such file %s." % fileName
        sys.exit(1)
        
def getGitRevision():
    try:
        with open (".." + os.sep + ".git" + os.sep + "HEAD") as file:
            gitRevision = file.read().strip()
        return gitRevision    
    except:
        print "[ERROR]: No such file .git%sHEAD in current directory." % os.sep
        sys.exit(1)

    
def inputSectionInfo(section, value):
    try:
        for line in fileinput.input(tmpReleaseNoteFile, inplace = True):
            print line.replace(section, value)
    except Exception, e:
        print "[ERROR]: " % str(e)

        
def getRestrictions():
    content = ""
    for keyValue in cf.items("resctrictions"):
        scope = keyValue[0].strip()
        content = content + "<b>%s</b><br>\n" % (scope)
        info = keyValue[1].strip()
        if ";" in info:
            for subInfo in info.split(";"):
                content = content + "%s<br>\n" % (subInfo)
        else:
            content = content + "%s<br>\n" % (info)
    content = content.strip()        
    return content

            
def getSectionInfo(section, confFile, cf):
    content = ""
    for keyValue in cf.items(section):
        content = content + "<tr><td width=150px>%s</td><td width=400px>%s</td></tr>\n" % (keyValue[0], keyValue[1])
    content = content.strip()
    return content

    
def adjustFormat():
    with open(tmpReleaseNoteFile, "r") as infile:
        with open(releaseNoteFile, "w") as outfile:
            for line in infile.readlines():
                if line.split():
                    outfile.writelines(line)

                    
if __name__ == '__main__':
    if len(sys.argv) != 2:
        print "[ERROR]: No such parameter."
        sys.exit(1)
    
    packageLink = sys.argv[1].strip()
    
    tmpReleaseNoteFile = "ReleaseNotes_temp.html"
    releaseNoteFile = "ReleaseNotes.html"
    releaseNoteTemplate = "ReleaseNotes_template.html"

    checkFile(releaseNoteTemplate)
    
    releaseDate = time.strftime("%Y-%m-%d", time.localtime())
    releaseTime = time.strftime("%H:%M:%S", time.localtime())
    
    confFile = "config.txt"
    cf = MyConfigParser()
    cf.read(confFile)

    
    projectName = cf.get("basic_infomation", "projectName")
    releaseVersion = cf.get("basic_infomation", "releaseVersion")
    #releaseStatus = cf.get("basic_infomation", "releaseStatus")
    emailAddress = cf.get("basic_infomation", "emailAddress")
    gitRevision = getGitRevision()    
    shutil.copy(releaseNoteTemplate, tmpReleaseNoteFile)
    
    # Basic Information
    inputSectionInfo("{{projectName}}", projectName)
    inputSectionInfo("{{releaseVersion}}", releaseVersion)
    inputSectionInfo("{{gitRevision}}", gitRevision)
    #inputSectionInfo("{{releaseStatus}}", releaseStatus)
    inputSectionInfo("{{releaseDate}}", releaseDate)
    inputSectionInfo("{{releaeTime}}", releaseTime)
    inputSectionInfo("{{emailAddress}}", emailAddress)
    inputSectionInfo("{{packageLink}}", packageLink)

    # ImportantNotes
    importantNotes = getSectionInfo("importantnotes", confFile, cf)
    inputSectionInfo("{{importantNotes}}", importantNotes)
       
    # restrictions
    restrictions = getRestrictions()
    inputSectionInfo("{{restrictions}}", restrictions)
    
    # Components
    components = getSectionInfo("components", confFile, cf)
    inputSectionInfo("{{components}}", components)
    
    # New Features
    newFeatures = getSectionInfo("new_features", confFile, cf)
    inputSectionInfo("{{newFeatures}}", newFeatures)
    
    # Fixed Issues
    fixedIssues = getSectionInfo("fixed_issues", confFile, cf)
    inputSectionInfo("{{fixedIssues}}", fixedIssues)
    
    adjustFormat()
    os.remove(tmpReleaseNoteFile)

