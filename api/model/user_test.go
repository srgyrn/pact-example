package model

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewUser(t *testing.T) {
	type args struct {
		name     string
		lastName string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{
			name:    "creates user",
			args:    args{"Eric", "Smith"},
			want:    &User{"Eric", "Smith", 0, nil},
			wantErr: false,
		},
		{
			name:    "fails when name is empty",
			args:    args{"", "smith"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "fails when last name is empty",
			args:    args{"eric", ""},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.name, tt.args.lastName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUserHandler(t *testing.T) {
	got := NewUserHandler()
	want := &UserHandler{nil, make(map[string]*User)}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewUserHandler() = %v, want %v", got, want)
	}
}

func TestUserHandler_AddToDB(t *testing.T) {
	type fields struct {
		Usr *User
		db  map[string]*User
	}

	testDb := getUserTestDb()

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "fails when user exists",
			fields: fields{
				&User{"john", "doe", 100, nil},
				testDb,
			},
			wantErr: true,
		},
		{
			name: "fails when user name is empty",
			fields: fields{
				&User{"", "doe", 100, nil},
				testDb,
			},
			wantErr: true,
		},
		{
			name: "fails when user last name is empty",
			fields: fields{
				&User{"jane", "", 100, nil},
				testDb,
			},
			wantErr: true,
		},
		{
			name: "adds user to db successfully",
			fields: fields{
				&User{"Eric", "Smith", 100, []int{7, 8, 9}},
				testDb,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			udb := &UserHandler{
				Usr: tt.fields.Usr,
				db:  tt.fields.db,
			}
			if err := udb.AddToDB(); (err != nil) != tt.wantErr {
				t.Errorf("AddToDB() error = %v, wantErr %v", err, tt.wantErr)
			}

			want := getUserTestDb()
			want[generateKeyForUser(tt.fields.Usr)] = tt.fields.Usr

			if !tt.wantErr && !reflect.DeepEqual(want, tt.fields.db) {
				t.Errorf("AddToDB() failed. want: %v, got: %v", want, testDb)
			}
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	type fields struct {
		Usr *User
		db  map[string]*User
	}

	testDb := getUserTestDb()

	tests := []struct {
		name   string
		fields fields
		key    string
		want   bool
	}{
		{
			name: "returns false when user not found",
			fields: fields{
				nil,
				testDb,
			},
			key:  "qwerty",
			want: false,
		},
		{
			name: "returns true when user found",
			fields: fields{
				nil,
				testDb,
			},
			key:  "jane-doe",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			udb := &UserHandler{
				Usr: tt.fields.Usr,
				db:  tt.fields.db,
			}
			if got := udb.Delete(tt.key); got != tt.want {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserHandler_Find(t *testing.T) {
	type fields struct {
		Usr *User
		db  map[string]*User
	}

	tests := []struct {
		name    string
		fields  fields
		key     string
		want    interface{}
		wantErr bool
	}{
		{
			name:    "finds user successfully",
			fields:  fields{nil, getUserTestDb()},
			key:     "john-doe",
			want:    &User{"John", "Doe", 100, nil},
			wantErr: false,
		},
		{
			name:    "fails when user does not exist",
			fields:  fields{nil, getUserTestDb()},
			key:     "eric-smith",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			udb := &UserHandler{
				Usr: tt.fields.Usr,
				db:  tt.fields.db,
			}
			got, err := udb.Find(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleUser_UpdateBalance() {
	usr, _ := NewUser("Jane", "Doe")
	got, _ := usr.UpdateBalance(5.95)

	fmt.Printf("%.2f", got)

	// Output:
	// 5.95
}

func getUserTestDb() map[string]*User {
	return map[string]*User{
		"john-doe": &User{"John", "Doe", 100, nil},
		"jane-doe": &User{"Jane", "Doe", 100, []int{1, 2, 3}},
	}
}
