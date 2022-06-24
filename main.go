package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	cli "github.com/urfave/cli/v2"
)

// # 				--label org.opencontainers.image.vendor="" \
// # --label org.opencontainers.image.licenses="" \
// # 				--label org.opencontainers.image.version="" \

// # ubuntu: bionic, focal, jammy (beta)
// # debian: stretch, buster, bullseye
// # centos: centos7
// # opensuseleap: 15.3, 15.4 (beta)
// # centos-stream: stream8 stream9
// # fedora: 35

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func buildBash() error {
	os.Chdir("./containers/bash")
	for _, version := range []string{"4.3", "4.4", "5.0", "5.1"} {
		dockerBuild := exec.Command("docker", []string{
			"build",
			"--file", "./" + version + ".Containerfile",
			"--tag", fmt.Sprintf("fox.%s:%s", "bash", version),
			"--tag", fmt.Sprintf("ghcr.io/hyperupcall/fox.%s:%s", "bash", version),
			"--label", "org.opencontainers.image.source=" + fmt.Sprintf("https://github.com/%s", "hyperupcall/containers"),
			".",
		}...)
		dockerBuild.Stdin = os.Stdin
		dockerBuild.Stdout = os.Stdout
		dockerBuild.Stderr = os.Stderr
		err := dockerBuild.Run()
		if err != nil {
			return fmt.Errorf("%s", err)
		}

		dockerPush := exec.Command("docker", "push", fmt.Sprintf("ghcr.io/hyperupcall/fox.%s:%s", "bash", version))
		dockerPush.Stdin = os.Stdin
		dockerPush.Stdout = os.Stdout
		dockerBuild.Stderr = os.Stderr
		err = dockerPush.Run()
		if err != nil {
			return fmt.Errorf("%s", err)
		}

	}

	return nil
}

func build(bypassCache bool, noPush bool, container string) error {
	rawCommit, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	commit := strings.TrimSpace(string(rawCommit))
	handle(err)

	fmt.Printf("Current Git Commit: %s\n", commit)

	hasRanAtLeastOnce := false
	for _, pair := range []struct {
		distro   string
		variants []string
	}{
		{
			distro:   "debian",
			variants: []string{"bullseye", "buster", "stretch"},
		},
		{
			distro:   "ubuntu",
			variants: []string{"jammy", "focal", "bionic"},
		},
	} {
		for _, variant := range pair.variants {
			if len(container) > 0 {
				arr := strings.Split(container, ":")
				specifiedDistro := arr[0]
				specifiedVariant := arr[1]

				if specifiedDistro != pair.distro || specifiedVariant != variant {
					continue
				}
			}
			hasRanAtLeastOnce = true

			rawDate, err := exec.Command("date", "--rfc-3339=seconds").Output()
			handle(err)
			date := strings.TrimSpace(string(rawDate))

			fmt.Println(pair.distro, variant)
			fmt.Printf("Current Date: %s\n", date)

			author := struct {
				name  string
				repo  string
				email string
			}{
				"Edwin Kofler",
				"hyperupcall/containers",
				"edwin@kofler.dev",
			}

			dockerArgs := []string{
				"build",
				"--build-arg", "ARG_DISTRO_VARIANT=" + variant,
				"--build-arg", "ARG_GIT_COMMIT=" + commit,
				"--file", "./" + pair.distro + "/Containerfile",
				"--label", "org.opencontainers.image.title=" + fmt.Sprintf("Fox build for %s", variant),
				"--label", "org.opencontainers.image.description=" + fmt.Sprintf("Fox build for %s", variant),
				"--label", "org.opencontainers.image.authors=" + fmt.Sprintf("%s <%s>", author.name, author.email),
				"--label", "org.opencontainers.image.url=" + fmt.Sprintf("https://github.com/%s", author.repo),
				"--label", "org.opencontainers.image.documentation=" + fmt.Sprintf("https://github.com/%s", author.repo),
				"--label", "org.opencontainers.image.source=" + fmt.Sprintf("https://github.com/%s", author.repo),
				"--label", "org.opencontainers.image.revision=" + commit,
				"--label", "org.opencontainers.image.created=" + date,
				"--tag", fmt.Sprintf("fox.%s", pair.distro),
				"--tag", fmt.Sprintf("ghcr.io/hyperupcall/fox.%s", pair.distro),
			}
			if bypassCache {
				dockerArgs = append(dockerArgs, []string{"--pull", "--no-cache"}...)
			}
			dockerArgs = append(dockerArgs, "./"+pair.distro)

			dockerBuild := exec.Command("docker", dockerArgs...)
			dockerBuild.Stdin = os.Stdin
			dockerBuild.Stdout = os.Stdout
			dockerBuild.Stderr = os.Stderr
			err = dockerBuild.Run()
			if err != nil {
				log.Fatalln(err)
			}

			if noPush {
				fmt.Println("Skipping push...")
			} else {
				dockerPush := exec.Command("docker", "push", fmt.Sprintf("ghcr.io/hyperupcall/fox.%s:latest", pair.distro))
				dockerPush.Stdin = os.Stdin
				dockerPush.Stdout = os.Stdout
				dockerPush.Stderr = os.Stderr
				err = dockerPush.Run()
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}

	if !hasRanAtLeastOnce {
		return fmt.Errorf("No containers were either build or pushed")
	}

	return nil
}

func main() {
	app := &cli.App{
		Name:    "Container Builder",
		Usage:   "Builder",
		Version: "0.1.0",
		Authors: []*cli.Author{
			{
				Name:  "Edwin Kofler",
				Email: "edwin@kofler.dev",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "build",
				Usage: "Build containers",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "bypass-cache",
						Usage: "Bypass the cache",
					},
					&cli.BoolFlag{
						Name:  "no-push",
						Usage: "Do not push to the OCI Registry",
					},
					&cli.StringFlag{
						Name:  "container",
						Usage: "Build a specific container",
					},
				},
				Action: func(ctx *cli.Context) error {
					return build(ctx.Bool("bypass-cache"), ctx.Bool("no-push"), ctx.String("container"))
				},
			},
			{
				Name:  "bash",
				Usage: "Do Bash things",
				Action: func(ctx *cli.Context) error {
					return buildBash()
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
