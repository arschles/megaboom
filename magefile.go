//go:build mage
// +build mage

package main

import (
	"context"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Hey mg.Namespace

func (Hey) DockerBuild(ctx context.Context) error {
	return sh.RunV(
		"docker",
		"build",
		"-t",
		"ghcr.io/arschles/hey:latest",
		"-f",
		"hey/Dockerfile",
		"./hey",
	)
}

func (Hey) DockerPush(ctx context.Context) error {
	return sh.RunV(
		"docker",
		"push",
		"ghcr.io/arschles/hey:latest",
	)
}

func Build(ctx context.Context) error {
	return sh.RunV(
		"go",
		"build",
		"-o",
		"./bin/megaboom",
		".",
	)
}

func DockerBuild(ctx context.Context) error {
	return sh.RunV(
		"docker",
		"build",
		"-t",
		"ghcr.io/arschles/megaboom:latest",
		".",
	)
}

func DockerBuildACR(ctx context.Context) error {
	return sh.RunV(
		"az",
		"acr",
		"build",
		"--image",
		"megaboom",
		"--registry",
		"testingkeda",
		"--file",
		"Dockerfile",
		".",
	)
}

func DockerPush(ctx context.Context) error {
	return sh.RunV(
		"docker",
		"push",
		"ghcr.io/arschles/megaboom:latest",
	)
}
