package main

import (
	"github.com/stretchr/testify/mock"
	"note/mocks"
	"testing"
	"time"
)

func Test_noteHelper_OpenNote(t *testing.T) {
	Now = func() time.Time {
		return time.Date(2022, time.October, 1, 10, 0, 0, 0, time.UTC)
	}
	defer func() { Now = time.Now }()
	noteExistsHelper := &mocks.FileHelper{}
	noteExistsHelper.On("ReadFile", mock.Anything).Return([]byte("notefile"), nil)
	noteExistsHelper.On("FileExists", mock.Anything).Return(true, nil)
	noteExistsHelper.On("AppendHomeDirectory", mock.Anything).Return(func(input string) (string, error) {
		return input, nil
	})
	noteExistsHelper.On("EditorOpenFile", mock.Anything, mock.Anything).Return(nil)

	noteMissingHelper := &mocks.FileHelper{}
	noteMissingHelper.On("ReadFile", mock.Anything).Return([]byte("notefile"), nil)
	noteMissingHelper.On("FileExists", mock.Anything).Return(false, nil)
	noteMissingHelper.On("AppendHomeDirectory", mock.Anything).Return(func(input string) (string, error) {
		return input, nil
	})
	noteMissingHelper.On("EditorOpenFile", mock.Anything, mock.Anything).Return(nil)
	noteMissingHelper.On("WriteFile", mock.Anything, mock.Anything).Return(nil)

	type fields struct {
		fileHelper FileHelper
	}
	type args struct {
		relativeWeek int
		config       NoteConfig
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		AssertCallBack func()
	}{
		{
			name:    "Note exists, relativeWeek: 0",
			fields:  fields{fileHelper: noteExistsHelper},
			args:    args{relativeWeek: 0, config: defaultConfig},
			wantErr: false,
			AssertCallBack: func() {
				noteExistsHelper.AssertCalled(t, "EditorOpenFile", defaultConfig.Editor, "~/notes/2022-09-26."+defaultConfig.Extension)
			},
		},
		{
			name:    "Note exists, relativeWeek: 1",
			fields:  fields{fileHelper: noteExistsHelper},
			args:    args{relativeWeek: 1, config: defaultConfig},
			wantErr: false,
			AssertCallBack: func() {
				noteExistsHelper.AssertCalled(t, "EditorOpenFile", defaultConfig.Editor, "~/notes/2022-10-03."+defaultConfig.Extension)
			},
		},
		{
			name:    "Note exists, relativeWeek: -1",
			fields:  fields{fileHelper: noteExistsHelper},
			args:    args{relativeWeek: -1, config: defaultConfig},
			wantErr: false,
			AssertCallBack: func() {
				noteExistsHelper.AssertCalled(t, "EditorOpenFile", defaultConfig.Editor, "~/notes/2022-09-19."+defaultConfig.Extension)
			},
		},
		{
			name:    "Note exists, relativeWeek: 0",
			fields:  fields{fileHelper: noteMissingHelper},
			args:    args{relativeWeek: 0, config: defaultConfig},
			wantErr: false,
			AssertCallBack: func() {
				noteMissingHelper.AssertCalled(t, "WriteFile", defaultConfig.Location+"/2022-09-26."+defaultConfig.Extension, mock.Anything)
				noteMissingHelper.AssertCalled(t, "EditorOpenFile", defaultConfig.Editor, "~/notes/2022-09-26."+defaultConfig.Extension)
			},
		},
		{
			name:    "Note exists, relativeWeek: 1",
			fields:  fields{fileHelper: noteMissingHelper},
			args:    args{relativeWeek: 1, config: defaultConfig},
			wantErr: false,
			AssertCallBack: func() {
				noteMissingHelper.AssertCalled(t, "WriteFile", defaultConfig.Location+"/2022-10-03."+defaultConfig.Extension, mock.Anything)
				noteMissingHelper.AssertCalled(t, "EditorOpenFile", defaultConfig.Editor, "~/notes/2022-10-03."+defaultConfig.Extension)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			noteHelper := noteHelper{
				fileHelper: tt.fields.fileHelper,
			}
			if err := noteHelper.OpenNote(tt.args.relativeWeek, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("OpenNote() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.AssertCallBack()
		})
	}
}
