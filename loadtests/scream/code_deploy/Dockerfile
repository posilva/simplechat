# FROM arm64v8/elixir:1.14.5-otp-25
FROM posilva/amz2-2023.0.20230503.0-elixir-v1.14.5-otp-26.0:latest


ARG MIX_ENV=prod
ARG HEX_MIRROR_URL=https://repo.hex.pm

RUN /usr/local/bin/mix local.hex --force && \
    /usr/local/bin/mix local.rebar --force

RUN uname -a 

RUN yum update && yum install -y git 
RUN yum  install -y unzip 

WORKDIR /opt/app

COPY config/ ./config
COPY mix.exs mix.lock /opt/app/
COPY . .

RUN mix deps.get --only $MIX_ENV

RUN mix release  --path  /opt/scream  --force --overwrite --quiet

ENTRYPOINT [ "/opt/scream/bin/scream", "start_iex" ]