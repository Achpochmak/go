// + build suite

package suite_tests

import (
	"HOMEWORK-1/internal/cli"
	"HOMEWORK-1/internal/module"
	"HOMEWORK-1/internal/repository"
	"HOMEWORK-1/internal/repository/transactor"
	"HOMEWORK-1/tests"

	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type PVZTestSuite struct {
	suite.Suite
	ctx  context.Context
	tm   *transactor.TransactionManager
	repo *repository.Repository
	pvz  module.Module
	cli  *cli.CLI
	pool *pgxpool.Pool
}

func TestPVZ(t *testing.T) {
	suite.Run(t, new(PVZTestSuite))
}

func (s *PVZTestSuite) SetupSuite() {
	tests.InitConfig()
	s.pool = tests.ConnectDB()
	s.tm = &transactor.TransactionManager{Pool: s.pool}
	s.repo = repository.NewRepository(s.tm)
	s.pvz = module.NewModule(module.Deps{
		Repository: s.repo,
		Transactor: s.tm,
	})

	s.ctx = context.Background()
	s.cli = cli.NewCLI(cli.Deps{Module: s.pvz})
}

func (s *PVZTestSuite) TearDownSuite() {
	s.pool.Close()
}

func (s *PVZTestSuite) SetupTest() {
	s.repo.TruncateTable(s.ctx)
}

func (s *PVZTestSuite) TearDownTest() {
	s.repo.TruncateTable(s.ctx)
}
