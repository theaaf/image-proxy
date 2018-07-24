FROM golang:1.10

RUN apt-get update

RUN apt-get install -y autoconf librsvg2-dev

RUN curl -L https://github.com/NixOS/patchelf/archive/27ffe8ae871e7a186018d66020ef3f6162c12c69.tar.gz -o /tmp/patchelf.tar.gz
RUN cd /tmp && tar -xvf ./patchelf.tar.gz && cd patchelf-27ffe8ae871e7a186018d66020ef3f6162c12c69 && ./bootstrap.sh && ./configure && make install

