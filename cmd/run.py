import json
import os
import random
import sys
import time
import urllib2


def daemon_status():
    api_host = os.getenv('WEB_HOST', "http://localhost:4399")
    evo_api = api_host + "/api/system/status"
    try:
        req = urllib2.urlopen(evo_api)
        return json.load(req)
    except:
        print("can't connect api server")
        return {"data": {
            "substrate": True}
        }


def main():
    argv = sys.argv
    op, log_file = [], ""
    endpoint_pool = ["wss://crayfish.subscan.network/"]
    if len(argv) == 1:
        os.system('./subscan -conf ../configs')
    elif argv[1] == "substrate":
        op = ["substrate"]
        log_file = "../log/substrate_log"
    while log_file:
        map(system_do, op)  # run,run,run
        print("start to listen daemon status :", time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()))
        if len(op) > 0:
            j = daemon_status()
            for i in range(len(op)):
                try:
                    if not j["data"][op[i]]:
                        s = './subscan -conf ../configs stop {daemon} && ./subscan -conf ../configs start {daemon}'
                        if op[i] == "substrate":
                            s = "CHAIN_WS_ENDPOINT=" + random.sample(endpoint_pool, 1)[0] + s
                        os.system(s.format(daemon=op[i]))
                except KeyError:
                    pass
        time.sleep(30)


def system_do(op):
    os.system('./subscan -conf ../configs start ' + op)


if __name__ == "__main__":
    main()
