# Copyright (c) 2022 Institute of Software, Chinese Academy of Sciences (ISCAS)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM debian:experimental as builder

WORKDIR /work/

RUN apt update; apt install golang-1.19 make curl wget -y

COPY . .

ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/lib/go-1.19/bin
RUN make -f Makefile

FROM debian:experimental

ENV LANG C.UTF-8
ENV TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive

RUN set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
		ca-certificates \
		netbase \
		tzdata \
		iputils-ping \
	; \
	apt remove docker.io -y ; \
	rm -rf /var/lib/apt/lists/*

COPY --from=builder /work/build/eulixspace-upgrade /usr/local/bin/eulixspace-upgrade
#COPY --from=builder /work/supervisord.conf /etc/supervisor/supervisord.conf

EXPOSE 5681

CMD ["/usr/local/bin/eulixspace-upgrade"]
