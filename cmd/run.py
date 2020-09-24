import json
import os
import sys
import time
import urllib2


def observer_status():
    api_host = os.getenv('WEB_HOST', "http://localhost:4399")
    api = api_host + "/api/system/status"
    try:
        return json.load(urllib2.urlopen(api))
    except:
        print("can't connect api server")
        return {"data": {"substrate": True}}


def system_do(op):
    os.system('./subscan -conf ../configs start ' + op)


def main():
    op = []
    if len(sys.argv) == 1:
        os.system('./subscan -conf ../configs')
    elif sys.argv[1] == "substrate":
        op = ["substrate"]
    map(system_do, op)  # run,run,run
    print("start to listen observer status :",
          time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()))
    while len(op) > 0:
        j = observer_status()
        for i in range(len(op)):
            try:
                if not j["data"][op[i]]:
                    s = './subscan stop {observer} && ./subscan start {observer}'
                    os.system(s.format(observer=op[i]))
            except KeyError:
                pass
        time.sleep(60)


if __name__ == "__main__":
    main()
