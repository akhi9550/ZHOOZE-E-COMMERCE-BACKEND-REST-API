package repository

import (
	"Zhooze/pkg/utils/models"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_GetUserDetails(t *testing.T) {
	tests := []struct {
		name    string
		args    int
		stub    func(mockSQL sqlmock.Sqlmock)
		want    models.UsersProfileDetails
		wantErr error
	}{
		{
			name: "success",
			args: 1,
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectQuery := `SELECT u.firstname,u.lastname,u.email,u.phone FROM users u WHERE u.id = ?`
				mockSQL.ExpectQuery(expectQuery).WillReturnRows(sqlmock.NewRows([]string{"firstname", "lastname", "email", "phone"}).AddRow("akhil", "c", "akhil89@gmail.com", "+919087678564"))
			},
			want: models.UsersProfileDetails{
				Firstname: "akhil",
				Lastname:  "c",
				Email:     "akhil89@gmail.com",
				Phone:     "+919087678564",
			},
			wantErr: nil,
		},
		{
			name: "error",
			args: 1,
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectQuery := `SELECT u.firstname,u.lastname,u.email,u.phone FROM users u WHERE u.id = ?`
				mockSQL.ExpectQuery(expectQuery).WillReturnError(errors.New("error"))
			},
			want:    models.UsersProfileDetails{},
			wantErr: errors.New("could not get user details"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()
			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})
			tt.stub(mockSQL)
			u := NewUserRepository(gormDB)

			result, err := u.UserDetails(tt.args)

			assert.Equal(t, tt.want, result)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestUserSignUp(t *testing.T) {
	type args struct {
		input models.UserSignUp
	}
	tests := []struct {
		name       string
		args       args
		beforeTest func(mockSQL sqlmock.Sqlmock)
		want       models.UserDetailsResponse
		wantErr    error
	}{
		{
			name: "success signup user",
			args: args{
				input: models.UserSignUp{
					Firstname: "Akhil",
					Lastname:  "c",
					Email:     "zhooze9550@gmail.com",
					Password:  "12345",
					Phone:     "+917565748990",
				},
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `INSERT INTO users\(firstname, lastname, email, password, phone\) VALUES\(\$1, \$2, \$3, \$4, \$5\) RETURNING id, firstname, lastname, email, phone`
				mockSQL.ExpectQuery(expectedQuery).
					WithArgs("Akhil", "c", "zhooze9550@gmail.com", "12345", "+917565748990").
					WillReturnRows(sqlmock.NewRows([]string{"id", "firstname", "lastname", "email", "phone"}).
						AddRow(1, "Akhil", "c", "zhooze9550@gmail.com", "+917565748990"))
			},
			want: models.UserDetailsResponse{
				Id:        1,
				Firstname: "Akhil",
				Lastname:  "c",
				Email:     "zhooze9550@gmail.com",
				Phone:     "+917565748990",
			},
			wantErr: nil,
		},
		{
			name: "error signup user",
			args: args{
				input: models.UserSignUp{},
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `INSERT INTO users\(firstname, lastname, email, password, phone\) VALUES\(\$1, \$2, \$3, \$4, \$5\) RETURNING id, firstname, lastname, email, phone`
				mockSQL.ExpectQuery(expectedQuery).
					WithArgs("", "", "", "", "").
					WillReturnError(errors.New("email should be unique"))
			},

			want:    models.UserDetailsResponse{},
			wantErr: errors.New("email should be unique"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()
			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})
			tt.beforeTest(mockSQL)
			u := NewUserRepository(gormDB)
			got, err := u.UserSignUp(tt.args.input)
			assert.Equal(t, tt.wantErr, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.UserSignUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FindUserByEmail(t *testing.T) {
	tests := []struct {
		name    string
		args    models.LoginDetail
		stub    func(sqlmock.Sqlmock)
		want    models.UserLoginResponse
		wantErr error
	}{
		{
			name: "success",
			args: models.LoginDetail{
				Email:    "akhilc567@gmail.com",
				Password: "1234",
			},
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `^SELECT \* FROM users(.+)$`
				mockSQL.ExpectQuery(expectedQuery).WithArgs("akhilc567@gmail.com").
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "firstname", "lastname", "email", "phone", "password"}).
						AddRow(1, 1, "Akhil", "c", "akhilc567@gmail.com", "+916282246077", "4321"))
			},

			want: models.UserLoginResponse{
				Id:        1,
				UserId:    1,
				Firstname: "Akhil",
				Lastname:  "c",
				Email:     "akhilc567@gmail.com",
				Phone:     "+916282246077",
				Password:  "4321",
			},
			wantErr: nil,
		},
		{
			name: "error",
			args: models.LoginDetail{
				Email:    "akhilc567@gmail.com",
				Password: "1234",
			},
			stub: func(mockSQL sqlmock.Sqlmock) {

				expectedQuery := `^SELECT \* FROM users(.+)$`

				mockSQL.ExpectQuery(expectedQuery).WithArgs("akhilc567@gmail.com").
					WillReturnError(errors.New("new error"))

			},

			want:    models.UserLoginResponse{},
			wantErr: errors.New("error checking user details"),
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()

			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})

			tt.stub(mockSQL)

			u := NewUserRepository(gormDB)

			result, err := u.FindUserByEmail(tt.args)

			assert.Equal(t, tt.want, result)
			assert.Equal(t, tt.wantErr, err)
		})
	}

}
func Test_FindIdFromPhone(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		stub    func(sqlmock.Sqlmock)
		want    int
		wantErr error
	}{
		{
			name: "success",
			args: "+919087678909",
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `SELECT id FROM users WHERE phone=\$1`
				mockSQL.ExpectQuery(expectedQuery).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			want:    1,
			wantErr: nil,
		},
		{
			name: "error",
			args: "+919087678909",
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `SELECT id FROM users WHERE phone=\$1`
				mockSQL.ExpectQuery(expectedQuery).
					WillReturnError(errors.New("error"))
			},
			want:    0,
			wantErr: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()

			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})

			tt.stub(mockSQL)
			u := NewUserRepository(gormDB)

			result, err := u.FindIdFromPhone(tt.args)
			assert.Equal(t, tt.want, result)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
func Test_GetAllAddress(t *testing.T) {
	testCase := []struct {
		name    string
		args    int
		stub    func(mockSQL sqlmock.Sqlmock)
		want    models.AddressInfoResponse
		wantErr error
	}{
		{
			name: "Success",
			args: 1,
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `SELECT \* FROM addresses WHERE user_id = \$1`
				mockSQL.ExpectQuery(expectedQuery).WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "house_name", "street", "city", "state", "pin"}).
						AddRow(1, "akhil", "chekkiyil house", "cheleri", "kannur", "kerala", "670604"))
			},
			want: models.AddressInfoResponse{
				ID:        1,
				Name:      "akhil",
				HouseName: "chekkiyil house",
				Street:    "cheleri",
				City:      "kannur",
				State:     "kerala",
				Pin:       "670604",
			},
			wantErr: nil,
		},
		{
			name: "failed",
			args: 1,
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `SELECT \* FROM addresses WHERE user_id = \$1`
				mockSQL.ExpectQuery(expectedQuery).WithArgs(1).
					WillReturnError(errors.New("error"))
			},
			want:    models.AddressInfoResponse{},
			wantErr: errors.New("error"),
		},
	}
	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()
			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})
			tt.stub(mockSQL)
			u := NewUserRepository(gormDB)
			result, err := u.GetAllAddres(tt.args)
			assert.Equal(t, tt.want, result)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
