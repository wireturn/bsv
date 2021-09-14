package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ConfigValidate(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		cfg *Config
		err error
	}{
		"valid db config (sqlite) should return no errors": {
			cfg: &Config{
				Db: &Db{
					Type: "sqlite",
				},
			},
			err: nil,
		}, "valid db config (postgres) should return no errors": {
			cfg: &Config{
				Db: &Db{
					Type: "postgres",
				},
			},
			err: nil,
		}, "valid db config (mysql) should return no errors": {
			cfg: &Config{
				Db: &Db{
					Type: "mysql",
				},
			},
			err: nil,
		}, "invalid db config should return no errors": {
			cfg: &Config{
				Db: &Db{
					Type: "mydb",
				},
			},
			err: errors.New("[db.type: value mydb failed to meet requirements]"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.cfg.Validate()
			if test.err == nil {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, test.err.Error())

		})
	}
}
