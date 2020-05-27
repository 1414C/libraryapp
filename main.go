package main

import (
	"flag"

	"github.com/1414C/libraryapp/appobj"

	_ "github.com/SAP/go-hdb/driver"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	prodFlag := flag.Bool("prod", false, "this flag should be set in a production environment")
	devFlag := flag.Bool("dev", false, "this flag should be set in a development environment")
	drFlag := flag.Bool("dr", false, "db destructive reset")
	rsFlag := flag.Bool("rs", false, "rebuild the Auth allocations to the Super UsrGroup")
	flag.Parse()

	a := appobj.AppObj{}
	a.Initialize(*devFlag, *prodFlag, *drFlag, *rsFlag)

	lsg := a.CreateLeadSetGet()
	a.Run(lsg)
}
