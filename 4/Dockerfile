# 我们需要引入到基础容器
FROM alpine:3.18

# 注意看这里我们的写法, 在sh中 && 可以表示下一条命令连续执行，而 \ 则是命令的分隔符号
# 思考：为什么这么写？
RUN mkdir /go && cd /go \
    && wget --no-check-certificate https://golang.google.cn/dl/go1.21.0.linux-amd64.tar.gz \
    && tar -C /usr/local -zxf go1.21.0.linux-amd64.tar.gz \
    && rm -rf /go/go1.21.0.linux-amd64.tar.gz \
    && mkdir /lib64 \
    && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 
   
# 配置系统环境变量
# 在dockerfile中对容器的系统环境变量配置统一采用ENV这个关键词定义
ENV GOPATH /go
ENV PATH /usr/local/go/bin:$GOPATH/bin:$PATH

# 这个命令可以让我们的docker容器在启动的时候就执行下面的命令
# 把原本在docker run中的命令放到dockerfile中，并示意启动容器的时候执行
# 但是如果在docker run后跟新的命令会代替CMD中的命令
CMD ["ping", "www.baodu.com"]