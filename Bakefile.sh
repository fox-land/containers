# shellcheck shell=bash

task.build() {
	cd './debian'

	local commit=
	commit=$(git rev-parse --short HEAD)

	local args=()
	for arg do case $arg in
	--bypass-cache)
		args+=('--pull' '--no-cache')
	esac done

	# ubuntu: bionic, focal, jammy (beta)
	# debian: stretch, buster, bullseye
	# centos: centos7 
	# opensuseleap: 15.3, 15.4 (beta)
	# centos-stream: stream8 stream9
	# fedora: 35
	
	local distro_variant=
	for distro_variant in bullseye buster stretch; do
		local date=
		date=$(date --rfc-3339=seconds)

		bake.info "Building $distro_variant"
		docker build \
			--build-arg ARG_DISTRO_VARIANT="$distro_variant" \
			--build-arg ARG_GIT_COMMIT="$commit" \
			--file "./Containerfile" \
			--label org.opencontainers.image.title="Fox build for $distro_variant" \
			--label org.opencontainers.image.description="Fox build for $distro_variant" \
			--label org.opencontainers.image.authors="Edwin Kofler <edwin@kofler.dev" \
			--label org.opencontainers.image.vendor="" \
			--label org.opencontainers.image.url="https://github.com/hyperupcall/containers" \
			--label org.opencontainers.image.documentation="https://github.com/hyperupcall/containers" \
			--label org.opencontainers.image.licenses="" \
			--label org.opencontainers.image.version="" \
			--label org.opencontainers.image.revision="$commit" \
			--label org.opencontainers.image.created="$date" \
			--tag "fox-debian-$distro_variant" \
			--tag "git-commit-$commit" \
			"${args[@]}" \
			.
	done
}