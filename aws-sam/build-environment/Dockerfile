FROM amazonlinux:1

RUN yum-config-manager --enable epel
RUN yum install -y \
    gcc \
    gcc-c++ \
    gettext-devel \
    libxml2-devel \
    libpng-devel \
    libjpeg-turbo-devel \
    pixman-devel \
    gobject-introspection \
    gobject-introspection-devel \
    unzip \
    wget

ENV PKG_CONFIG_PATH=/usr/lib/pkgconfig
ENV LD_LIBRARY_PATH=/lib:/usr/lib:/usr/local/lib
ENV LDFLAGS=-L/usr/lib

WORKDIR /tmp
RUN wget https://ftp.pcre.org/pub/pcre/pcre-8.42.zip
RUN unzip pcre-8.42.zip && cd pcre-8.42 && ./configure --prefix=/usr --enable-unicode-properties && make && make install

WORKDIR /tmp
RUN wget https://download.savannah.gnu.org/releases/freetype/freetype-2.8.1.tar.gz
RUN tar xvf freetype-2.8.1.tar.gz && cd freetype-2.8.1 && ./configure --prefix=/usr && make && make install

WORKDIR /tmp
RUN wget https://download.gnome.org/sources/glib/2.48/glib-2.48.0.tar.xz
RUN tar xvfJ glib-2.48.0.tar.xz && cd glib-2.48.0 && ./configure --prefix=/usr && make && make install

WORKDIR /tmp
RUN wget https://ftp.gnome.org/pub/GNOME/sources/libcroco/0.6/libcroco-0.6.8.tar.xz
RUN tar xvfJ libcroco-0.6.8.tar.xz && cd libcroco-0.6.8 && ./configure --prefix=/usr && make && make install

WORKDIR /tmp
RUN wget https://ftp.gnome.org/pub/GNOME/sources/gdk-pixbuf/2.28/gdk-pixbuf-2.28.2.tar.xz
RUN tar xvfJ gdk-pixbuf-2.28.2.tar.xz && cd gdk-pixbuf-2.28.2 && ./configure --prefix=/usr --without-libtiff && make && make install

WORKDIR /tmp
RUN wget https://www.freedesktop.org/software/fontconfig/release/fontconfig-2.10.91.tar.gz
RUN tar xvf fontconfig-2.10.91.tar.gz && cd fontconfig-2.10.91 && ./configure --prefix=/usr --enable-libxml2 && make && make install

RUN curl https://sh.rustup.rs -sSf | bash -s -- -y
ENV PATH=$PATH:/root/.cargo/bin

WORKDIR /tmp
RUN wget https://ftp.gnu.org/gnu/automake/automake-1.15.1.tar.xz
RUN tar xvfJ automake-1.15.1.tar.xz && cd automake-1.15.1 && ./configure --prefix=/usr && make && make install

WORKDIR /tmp
RUN wget https://cairographics.org/releases/cairo-1.16.0.tar.xz
RUN tar xvfJ cairo-1.16.0.tar.xz && cd cairo-1.16.0 && ./configure --prefix=/usr && make && make install

WORKDIR /tmp
RUN yum install -y bzip2
RUN wget https://www.freedesktop.org/software/harfbuzz/release/harfbuzz-1.4.2.tar.bz2
RUN tar xvf harfbuzz-1.4.2.tar.bz2 && cd harfbuzz-1.4.2 && ./configure --prefix=/usr && make && make install

WORKDIR /tmp
RUN wget https://ftp.gnome.org/pub/GNOME/sources/pango/1.38/pango-1.38.1.tar.xz
RUN tar xvfJ pango-1.38.1.tar.xz && cd pango-1.38.1 && ./configure --prefix=/usr && make && make install

WORKDIR /tmp
RUN wget https://ftp.gnome.org/pub/GNOME/sources/librsvg/2.42/librsvg-2.42.7.tar.xz
RUN tar xvfJ librsvg-2.42.7.tar.xz && cd librsvg-2.42.7 && ./configure --prefix=/usr && make && make install

RUN yum install -y golang
ENV GOPATH=/go
RUN mkdir -p /go/bin
RUN wget -O - https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
ENV PATH=$PATH:/go/bin

RUN wget -O /tmp/patchelf.tar.gz https://github.com/NixOS/patchelf/archive/27ffe8ae871e7a186018d66020ef3f6162c12c69.tar.gz
RUN cd /tmp && tar -xvf ./patchelf.tar.gz && cd patchelf-27ffe8ae871e7a186018d66020ef3f6162c12c69 && ./bootstrap.sh && ./configure && make install
