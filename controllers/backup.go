package controllers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ExecBackup ExecBackup
func ExecBackup(host string, port string, username string, password string, dbname string, encoding string, tipo string, saveFile string) {
	cmd := exec.Command("/usr/bin/pg_dump", "--host="+host, "--port="+port, "--username", username, "--role=postgres", "--no-password", "--format="+tipo, "--blobs", "--encoding="+encoding, "--file="+saveFile, dbname)

	buffer := bytes.Buffer{}
	buffer.Write([]byte(password))
	cmd.Stdin = &buffer

	err := cmd.Run()
	if err != nil {
		fmt.Print("\nDatabase backup error:\t", dbname, "\n")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// ListTables list all tables from server except: 'template0', 'template1', 'postgres'
func ListTables(host string, port string, username string, password string) (tables []string) {
	cmd := exec.Command("/usr/bin/psql", "--host="+host, "--port="+port, "--username="+username, "--no-password", "-t", "-c", "SELECT datname FROM pg_database where datname not in ('template0', 'template1', 'postgres');")

	var inb, outb bytes.Buffer
	inb.Write([]byte(password))

	cmd.Stdout = &outb

	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	listOfTables := strings.Split(outb.String(), "\n")

	for _, table := range listOfTables {
		if table != "" {
			tables = append(tables, strings.TrimSpace(table))
		}
	}

	return tables
}
