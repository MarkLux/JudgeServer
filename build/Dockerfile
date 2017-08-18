FROM ubuntu:14.04
MAINTAINER MarkLux
# ADD sources.list /etc/apt/sources.list
RUN apt-get update&&apt-get -y install \
    gcc g++\
    libseccomp-dev\
    git\
    wget
RUN mkdir -p /home/judge
WORKDIR /home/judge
RUN wget https://www.golangtc.com/static/go/1.8.3/go1.8.3.linux-amd64.tar.gz&&\
    tar -xvzf go1.8.3.linux-amd64.tar.gz -C /usr/local&&\
    rm -rf go1.8.3.linux-amd64.tar.gz&&\
    mkdir -p /home/judge/go /home/judge/testcases /home/judge/submissions
    # echo "export GOROOT=/usr/local/go\nexport PATH=$PATH:/usr/local/go/bin\nexport GOPATH=/home/judge/go" >> ~/.bashrc
ENV GOROOT /usr/local/go
ENV GOPATH /home/judge/go
ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin
RUN go get -u -v github.com/MarkLux/Judger_GO&&\
    go get -u -v github.com/MarkLux/JudgeServer&&\
    go get -u -v github.com/gin-gonic/gin&&\
    go get -u -v github.com/satori/go.uuid&&\
    # go get -u -v github.com/shirou/gopsutil&&\
    mkdir -p /home/judge/go/src/golang.org/x /usr/lib/judger /home/judge/go/src/github.com/shirou
WORKDIR /home/judge/go/src/golang.org/x
RUN git clone https://github.com/golang/sys.git
WORKDIR /home/judge/go/src/github.com/shirou
RUN git clone https://github.com/shirou/gopsutil.git
WORKDIR /home/judge/go/src/github.com/MarkLux/Judger_GO
RUN mv libjudger.so /usr/lib/judger/libjudger.so
WORKDIR /home/judge/go/src/github.com/MarkLux/JudgeServer
CMD go run main.go
EXPOSE 8090