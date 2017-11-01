FROM ubuntu:14.04
MAINTAINER MarkLux
# ADD sources.list /etc/apt/sources.list
RUN apt-get update&&apt-get -y install \
    gcc g++\
    libseccomp-dev\
    git\
    wget\
    rsync\
    software-properties-common\
    python
RUN add-apt-repository ppa:openjdk-r/ppa
RUN apt-get update&&apt-get -y install openjdk-7-jdk
RUN mkdir -p /home/judge
WORKDIR /home/judge
RUN wget https://studygolang.com/dl/golang/go1.9.2.linux-amd64.tar.gz&&\
    tar -xvzf go1.9.2.linux-amd64.tar.gz -C /usr/local&&\
    rm -rf go1.9.2.linux-amd64.tar.gz&&\
    mkdir -p /home/judge/go /home/judge/testcases /home/judge/submissions
    # echo "export GOROOT=/usr/local/go\nexport PATH=$PATH:/usr/local/go/bin\nexport GOPATH=/home/judge/go" >> ~/.bashrc
ENV GOROOT /usr/local/go
ENV GOPATH /home/judge/go
ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin
RUN mkdir -p /home/judge/go/src/golang.org/x /usr/lib/judger /home/judge/go/src/github.com/shirou
WORKDIR /home/judge/go/src/golang.org/x
RUN git clone https://github.com/golang/sys.git
WORKDIR /home/judge/go/src/github.com/shirou
RUN git clone https://github.com/shirou/gopsutil.git
RUN go get -v github.com/MarkLux/Judger_GO&&\
    go get -v github.com/MarkLux/JudgeServer&&\
    go get -u -v github.com/gin-gonic/gin&&\
    go get -u -v github.com/satori/go.uuid
WORKDIR /home/judge/go/src/github.com/MarkLux/Judger_GO
RUN mv libjudger.so /usr/lib/judger/libjudger.so
RUN apt-get clean \
    && apt-get autoclean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* 
WORKDIR /home/judge/go/src/github.com/MarkLux/JudgeServer
CMD go run main.go
EXPOSE 8090