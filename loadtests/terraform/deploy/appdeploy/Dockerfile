FROM public.ecr.aws/amazonlinux/amazonlinux:2.0.20230515.0-arm64v8

ENV OTP_VERSION="26.0" \
    REBAR3_VERSION="3.20.0"
ENV OTP_DOWNLOAD_URL="https://github.com/erlang/otp/archive/OTP-${OTP_VERSION}.tar.gz" \
	OTP_DOWNLOAD_SHA256="29a50b17cb4bc326f01e4cbb2badb6a31ec586525f7c2f29166a0c6229b53f24" 

LABEL org.opencontainers.image.version=$OTP_VERSION

# docker buildx build  --load --platform linux/arm64 -t posilva/amz2-2023.0.20230515.0-elixir-v1.14.5-otp-26.0:latest  --progress plain  -f Dockerfile .
# docker run -it posilva/amz2-2023.0.20230515.0-elixir-v1.14.5-otp-26.0:latest 

WORKDIR /opt/

RUN yum install openssl-devel ncurses-devel -y
RUN yum install unixODBC-devel lksctp-tools-devel gtk3-devel unixODBC lksctp-tools gtk3 tar -y
RUN yum groupinstall "Development Tools" -y 
RUN set -xe \
	&& curl -fSL -o otp-src.tar.gz "$OTP_DOWNLOAD_URL" \
	&& echo "$OTP_DOWNLOAD_SHA256  otp-src.tar.gz" | sha256sum -c - \
	&& export ERL_TOP="/usr/src/otp_src_${OTP_VERSION%%@*}" \
	&& mkdir -vp $ERL_TOP \
	&& tar -xzf otp-src.tar.gz -C $ERL_TOP --strip-components=1 \
	&& rm otp-src.tar.gz \
	&& ( cd $ERL_TOP \
	  && ./otp_build autoconf \
	  && ./configure --build="aarch64-linux-gnu" \
      &&  make -j16 \
	  && make install ) \
    && find /usr/local -name examples | xargs rm -rf \
	&& rm -rf $ERL_TOP \
	&& echo "$ERL_TOP" 



ENV ELIXIR_VERSION=v1.14.5 \
    LANG=C.UTF-8

ENV ELIXIR_DOWNLOAD_URL="https://github.com/elixir-lang/elixir/archive/${ELIXIR_VERSION}.tar.gz" \
    ELIXIR_DOWNLOAD_SHA256="2ea249566c67e57f8365ecdcd0efd9b6c375f57609b3ac2de326488ac37c8ebd"


RUN set -xe \
 pwd \
 && curl -fSL -o /opt/elixir-src.tar.gz $ELIXIR_DOWNLOAD_URL \
 && ls -la 

RUN echo "$ELIXIR_DOWNLOAD_SHA256  elixir-src.tar.gz" | sha256sum -c - \
    && mkdir -p /usr/local/src/elixir \
    && tar -xzC /usr/local/src/elixir --strip-components=1 -f elixir-src.tar.gz \
    && rm elixir-src.tar.gz \
    && cd /usr/local/src/elixir \
    && make install clean \
    && find /usr/local/src/elixir/ -type f -not -regex "/usr/local/src/elixir/lib/[^\/]*/lib.*" -exec rm -rf {} + \
    && find /usr/local/src/elixir/ -type d -depth -empty -delete

RUN yum install yum-utils -y
RUN yum autoremove unixODBC-devel lksctp-tools-devel gtk3-devel make -y
RUN package-cleanup --leaves --all
RUN yum autoremove yum-utils -y 
RUN yum autoremove -y
RUN yum clean packages
RUN yum clean metadata
RUN yum clean all
RUN rm -rf /var/cache/dnf
RUN yum groupremove "Development Tools" -y 

COPY ./scream /opt/scream

CMD [ "/bin/bash" ]




