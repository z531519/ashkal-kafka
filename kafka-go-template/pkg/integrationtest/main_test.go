//go:build integration
// +build integration

package integrationtest

import (
	"context"
	"os"
	"testing"

	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

////
// Integration Test for the CustomerHandler
//
// The integration test can be run with:
// go test ./... -tags integration
//
// This test use testcontainers

func TestMain(m *testing.M) {
	identifier := tc.StackIdentifier("integration-test")
	compose, err := tc.NewDockerComposeWith(tc.WithStackFiles("../../docker-compose.yml"), identifier)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	compose.Up(ctx, tc.Wait(true))

	code := m.Run()

	err = compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal)
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}
