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
		name           string
		fields         fields
		want           NoteConfig
		wantErr        bool
		AssertCallBack func()
	}{
		{
			name:    "Happy Path",
			fields:  fields{fileHelper: happyPathFileHelper},
			want:    defaultConfig,
			wantErr: false,
			AssertCallBack: func() {
				happyPathFileHelper.AssertCalled(t, "ReadFile", ConfigPath)
			},
		},
		{
			name:    "Missing file",
			fields:  fields{fileHelper: missingFileFileHelper},
			want:    NoteConfig{},
			wantErr: true,
			AssertCallBack: func() {
				missingFileFileHelper.AssertCalled(t, "ReadFile", ConfigPath)
			},
		},
		{
			name:    "Invalid json",
			fields:  fields{fileHelper: invalidJsonFileHelper},
			want:    NoteConfig{},
			wantErr: true,
			AssertCallBack: func() {
				invalidJsonFileHelper.AssertCalled(t, "ReadFile", ConfigPath)
			},
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
			tt.AssertCallBack()
		})
	}
}

func Test_configHelper_Setup(t *testing.T) {
	config, _ := json.MarshalIndent(defaultConfig, "", "  ")
	missingFileHelper := &mocks.FileHelper{}
	missingFileHelper.On("ReadFile", ConfigPath).Return(nil, errors.New("missing file"))
	missingFileHelper.On("WriteFile", mock.Anything, mock.Anything).Return(nil, nil)
	missingFileHelper.On("AppendHomeDirectory", ConfigPath).Return(ConfigPath, nil)
	missingFileHelper.On("EditorOpenFile", mock.Anything, mock.Anything).Return(nil)

	mockFileHelper := &mocks.FileHelper{}
	mockFileHelper.On("ReadFile", ConfigPath).Return(config, nil)
	mockFileHelper.On("AppendHomeDirectory", ConfigPath).Return(ConfigPath, nil)
	mockFileHelper.On("EditorOpenFile", mock.Anything, mock.Anything).Return(nil)

	type fields struct {
		fileHelper FileHelper
	}

	tests := []struct {
		name           string
		fields         fields
		wantErr        bool
		AssertCallBack func()
	}{
		{
			name:    "Config file does not exist, write default",
			fields:  fields{fileHelper: missingFileHelper},
			wantErr: false,
			AssertCallBack: func() {
				missingFileHelper.AssertCalled(t, "WriteFile", ConfigPath, config)
				missingFileHelper.AssertCalled(t, "WriteFile", defaultConfig.Template, mock.Anything)
				missingFileHelper.AssertCalled(t, "EditorOpenFile", defaultConfig.Editor, ConfigPath)
			},
		},
		{
			name:    "Config file exist, do nothing just open it",
			fields:  fields{fileHelper: mockFileHelper},
			wantErr: false,
			AssertCallBack: func() {
				mockFileHelper.AssertNotCalled(t, "WriteFile", mock.Anything, mock.Anything)
				mockFileHelper.AssertCalled(t, "EditorOpenFile", defaultConfig.Editor, ConfigPath)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configHelper := configHelper{
				fileHelper: tt.fields.fileHelper,
			}
			if err := configHelper.Setup(); (err != nil) != tt.wantErr {
				t.Errorf("Setup() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.AssertCallBack()
		})
	}
}
