package model

import (
	"reflect"
	"testing"
)

func TestNewVoucher(t *testing.T) {
	type args struct {
		balance float32
		userKey string
	}
	tests := []struct {
		name    string
		args    args
		want    Voucher
		wantErr bool
	}{
		{
			name: "fails when user key is not provided",
			args: args{
				balance: 100,
				userKey: "",
			},
			want:    Voucher{},
			wantErr: true,
		},
		{
			name: "creates voucher account successfully",
			args: args{
				balance: 150,
				userKey: "eric-smith",
			},
			want: Voucher{
				Balance:  150,
				Currency: "USD",
				userKey:  "eric-smith",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVoucher(tt.args.balance, tt.args.userKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVoucher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVoucher() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewVoucherDB(t *testing.T) {
	want := &VoucherHandler{
		Account: nil,
		db:      make(map[string]*Voucher),
	}
	got := NewVoucherHandler()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("NewVoucherHandler failed. got: %v, want: %v", got, want)
	}
}

func TestVoucherDB_AddToDB(t *testing.T) {
	type fields struct {
		V  *Voucher
		db map[string]*Voucher
	}

	testDB := getVoucherTestDB()

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "fails when currency is other than the default currency",
			fields: fields{
				V: &Voucher{
					Balance:  150,
					Currency: "EUR",
					userKey:  "eric-smith",
				},
				db: testDB,
			},
			wantErr: true,
		},
		{
			name: "adds new voucher account successfully",
			fields: fields{
				V: &Voucher{
					Balance:  150,
					Currency: "USD",
					userKey:  "eric-smith",
				},
				db: testDB,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VoucherHandler{
				Account: tt.fields.V,
				db:      tt.fields.db,
			}
			if err := v.AddToDB(); (err != nil) != tt.wantErr {
				t.Errorf("AddToDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVoucherDB_Delete(t *testing.T) {
	type fields struct {
		V  *Voucher
		db map[string]*Voucher
	}

	tests := []struct {
		name   string
		fields fields
		key    string
		want   bool
	}{
		{
			name: "returns false when account is not found",
			fields: fields{
				V:  nil,
				db: getVoucherTestDB(),
			},
			key:  "qwerty",
			want: false,
		},
		{
			name: "removes account successfully",
			fields: fields{
				V:  nil,
				db: getVoucherTestDB(),
			},
			key:  "john-doe-usd",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VoucherHandler{
				Account: tt.fields.V,
				db:      tt.fields.db,
			}
			if got := v.Delete(tt.key); got != tt.want {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}

			want := len(getVoucherTestDB()) - 1

			if tt.want && len(tt.fields.db) != want {
				t.Errorf("Delete failed.")
			}
		})
	}
}

func TestVoucherDB_Find(t *testing.T) {
	type fields struct {
		V  *Voucher
		db map[string]*Voucher
	}

	testDB := getVoucherTestDB()

	tests := []struct {
		name    string
		fields  fields
		key     string
		want    *Voucher
		wantErr bool
	}{
		{
			name: "returns error when key is empty",
			fields: fields{
				V:  nil,
				db: testDB,
			},
			key:     "",
			want:    nil,
			wantErr: true,
		},
		{
			name: "returns error when key is a single space",
			fields: fields{
				V:  nil,
				db: testDB,
			},
			key:     " ",
			want:    nil,
			wantErr: true,
		},
		{
			name: "returns error when key is not found",
			fields: fields{
				V:  nil,
				db: testDB,
			},
			key:     "test-user",
			want:    nil,
			wantErr: true,
		},
		{
			name: "finds voucher account successfully",
			fields: fields{
				V:  nil,
				db: testDB,
			},
			key: "john-doe-usd",
			want: &Voucher{
				Balance:  200,
				Currency: DefaultCurrency,
				userKey:  "john-doe",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VoucherHandler{
				Account: tt.fields.V,
				db:      tt.fields.db,
			}
			err := v.Find(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := v.Account
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVoucher_UpdateBalance(t *testing.T) {
	vdb := VoucherHandler{
		Account: &Voucher{
			Balance:  100,
			Currency: DefaultCurrency,
			userKey:  "jane-doe",
		},
		db: getVoucherTestDB(),
	}

	want := float32(105.95)
	got, err := vdb.UpdateBalance(5.95)

	if err != nil || got != want {
		t.Errorf("UpdateBalance failed. Got: %v, want: %v", got, want)
	}
}

func getVoucherTestDB() map[string]*Voucher {
	return map[string]*Voucher{
		"jane-doe-usd": &Voucher{
			Balance:  100,
			Currency: DefaultCurrency,
			userKey:  "jane-doe",
		},
		"john-doe-usd": &Voucher{
			Balance:  200,
			Currency: DefaultCurrency,
			userKey:  "john-doe",
		},
	}
}
