# Judge Server

> a judge server implement in golang referred to  [QDU JudgeServer](https://github.com/QingdaoU/JudgeServer) 

## install & depencies

packages needed:

* Judger_GO (https://github.com/MarkLux/Judger_GO)
* Gin (https://github.com/gin-gonic/gin)
* Go.uuid (https://github.com/satori/go.uuid)
* gopsutil (https://github.com/shirou/gopsutil)

## prepare

Edit the `config/consts.go` to set the directories needed rightly (REMEMBER TO CHECK THE DIR PERMISSIONS !)

Install gcc,g++,python,javac language support before starting judge.

## test cases

only `.in`&`.out` files in pair with same file name can be recognized.

## usage

The package run in the way of gin webserver (defautly listen port 8090,you can change it in `main.go`)

Just run `go run main.go` to start the server!

## API

*this judge server currently doesn't support special judge,please waiting for further development.*

### ping ([GET] /ping)

get the server info.

response example:

```
{
    "action": "pong",
    "cpu": [
        0.13249921668799483
    ],
    "cpu_core": 2,
    "hostname": "0b35f2d21276",
    "judger_version": "0.1.0",
    "memory": 22.187548851023916
}
```

### judge(not for specail judge) ([POST] /judge)

request example:

```
{
	"src":"# include<stdio.h> \n int main()\n  {\n printf(\"1\");\n return 0;\n}",
	"language":"c",
	"max_cpu_time":1000,
	"max_memory":395671011,
	"test_case_id":1001
}
```

response example:

```
{
    "code": 0,
    "data": {
        "Passed": [
            {
                "CpuTime": 0,
                "Result": 0,
                "RealTime": 9,
                "Memory": 4182016,
                "Signal": 0,
                "Error": 0,
                "OutputMD5": "c4ca4238a0b923820dcc509a6f75849b"
            }
        ],
        "UnPassed": [
            {
                "CpuTime": 0,
                "Result": -1,
                "RealTime": 8,
                "Memory": 4177920,
                "Signal": 0,
                "Error": 0,
                "OutputMD5": "c4ca4238a0b923820dcc509a6f75849b"
            }
        ]
    }
}
```

data contains results for each test case (`with a .in and a .out file`),with Accpted ones in `Passed` array and others in `UnPassed` array.

## build & images

you can use `build/Dockerfile` to build the docker image to use

or use this repsository that I built on docker hub: https://hub.docker.com/r/marklux/judge_server/
