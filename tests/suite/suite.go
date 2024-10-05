package suite

import (
	"context"
	ssov1 "github.com/GolangLessons/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"sso/internal/config"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Config     *config.Config
	AuthClient ssov1.AuthClient
}

const (
	grpcHost = "localhost"
)

// New creates new test suite.
func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath(configPath())

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials())) // only for tests
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Config:     cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func configPath() string {
	const key = "CONFIG_PATH"

	//if v := os.Getenv(key); v != "" {
	//	return v
	//}

	return "../config/local_tests.yaml"
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
