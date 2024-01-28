package repository

import (
	"Zhooze/pkg/utils/models"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//	func TestUserSignUp(t *testing.T) {
//		type args struct {
//			input models.UserSignUp
//		}
//		tests := []struct {
//			name       string
//			args       args
//			beforeTest func(sqlmock.Sqlmock)
//			want       models.UserDetailsResponse
//			wantErr    error
//		}{
//			{
//				name: "success signup user",
//				args: args{
//					input: models.UserSignUp{Firstname: "Akhil", Lastname: "c", Email: "zhooze.9550@gmail.com", Phone: "7565748990", Password: "12345"},
//				},
//				beforeTest: func(mockSQL sqlmock.Sqlmock) {
//					expectedQuery := `^INSERT INTO users (.+)$`
//					mockSQL.ExpectQuery(expectedQuery).WithArgs("Akhil", "c", "zhooze.9550@gmail.com", "7565748990", "678790").
//						WillReturnRows(sqlmock.NewRows([]string{"id", "firstname", "lastname", "email", "phone"}).
//							AddRow(1, "Akhil", "c", "zhooze.9550@gmail.com", "12345", "7565748990"))
//				},
//				want:    models.UserDetailsResponse{Id: 1, Firstname: "Akhil", Lastname: "c", Email: "zhooze.9550@gmail.com", Phone: "7565748990"},
//				wantErr: nil,
//			},
//			{
//				name: "error signup user",
//				args: args{
//					input: models.UserSignUp{Firstname: "", Lastname: "", Email: "", Phone: "", Password: ""},
//				},
//				beforeTest: func(mockSQL sqlmock.Sqlmock) {
//					expectedQuery := "Query '(?i)^INSERT INTO users (firstname,lastname,email,password,phone)VALUES($1,$2,$3,$4,$5)RETURNING id,firstname,lastname,email,password,phone$'"
//					mockSQL.ExpectQuery(expectedQuery).WithArgs("Akhil", "c", "zhooze.9550@gmail.com", "678790", "7565748990").
//						WillReturnRows(sqlmock.NewRows([]string{}).AddRow()).
//						WillReturnError(errors.New("email should be unique"))
//				},
//				want:    models.UserDetailsResponse{},
//				wantErr: errors.New("email should be unique"),
//			},
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				mockDB, mockSQL, _ := sqlmock.New()
//				defer mockDB.Close()
//				gormDB, _ := gorm.Open(postgres.New(postgres.Config{
//					Conn: mockDB,
//				}), &gorm.Config{})
//				tt.beforeTest(mockSQL)
//				u := NewUserRepository(gormDB)
//				got, err := u.UserSignUp(tt.args.input)
//				assert.Equal(t, tt.wantErr, err)
//				if !reflect.DeepEqual(got, tt.want) {
//					t.Errorf("userRepo.UserSignUp() = %v, want %v", got, tt.want)
//				}
//			})
//		}
//	}
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
						AddRow(1, 1, "Akhil", "c", "akhilc567@gmail.com", "6282246077", "4321"))
			},

			want: models.UserLoginResponse{
				Id:        1,
				UserId:    1,
				Firstname: "Akhil",
				Lastname:  "c",
				Email:     "akhilc567@gmail.com",
				Phone:     "6282246077",
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
func Test_EditPhone(t *testing.T) {
	tests := []struct {
		name string
		args struct {
			id    int
			phone string
		}
		stub    func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success",
			args: struct {
				id    int
				phone string
			}{id: 1, phone: "9282246077"},
			stub: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectExec("UPDATE users SET phone = ? WHERE id = ?").
					WithArgs("9282246077", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			wantErr: false,
		},
		{
			name: "error",
			args: struct {
				id    int
				phone string
			}{id: 1, phone: "9282246077"},
			stub: func(mockSQL sqlmock.Sqlmock) {
				mockSQL.ExpectExec("UPDATE users SET phone = ? WHERE id = ?").
					WithArgs("9282246077", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			wantErr: true,
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

			err := u.UpdateUserPhone(tt.args.phone, tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("EditPhone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
