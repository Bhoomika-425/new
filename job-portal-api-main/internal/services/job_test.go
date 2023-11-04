package service

import (
	"context"

	"project/internal/auth"
	"project/internal/models"
	"project/internal/repository"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestService_ViewJobById(t *testing.T) {
	type args struct {
		ctx context.Context
		jid uint64
	}
	tests := []struct {
		name string
		// s       *Service
		args             args
		want             models.Jobs
		wantErr          bool
		mockRepoResponse func() (models.Jobs, error)
	}{
		// {name: "error",
		// 	want: models.Jobs{},
		// 	args: args{
		// 		ctx: context.Background(),
		// 		jid: 5,
		// 	},
		// 	wantErr: true,
		// 	mockRepoResponse: func() (models.Jobs, error) {
		// 		return models.Jobs{}, errors.New("test error")
		// 	},
		// },
		{name: "success",
			want: models.Jobs{
				Company: models.Company{
					Name:     "tcs",
					Location: "bang",
					Field:    "software",
				},
				Cid:          2,
				Name:         "developer",
				Salary:       "30000",
				NoticePeriod: "3 weeks",
			},
			args: args{
				ctx: context.Background(),
				jid: 15,
			},
			wantErr: false,
			mockRepoResponse: func() (models.Jobs, error) {
				return models.Jobs{
					Company: models.Company{
						Name:     "tcs",
						Location: "bang",
						Field:    "software",
					},
					Cid:          2,
					Name:         "developer",
					Salary:       "30000",
					NoticePeriod: "3 weeks",
				}, nil
			},
		},
		{
			name: "invalid id",
			want: models.Jobs{},
			args: args{
				ctx: context.Background(),
				jid: 5,
			},
			mockRepoResponse: func() (models.Jobs, error) {
				return models.Jobs{}, nil
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().Jobbyjid(tt.args.ctx, tt.args.jid).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(mockRepo, &auth.Auth{})
			got, err := s.ViewJobById(tt.args.ctx, tt.args.jid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ViewJobById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ViewJobById() = %v, want %v", got, tt.want)
			}
		})
	}
}
