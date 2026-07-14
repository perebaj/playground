// quack-probe — debugging only: inspect what's loaded in go-duckdb's
// embedded DuckDB after attempting to install + load Quack.
package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/marcboeker/go-duckdb/v2"
)

func main() {
	db, err := sql.Open("duckdb", "?allow_unsigned_extensions=true")
	must(err)
	defer db.Close()

	for _, s := range []string{
		"INSTALL quack FROM 'http://extensions.duckdb.org'",
		"LOAD quack",
	} {
		fmt.Printf("\n>> %s\n", s)
		if _, err := db.ExecContext(context.Background(), s); err != nil {
			fmt.Printf("   ERR: %v\n", err)
		} else {
			fmt.Println("   ok")
		}
	}

	rows, err := db.QueryContext(context.Background(),
		"SELECT extension_name, loaded, installed, install_mode, extension_version FROM duckdb_extensions() WHERE extension_name = 'quack'")
	must(err)
	defer rows.Close()

	cols, _ := rows.Columns()
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		must(rows.Scan(ptrs...))
		fmt.Println("\nQuack row:")
		for i, c := range cols {
			fmt.Printf("  %s = %v\n", c, vals[i])
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
