package main

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/mock"
	"note/mocks"
	"reflect"
	"testing"
)

func Test_configHelper_ReadConfig(t *testing.T) {
	happyPathFileHelper := &mocks.FileHelper{}
	config, _ := json.Marshal(defaultConfig)
	happyPathFileHelper.On("ReadFile", mock.Anything).Return(config, nil)
	missingFileFileHelper := &mocks.FileHelper{}
	missingFileFileHelper.On("ReadFile", mock.Anything).Return(nil, errors.New("missing file"))
	invalidJsonFileHelper := &mocks.FileHelper{}
	invalidJsonFileHelper.On("ReadFile", mock.Anything).Return([]byte("not json"), nil)
	type fields struct {
		fileHelper FileHelper
	}
	tests := []struct {
		name    string
		fields  fields
		want    NoteConfig
		wantErr bool
	}{
		{
			name:    "Happy Path",
			fields:  struct{ fileHelper FileHelper }{fileHelper: happyPathFileHelper},
			want:    defaultConfig,
			wantErr: false,
		},
		{
			name:    "Missing file",
			fields:  struct{ fileHelper FileHelper }{fileHelper: missingFileFileHelper},
			want:    NoteConfig{},
			wantErr: true,
		},
		{
			name:    "Invalid json",
			fields:  struct{ fileHelper FileHelper }{fileHelper: invalidJsonFileHelper},
			want:    NoteConfig{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configHelper := configHelper{
				fileHelper: tt.fields.fileHelper,
			}
			got, err := configHelper.ReadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
