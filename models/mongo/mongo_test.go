package mongo_test

import (
	"testing"

	"github.com/ocakhasan/getir-api-task/models/mongo"

	"github.com/ocakhasan/getir-api-task/controllers/errors"

	"github.com/ocakhasan/getir-api-task/controllers/requests"
)

func TestCreatePipe(t *testing.T) {
	type args struct {
		filter requests.MongoRequestBody
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "start date is later than end date",
			args: args{
				filter: requests.MongoRequestBody{
					StartDate: "2020-01-02",
					EndDate:   "2019-01-02",
					MinCount:  2000,
					MaxCount:  3000,
				},
			},
			wantErr: errors.ErrorInvalidBody,
		},
		{
			name: "min count is higher than max count",
			args: args{
				filter: requests.MongoRequestBody{
					StartDate: "2020-01-02",
					EndDate:   "2021-01-02",
					MinCount:  3000,
					MaxCount:  2000,
				},
			},
			wantErr: errors.ErrorInvalidBody,
		},
		{
			name: "everything is perfect",
			args: args{
				filter: requests.MongoRequestBody{
					StartDate: "2020-01-02",
					EndDate:   "2021-01-02",
					MinCount:  2000,
					MaxCount:  3000,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mongo.CreatePipe(tt.args.filter)
			if err != tt.wantErr {
				t.Errorf("CreatePipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
