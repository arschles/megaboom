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

type CLI mg.Namespace

func (CLI) Build(ctx context.Context) error {
	return sh.RunV(
		"go",
		"build",
		"-o",
		"./bin/cli",
		"./cmd/cli",
	)
}

type Worker mg.Namespace

func (Worker) Build(ctx context.Context) error {
	return sh.RunV(
		"go",
		"build",
		"-o",
		"./bin/megaboom-worker",
		"./cmd/worker",
	)
}

func (Worker) DockerBuildACR(ctx context.Context) error {
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

func (Worker) DockerPush(ctx context.Context) error {
	return sh.RunV(
		"docker",
		"push",
		"ghcr.io/arschles/megaboom:latest",
	)
}
