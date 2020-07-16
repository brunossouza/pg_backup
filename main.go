package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"pg_backup/controllers"
	"runtime"
)

var (
	cfg     controllers.Config
	verbose bool
)

func helpMenu() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Programa com a finalidade de facilitar a realização do backup do banco de dados Postgresql."+
			"\n\nAjuda:\n\t-h\tExibi está ajuda\n")

		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stderr, "\t-%v\t%v\n", f.Name, f.Usage) // f.Name, f.Value
		})

		fmt.Println()
		os.Exit(0)
	}

	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.Parse()
}

func init() {
	system := runtime.GOOS

	if system == "linux" {
		if !controllers.FileExists("/usr/bin/pg_dump") {
			log.Fatalln("commando não encontrado /usr/bin/pg_dump, necessario à execução")
		}
		if !controllers.FileExists("/usr/bin/psql") {
			log.Fatalln("commando não encontrado /usr/bin/psql, necessario à execução")
		}
	}

	helpMenu()

	controllers.ReadFile(&cfg)
	controllers.CreateBackupDirectory(cfg.Backup.Directory)

}

func main() {
	var tables []string

	backupFormats := map[string]string{
		"plain":     "sql",
		"custom":    "backup",
		"tar":       "tar",
		"directory": "directory",
	}

	if cfg.Backup.Databases[0] == "all" {
		tables = controllers.ListTables(cfg.Server.Host, cfg.Server.Port, cfg.Database.Username, cfg.Database.Password)
	} else {
		tables = cfg.Backup.Databases
	}

	if verbose {
		fmt.Printf("Total de bancos encontradas:\t%v\n\n", len(tables))
	}

	for _, table := range tables {
		var path = fmt.Sprintf("%s/%s.%s", cfg.Backup.Directory, table, backupFormats[cfg.Backup.BackupFormat])
		if verbose {
			fmt.Printf("Banco:\t%s\nSalvo em:\t%s\n\n", table, path)
		}
		controllers.ExecBackup(cfg.Server.Host, cfg.Server.Port, cfg.Database.Username, cfg.Database.Password, table, cfg.Backup.BackupFormat, path)
	}
}
