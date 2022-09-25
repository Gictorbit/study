package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

const query = `
SELECT
    user_id,role_ids,ip,log_name
FROM
    auditlog_tbl
WHERE
    user_id!=0 AND company_id=$1
ORDER BY
    registered_at DESC
LIMIT $2;
`

type AuditLogg struct {
	UserID    uint32
	Roles     []uint32
	IP        string
	RpcMethod string
}

var dialCount int

func main() {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:8989"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "admin",
			Password: "admin",
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			dialCount++
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: false,
		Debugf: func(format string, v ...interface{}) {
			log.Println(format, v)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:      time.Duration(10) * time.Second,
		MaxOpenConns:     5,
		MaxIdleConns:     5,
		ConnMaxLifetime:  time.Duration(10) * time.Minute,
		ConnOpenStrategy: clickhouse.ConnOpenInOrder,
	})
	if err != nil {
		panic(err)
	}
	for i := 0; i < 3; i++ {
		audits, err := getAudits(conn, 2, 20)
		if err != nil {
			panic(err)
		}
		for _, aud := range audits {
			fmt.Println(aud.UserID, aud.IP, aud.Roles, aud.RpcMethod)
		}
		fmt.Println("###########################")
	}
	fmt.Println("dialcount", dialCount)
}

func getAudits(conn driver.Conn, companyID, limit uint32) ([]*AuditLogg, error) {
	rows, err := conn.Query(context.Background(), query, companyID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	audits := make([]*AuditLogg, 0)
	for rows.Next() {
		audit := &AuditLogg{}
		err := rows.Scan(
			&audit.UserID,
			&audit.Roles,
			&audit.IP,
			&audit.RpcMethod,
		)
		if err != nil {
			return nil, err
		}
		audits = append(audits, audit)
	}
	return audits, nil
}
