import os
import sys
import time
import json
import ConfigParser
import codecs

# sys.argv[0]:team_name  
# sys.argv[1]:package_name 
# sys.argv[3]:package_type 


class CaseSensConfigParser(ConfigParser.ConfigParser):
    # to be case sensitive
    def optionxform(self, optionstr):
        return optionstr

if __name__ == '__main__':
    if len(sys.argv) < 4 :
        print "Three arguments needed:"
        print "  1) team name, this value can be: ConsoleFramework|ECS_UI|APICOM|VPC|AutoScalling|CloudEye|MeterTicket|ApiGateway "
        print "  2) package name, this value can be: ConsoleFramework|ECS_UI|APICOM|VPC|AutoScalling|CloudEye|MeterTicket|ApiGateway "
        print "  3) package_type, this value can be: service|3rd|security"
        sys.exit()
    cfgfile = "config.txt"
    cf = CaseSensConfigParser()
    cf.read(cfgfile)
    sections = cf.sections()
    result = {}
   
    team_name = sys.argv[1] 
    package_name = sys.argv[2]
    package_type = sys.argv[3]
    release_date = time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(time.time())) 

    release_version = cf.get("basic_infomation","releaseVersion")

    release_note = {}
    playbook = {"playbook_package_name": cf.get("components","playbook_package_name")}
    components = {}
    coms = cf.items("components")
    for com in coms:
        components[com[0]] = com[1]
    # return result
    release_note["components"] = components
    release_note["playbook"] = playbook
    result = {"team_name": team_name, "release_date": release_date, "release_version": release_version}
    result = json.dumps(result, encoding="utf-8", ensure_ascii=False)

    jsonfile = team_name + "-" + release_version + "-release.json"
    
    jsonfile = codecs.open(jsonfile, "w", "utf-8")
    jsonfile.writelines(result)
    jsonfile.close()