# THIS FILE IS AUTOGENERATED! DO NOT EDIT!

FROM debian:12

RUN \
	useradd \
		--comment 'User' \
		--home-dir '/home/user' \
		--expiredate '' \
		--inactive '-1' \
		--create-home \
		--password 'password' \
		--shell '/bin/bash' \
		--user-group \
		'user'; \
	apt-get update -y ; \
	apt-get install -y build-essential git sudo autoconf
USER user
WORKDIR '/home/user'

RUN \
	git clone https://git.savannah.gnu.org/git/bash.git ; \
	cd bash ; \
	git switch --detach bash-4.4
WORKDIR '/home/user/bash'
RUN ./configure --prefix=/usr/local
RUN make
# RUN echo 'password' | sudo -S make install

	# ENTRYPOINT [ "/usr/bin/bash" ]
