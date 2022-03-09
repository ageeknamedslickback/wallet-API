package database_test

import (
	"os"
	"testing"

	"github.com/ageeknamedslickback/wallet-API/wallet/infrastructure/database"
	"github.com/brianvoe/gofakeit/v6"
)

func TestConnectToDatabase(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "happy case",
			wantErr: false,
		},
		{
			name:    "sad case - non existent database",
			wantErr: true,
		},
		{
			name:    "sad case - wrong user password",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad case - non existent database" {
				os.Setenv("DB_NAME", gofakeit.Name())
			}

			if tt.name == "sad case - wrong user password" {
				os.Setenv("DB_PASS", gofakeit.FarmAnimal())
			}

			db, err := database.ConnectToDatabase()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConnectToDatabase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && db == nil {
				t.Fatalf("expected a *gorm.DB object")
			}
		})
	}
}
