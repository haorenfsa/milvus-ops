package ctrl

import (
	"context"

	"github.com/pkg/errors"
	"github.com/tevino/log"
)

func wrapAndLog(ctx context.Context, err error, msg string) error {
	traceID, ok := ctx.Value("traceID").(string)
	if !ok {
		log.Error("get trace id failed")
	}
	err = errors.Wrap(err, "list namespaces error")
	if err != nil {
		log.Errorf("[%s] %s", traceID, err)
	}
	return err
}
