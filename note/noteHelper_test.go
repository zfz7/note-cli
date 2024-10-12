package main

import (
	"github.com/stretchr/testify/mock"
	"note/mocks"
	"testing"
	"time"
)

func Test_noteHelper_OpenNote(t *testing.T) {
	Now = func() time.Time {
		return time.Date(2022, time.October, 1, 10, 15, 10, 0, time.UTC)
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
	dayConfig := NoteConfig{
		Editor:    defaultConfig.Editor,
		Location:  defaultConfig.Location,
		Template:  defaultConfig.Template,
		Extension: defaultConfig.Extension,
		Interval:  "day",
	}
	monthConfig := NoteConfig{
		Editor:    defaultConfig.Editor,
		Location:  defaultConfig.Location,
		Template:  defaultConfig.Template,
		Extension: defaultConfig.Extension,
		Interval:  "MONTH",
	}
	type fields struct {
		fileHelper FileHelper
	}
	type args struct {
		relativeInterval int
		config           NoteConfig
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		AssertCallBack func()
	}{
		{
			name:    "Note exists, relativeInterval: 0, interval: week",
			fields:  fields{fileHelper: noteExistsHelper},
			args:    args{relativeInterval: 0, config: defaultConfig},
			wantErr: false,
			AssertCallBack: func() {
				noteExistsHelper.AssertCalled(t, "EditorOpenFile", defaultConfig.Editor, "~/notes/2022-09-26."+defaultConfig.Extension)
			},
		},
		{
			name:    "Note exists, relativeInterval: 1, interval: week",
			fields:  fields{fileHelper: noteExistsHelper},
			args:    args{relativeInterval: 1, config: defaultConfig},
			wantErr: false,
			AssertCallBack: func() {
				noteExistsHelper.AssertCalled(t, "EditorOpenFile", defaultConfig.Editor, "~/notes/2022-10-03."+defaultConfig.Extension)
			},
		},
		{
			name:    "Note exists, relativeInterval: -1, interval: week",
			fields:  fields{fileHelper: noteExistsHelper},
			args:    args{relativeInterval: -1, config: defaultConfig},
			wantErr: false,
			AssertCallBack: func() {
				noteExistsHelper.AssertCalled(t, "EditorOpenFile", defaultConfig.Editor, "~/notes/2022-09-19."+defaultConfig.Extension)
			},
		},
		{
			name:    "Note exists, relativeInterval: 0, interval: day",
			fields:  fields{fileHelper: noteExistsHelper},
			args:    args{relativeInterval: 0, config: dayConfig},
			wantErr: false,
			AssertCallBack: func() {
				noteExistsHelper.AssertCalled(t, "EditorOpenFile", dayConfig.Editor, "~/notes/2022-10-01."+dayConfig.Extension)
			},
		},
		{
			name:    "Note exists, relativeInterval: 1, interval: day",
			fields:  fields{fileHelper: noteExistsHelper},
			args:    args{relativeInterval: 1, config: dayConfig},
			wantErr: false,
			AssertCallBack: func() {
				noteExistsHelper.AssertCalled(t, "EditorOpenFile", dayConfig.Editor, "~/notes/2022-10-02."+dayConfig.Extension)
			},
		},
		{
			name:    "Note exists, relativeInterval: -1, interval: day",
			fields:  fields{fileHelper: noteExistsHelper},
			args:    args{relativeInterval: -1, config: dayConfig},
			wantErr: false,
			AssertCallBack: func() {
				noteExistsHelper.AssertCalled(t, "EditorOpenFile", dayConfig.Editor, "~/notes/2022-09-30."+dayConfig.Extension)
			},
		},
		{
			name:    "Note exists, relativeInterval: 0, interval: month",
			fields:  fields{fileHelper: noteExistsHelper},
			args:    args{relativeInterval: 0, config: monthConfig},
			wantErr: false,
			AssertCallBack: func() {
				noteExistsHelper.AssertCalled(t, "EditorOpenFile", monthConfig.Editor, "~/notes/2022-10-01."+monthConfig.Extension)
			},
		},
		{
			name:    "Note exists, relativeInterval: 1, interval: month",
			fields:  fields{fileHelper: noteExistsHelper},
			args:    args{relativeInterval: 1, config: monthConfig},
			wantErr: false,
			AssertCallBack: func() {
				noteExistsHelper.AssertCalled(t, "EditorOpenFile", monthConfig.Editor, "~/notes/2022-11-01."+monthConfig.Extension)
			},
		},
		{
			name:    "Note exists, relativeInterval: -1, interval: month",
			fields:  fields{fileHelper: noteExistsHelper},
			args:    args{relativeInterval: -1, config: monthConfig},
			wantErr: false,
			AssertCallBack: func() {
				noteExistsHelper.AssertCalled(t, "EditorOpenFile", dayConfig.Editor, "~/notes/2022-09-01."+dayConfig.Extension)
			},
		},
		{
			name:    "Note does not exists, relativeInterval: 0, interval: week",
			fields:  fields{fileHelper: noteMissingHelper},
			args:    args{relativeInterval: 0, config: defaultConfig},
			wantErr: false,
			AssertCallBack: func() {
				noteMissingHelper.AssertCalled(t, "WriteFile", defaultConfig.Location+"/2022-09-26."+defaultConfig.Extension, mock.Anything)
				noteMissingHelper.AssertCalled(t, "EditorOpenFile", defaultConfig.Editor, "~/notes/2022-09-26."+defaultConfig.Extension)
			},
		},
		{
			name:    "Note does not exists, relativeInterval: 1, interval: week",
			fields:  fields{fileHelper: noteMissingHelper},
			args:    args{relativeInterval: 1, config: defaultConfig},
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
			if err := noteHelper.OpenNote(tt.args.relativeInterval, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("OpenNote() error = %v, wantErr %v", err, tt.wantErr)
			}
			tt.AssertCallBack()
		})
	}
}
