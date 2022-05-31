# shellcheck shell=bash

# 				--label org.opencontainers.image.vendor="" \
# --label org.opencontainers.image.licenses="" \
# 				--label org.opencontainers.image.version="" \

# ubuntu: bionic, focal, jammy (beta)
# debian: stretch, buster, bullseye
# centos: centos7
# opensuseleap: 15.3, 15.4 (beta)
# centos-stream: stream8 stream9
# fedora: 35

task.build() {
	local -a args=()
	local -a docker_args=()

	for arg do case $arg in
	--bypass-cache)
		docker_args+=('--pull' '--no-cache') ;;
	*)
		args+=("$arg") ;;
	esac done

	local commit=
	commit=$(git rev-parse --short HEAD)

	# Build one or more containers
	local -A distro_pairs=(
		[debian]='bullseye:buster:stretch'
		[ubuntu]='jammy:focal:bionic'
	)
	if [ -z "${args[0]}" ]; then
		local key= value=
		for key in "${!distro_pairs[@]}"; do
			local value="${distro_pairs[$key]}"
			util.build_container "$key|$value"
		done
	else
		local key='debian'
		local value="${distro_pairs[$key]}"
		util.build_container "$key|$value"
	fi
}


util.build_container() {
	local pair="$1"

	distro=${pair%%|*}
	IFS=':' read -ra distro_variants <<< "${pair#*|}"

	local distro_variant=
	for distro_variant in "${distro_variants[@]}"; do
		local date=
		date=$(date --rfc-3339=seconds)

		bake.info "Building $distro $distro_variant"
		docker build \
			--build-arg ARG_DISTRO_VARIANT="$distro_variant" \
			--build-arg ARG_GIT_COMMIT="$commit" \
			--file "./$distro/Containerfile" \
			--label org.opencontainers.image.title="Fox build for $distro_variant" \
			--label org.opencontainers.image.description="Fox build for $distro_variant" \
			--label org.opencontainers.image.authors="Edwin Kofler <edwin@kofler.dev" \
			--label org.opencontainers.image.url="https://github.com/hyperupcall/containers" \
			--label org.opencontainers.image.documentation="https://github.com/hyperupcall/containers" \
			--label org.opencontainers.image.source="https://github.com/hyperupcall/containers" \
			--label org.opencontainers.image.revision="$commit" \
			--label org.opencontainers.image.created="$date" \
			--tag "fox.$distro-$distro_variant" \
			--tag "ghcr.io/hyperupcall/fox.$distro" \
			"${docker_args[@]}" \
			"./$distro"

			docker push "ghcr.io/hyperupcall/fox.$distro:latest"
	done
}
