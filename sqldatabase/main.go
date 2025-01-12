package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools/sqldatabase"
	"github.com/tmc/langchaingo/tools/sqldatabase/sqlite3"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	llm, err := openai.New()
	if err != nil {
		return err
	}

	const dns = "./data.db"
	db, err := sqldatabase.NewSQLDatabaseWithDSN(sqlite3.EngineName, dns, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	chain := chains.NewSQLDatabaseChain(llm, 100, db)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Waiting for input >\n")
	for scanner.Scan() {
		in := scanner.Text()
		out, err := chains.Run(ctx, chain, in)
		if err != nil {
			return err
		}
		fmt.Print("\n< Answer from the system\n")
		fmt.Print(out + "\n\n")
		fmt.Print("Waiting for input >\n")
	}

	return nil
}
