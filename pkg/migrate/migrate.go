package migrate

import (
	"log"
	"os"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/olekukonko/tablewriter"
	"github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"

	"personaapp/pkg/flag"
	"personaapp/pkg/postgresql"
)

type Config struct {
	Postgres postgresql.Config
	Limit    int
}

func (c *Config) Flags() *pflag.FlagSet {
	f := pflag.NewFlagSet("ServerConfig", pflag.PanicOnError)

	f.AddFlagSet(c.Postgres.Flags("PostgresConfig", "postgres"))
	f.IntVar(&c.Limit, "environment", 0, "Limit number of migrations to apply")

	return f
}

type migration struct {
	logger *zap.SugaredLogger
	pg     *postgresql.Storage
}

func newMigration(logger *zap.SugaredLogger, pg *postgresql.Storage, migrationTable string) *migration {
	migrate.SetTable(migrationTable)
	return &migration{logger: logger, pg: pg}
}

func (m *migration) migrate(migrations []*migrate.Migration, direction migrate.MigrationDirection, limit int) error {
	applied, err := migrate.ExecMax(m.pg.DB, "postgres", migrate.MemoryMigrationSource{Migrations: migrations}, direction, limit)
	if err != nil {
		return errors.WithStack(err)
	}

	m.logger.Infof("Migrations applied: %d", applied)
	return nil
}

func (m *migration) printStatus(migrations []*migrate.Migration) error {
	const dialect = "postgres"
	records, err := migrate.GetMigrationRecords(m.pg.DB, dialect)
	if err != nil {
		return errors.WithStack(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Migration", "Applied"})
	table.SetColWidth(80)

	type row struct {
		AppliedAt time.Time
	}
	migrationsMap := make(map[string]*row)
	for _, m := range migrations {
		migrationsMap[m.Id] = &row{}
	}

	for _, r := range records {
		migrationsMap[r.Id] = &row{r.AppliedAt}
	}

	for id, m := range migrationsMap {
		appliedAt := "no"
		if !m.AppliedAt.IsZero() {
			appliedAt = m.AppliedAt.UTC().Format(time.RFC3339Nano)
		}
		table.Append([]string{id, appliedAt})
	}
	table.Render()
	return nil
}

func Command(migrationTable string, migrations []*migrate.Migration) *cobra.Command {
	const (
		up     = "up"
		down   = "down"
		status = "status"
	)

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migrations",
	}
	var config Config
	cmd.Flags().AddFlagSet(config.Flags())

	run := func(cmd *cobra.Command, args []string) error {
		if err := flag.BindEnv(cmd); err != nil {
			return errors.WithStack(err)
		}

		logger, _ := zap.NewProduction()
		defer func() {
			err := logger.Sync() // flushes buffer, if any
			if err != nil {
				log.Println(err) // todo think about errors mapper/parser service
			}
		}()
		sugar := logger.Sugar()

		sugar.Info("starting migration")
		defer sugar.Info("stopping migration")

		pg, err := postgresql.New(&config.Postgres)
		if err != nil {
			return errors.WithStack(err)
		}

		m := newMigration(sugar, pg, migrationTable)
		switch cmd.Name() {
		case up:
			return m.migrate(migrations, migrate.Up, config.Limit)
		case down:
			return m.migrate(migrations, migrate.Down, config.Limit)
		case status:
			return m.printStatus(migrations)
		default:
			return errors.Errorf("no handler for command: %s", cmd.Name())
		}
	}

	for _, s := range []*cobra.Command{
		{
			Use:   "up",
			Short: "Up migration",
			RunE:  run,
		},
		{
			Use:   "down",
			Short: "Down migration",
			RunE:  run,
		},
		{
			Use:   "status",
			Short: "Print database status",
			RunE:  run,
		},
	} {
		cmd.AddCommand(s)
	}

	return cmd
}
