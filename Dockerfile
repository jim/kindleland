FROM ubuntu:latest
RUN apt-get update
RUN apt-get -y install gperf help2man bison texinfo flex gawk git build-essential autoconf libncurses5-dev curl wget file
WORKDIR /root/src
RUN git clone https://github.com/koreader/koxtoolchain.git
WORKDIR /root/src/koxtoolchain
ENV CT_EXPERIMENTAL=y
ENV CT_ALLOW_BUILD_AS_ROOT=y
ENV CT_ALLOW_BUILD_AS_ROOT_SURE=y
# RUN ./gen-tc.sh kindle

WORKDIR /root/src
RUN git clone https://github.com/NiLuJe/FBInk
WORKDIR /root/src/FBInk
RUN git submodule update --init

