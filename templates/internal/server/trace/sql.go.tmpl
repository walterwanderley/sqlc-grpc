package trace

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/ngrok/sqlmw"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

const driverName = "sqltrace"

func OpenDB(driver driver.Driver, dsn string) (*sql.DB, error) {
	sql.Register(driverName, sqlmw.Driver(driver, new(sqlInterceptor)))
	return sql.Open(driverName, dsn)
}

type sqlInterceptor struct {
	sqlmw.NullInterceptor
}

func (in *sqlInterceptor) ConnExecContext(ctx context.Context, conn driver.ExecerContext, query string, args []driver.NamedValue) (driver.Result, error) {
	ctx, span := otel.Tracer("").Start(ctx, "DB.ConnExecContext")
	defer span.End()
	span.SetAttributes(attribute.String("query", query), attribute.String("query.args", printArgs(args)))

	result, err := conn.ExecContext(ctx, query, args)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return result, err
}

func (in *sqlInterceptor) ConnQueryContext(ctx context.Context, conn driver.QueryerContext, query string, args []driver.NamedValue) (context.Context, driver.Rows, error) {
	ctx, span := otel.Tracer("").Start(ctx, "DB.ConnQueryContext")
	defer span.End()
	span.SetAttributes(attribute.String("query", query), attribute.String("query.args", printArgs(args)))

	rows, err := conn.QueryContext(ctx, query, args)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return ctx, rows, err
}

func printArgs(args []driver.NamedValue) string {
	var b strings.Builder
	for i, a := range args {
		b.WriteString(fmt.Sprintf("%d%s=\"%v\" ", i+1, a.Name, a.Value))
	}
	return b.String()
}
