package main

import (
	"context"
	"flag"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"kratos-k8s-job/api/scheduler/v1"
	"os"

	"kratos-k8s-job/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	logger log.Logger

	id, _ = os.Hostname()

	bc   conf.Bootstrap
	done = make(chan bool)
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	app := kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
		kratos.AfterStart(runCommand),
	)
	return app
}

func main() {
	flag.Parse()
	logger = log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	//when the channel gets the message then stop the app.
	go func() {
		<-done
		app.Stop()
	}()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}

}

func runCommand(context.Context) error {
	msg := "runCommand msg"

	l := log.NewHelper(logger)

	l.Info(msg)

	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9000"),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := v1.NewJobClient(conn)
	reply, err := client.ExecuteJob(context.Background(), &v1.ExecuteJobRequest{})
	if err != nil {
		l.Error(err)
	}
	l.Infof("[grpc] ExecuteJob %+v", reply)

	done <- true
	return nil
}
