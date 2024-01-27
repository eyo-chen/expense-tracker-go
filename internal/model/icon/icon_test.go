package icon

import (
	"database/sql"
	"testing"

	"github.com/OYE0303/expense-tracker-go/internal/domain"
	"github.com/OYE0303/expense-tracker-go/internal/usecase"
	"github.com/OYE0303/expense-tracker-go/pkg/dockerutil"
	"github.com/OYE0303/expense-tracker-go/pkg/testutil"
	"github.com/stretchr/testify/suite"
)

type IconSuite struct {
	suite.Suite
	db    *sql.DB
	model usecase.IconModel
}

func TestIconSuite(t *testing.T) {
	suite.Run(t, new(IconSuite))
}

func (s *IconSuite) SetupSuite() {
	port := dockerutil.RunDocker()
	db := testutil.ConnToDB(port)
	s.model = NewIconModel(db)
	s.db = db
}

func (s *IconSuite) TearDownSuite() {
	s.db.Close()
	dockerutil.PurgeDocker()
}

func (s *IconSuite) TestGetByID() {
	tests := []struct {
		Desc        string
		ID          int64
		SetupFun    func() error
		Expected    *domain.Icon
		ExpectedErr error
	}{
		{
			Desc: "Get icon successfully",
			ID:   1,
			SetupFun: func() error {
				stmt := `INSERT INTO icons (url) VALUES (?)`
				_, err := s.db.Exec(stmt, "https://test.com")
				return err
			},
			Expected: &domain.Icon{
				ID:  1,
				URL: "https://test.com",
			},
			ExpectedErr: nil,
		},
		{
			Desc:        "Not found",
			ID:          2,
			Expected:    nil,
			ExpectedErr: domain.ErrIconNotFound,
		},
	}

	for _, test := range tests {
		s.T().Run(test.Desc, func(t *testing.T) {
			if test.SetupFun != nil {
				err := test.SetupFun()
				s.NoError(err, test.Desc)
			}

			icon, err := s.model.GetByID(test.ID)
			if test.ExpectedErr != nil {
				s.EqualError(err, test.ExpectedErr.Error(), test.Desc)
				return
			}

			s.Equal(test.Expected, icon, test.Desc)
			s.Equal(test.Expected.ID, icon.ID, test.Desc)
		})
	}
}
