package model

import (
	"reflect"
	"testing"
)

func TestOrderDB_AddToDB(t *testing.T) {
	type fields struct {
		Ord Order
		db  map[string]*Order
	}

	testDB := getOrderTestDb()

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "fails when payment way is missing",
			fields: fields{
				Ord: Order{
					ID:                  5,
					Total:               200,
					PaymentWay:          0,
					ShippingCountryZone: ZoneEurope,
					IsDeleted:           false,
				},
				db: testDB,
			},
			wantErr: true,
		},
		{
			name: "fails when country zone is missing",
			fields: fields{
				Ord: Order{
					ID:                  5,
					Total:               200,
					PaymentWay:          CreditCard,
					ShippingCountryZone: 0,
					IsDeleted:           false,
				},
				db: testDB,
			},
			wantErr: true,
		},
		{
			name: "adds order successfully",
			fields: fields{
				Ord: Order{
					ID:                  5,
					Total:               200,
					PaymentWay:          CreditCard,
					ShippingCountryZone: ZoneEurope,
					IsDeleted:           false,
				},
				db: testDB,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderDB{
				Ord: tt.fields.Ord,
				db:  tt.fields.db,
			}
			if err := o.AddToDB(); (err != nil) != tt.wantErr {
				t.Errorf("AddToDB() error = %v, wantErr %v", err, tt.wantErr)
			}

			want := getOrderTestDb()
			want[string(len(want)+1)] = &tt.fields.Ord

			if !tt.wantErr && !reflect.DeepEqual(want, tt.fields.db) {
				t.Errorf("AddToDB failed, want: %v, got: %v", want, tt.fields.db)
			}
		})
	}
}

func TestOrderDB_Delete(t *testing.T) {
	type fields struct {
		Ord Order
		db  map[string]*Order
	}

	tests := []struct {
		name   string
		fields fields
		key    string
		want   bool
	}{
		{
			name: "return false when order not found",
			fields: fields{
				Ord: Order{},
				db:  getOrderTestDb(),
			},
			key:  "1000",
			want: false,
		},
		{
			name: "returns true when order is found",
			fields: fields{
				Ord: Order{},
				db:  getOrderTestDb(),
			},
			key:  "1",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderDB{
				Ord: tt.fields.Ord,
				db:  tt.fields.db,
			}
			if got := o.Delete(tt.key); got != tt.want {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}

			ord := tt.fields.db[tt.key]

			if tt.want && !ord.IsDeleted {
				t.Errorf("Delete failed, want %v, got %v", true, ord.IsDeleted)
			}
		})
	}
}

func TestOrderDB_Find(t *testing.T) {
	type fields struct {
		Ord Order
		db  map[string]*Order
	}

	testDB := getOrderTestDb()

	tests := []struct {
		name    string
		fields  fields
		key     string
		want    interface{}
		wantErr bool
	}{
		{
			name: "fails when key is not found",
			fields: fields{
				Ord: Order{},
				db:  testDB,
			},
			key:     "1000",
			want:    nil,
			wantErr: true,
		},
		{
			name: "returns order when found",
			fields: fields{
				Ord: Order{},
				db:  testDB,
			},
			key: "1",
			want: Order{
				ID:                  1,
				Total:               100,
				PaymentWay:          CreditCard,
				ShippingCountryZone: ZoneEurope,
				IsDeleted:           false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderDB{
				Ord: tt.fields.Ord,
				db:  tt.fields.db,
			}
			got, err := o.Find(tt.key)
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

func getOrderTestDb() map[string]*Order {
	return map[string]*Order{
		"1": &Order{
			ID:                  1,
			Total:               100,
			PaymentWay:          CreditCard,
			ShippingCountryZone: ZoneEurope,
			IsDeleted:           false,
		},
		"2": &Order{
			ID:                  2,
			Total:               200,
			PaymentWay:          CashOnDelivery,
			ShippingCountryZone: ZoneAmerica,
			IsDeleted:           true,
		},
		"3": &Order{
			ID:                  3,
			Total:               300,
			PaymentWay:          Paypal,
			ShippingCountryZone: ZoneMena,
			IsDeleted:           false,
		},
		"4": &Order{
			ID:                  4,
			Total:               400,
			PaymentWay:          CashOnDelivery,
			ShippingCountryZone: ZoneMena,
			IsDeleted:           false,
		},
	}
}
