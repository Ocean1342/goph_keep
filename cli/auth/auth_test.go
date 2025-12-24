package auth

import (
	"context"
	"go.uber.org/mock/gomock"
	"goph_keeper/cli/mocks"
	"goph_keeper/cli/remote"
	"goph_keeper/cli/storage"
	"goph_keeper/pkg/auth"
	"testing"
	"time"
)

func TestAuth_Register(t *testing.T) {
	type fields struct {
		store  storage.IStore
		sender remote.IRemote
	}
	type args struct {
		ctx    context.Context
		login  string
		pass   string
		phrase string
	}
	ctx := context.Background()
	token, err := auth.New("secret", time.Hour).
		CreateToken(&auth.User{
			ID: 1, Login: "test", Password: "test", Phrase: "test"},
		)
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			fields: fields{
				store: func() storage.IStore {
					ctrl := gomock.NewController(t)
					m := mocks.NewMockIStore(ctrl)
					m.EXPECT().
						StoreInitialUserData(ctx, "test", "test", "test").
						Return(nil)
					m.EXPECT().
						StoreToken(ctx, "test", token).Return(nil)
					return m
				}(),
				sender: func() remote.IRemote {
					ctrl := gomock.NewController(t)
					m := mocks.NewMockIRemote(ctrl)
					m.EXPECT().
						Register(ctx, "test", "test", "test").
						Return(token, nil)
					return m
				}(),
			},
			args: args{
				ctx:    ctx,
				login:  "test",
				pass:   "test",
				phrase: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Auth{
				store:  tt.fields.store,
				sender: tt.fields.sender,
			}
			if err := a.Register(tt.args.ctx, tt.args.login, tt.args.pass, tt.args.phrase); (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
