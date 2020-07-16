FROM alpine:3.12

RUN echo $'\
@edge http://dl-cdn.alpinelinux.org/alpine/edge/main\n\
@edge http://dl-cdn.alpinelinux.org/alpine/edge/testing\n\
@edge http://dl-cdn.alpinelinux.org/alpine/edge/community' >> /etc/apk/repositories

RUN apk add --no-cache \
      python3 py3-pip py3-setuptools py3-wheel \
      py3-virtualenv \
      py3-pillow \
      py3-aiohttp \
      py3-magic \
      py3-ruamel.yaml \
      py3-commonmark@edge \
      # Other dependencies
      ca-certificates \
      su-exec \
      # encryption
      olm-dev \
      py3-cffi \
	  py3-pycryptodome \
      py3-unpaddedbase64 \
      py3-future

COPY requirements.txt /opt/mautrix-twitter/requirements.txt
COPY optional-requirements.txt /opt/mautrix-twitter/optional-requirements.txt
WORKDIR /opt/mautrix-twitter
RUN apk add --virtual .build-deps python3-dev libffi-dev build-base \
 && pip3 install -r requirements.txt -r optional-requirements.txt \
 && apk del .build-deps

COPY . /opt/mautrix-twitter
RUN apk add git && pip3 install .[e2be] && apk del git \
  # This doesn't make the image smaller, but it's needed so that the `version` command works properly
  && cp mautrix_twitter/example-config.yaml . && rm -rf mautrix_twitter

VOLUME /data
ENV UID=1337 GID=1337

CMD ["/opt/mautrix-twitter/docker-run.sh"]
