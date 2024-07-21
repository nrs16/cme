package server

import (
	"bufio"
	"fmt"
	"log"
	"nrs16/cme/config"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocql/gocql"
)

func InitializeDatabase(conf config.Config) {
	fmt.Println("INITIALIZING DB")
	cluster := gocql.NewCluster(conf.Database.Host)
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	defer session.Close()

	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalf("error getting file executable path: %s", err.Error())
	}
	executableDir := filepath.Dir(executablePath)

	if !keyspaceExists(session, conf.Database.KeySpace) {
		err := executeCQLFile(session, executableDir+"/schema.cql")
		if err != nil {
			log.Fatalf("Failed to execute CQL file: %v", err)
		}
	}
	// Your application logic here
}

func executeCQLFile(session *gocql.Session, filePath string) error {

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var sb strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "--") || len(strings.TrimSpace(line)) == 0 {
			continue
		}
		fmt.Println(line)
		sb.WriteString(line)
		if strings.HasSuffix(line, ";") {
			err = session.Query(sb.String()).Exec()
			if err != nil {
				return fmt.Errorf("failed to execute query: %w", err)
			}
			sb.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return nil
}

func keyspaceExists(session *gocql.Session, keyspace string) bool {
	query := fmt.Sprintf("SELECT keyspace_name FROM system_schema.keyspaces WHERE keyspace_name='%s'", keyspace)
	var ks string
	if err := session.Query(query).Scan(&ks); err != nil {
		if err == gocql.ErrNotFound {
			return false
		}
		log.Fatalf("Failed to query keyspaces: %v", err)
	}
	return ks == keyspace
}
