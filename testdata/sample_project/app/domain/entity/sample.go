package entity

import (
	"github.com/kazdevl/sample_project/app/domain/model"
	gutil "github.com/kazdevl/sample_project/app/util"
)

type SampleEntity struct {
	m *model.SampleModel
}

func NewSampleEntity(id int, name string) *SampleEntity {
	return &SampleEntity{
		m: &model.SampleModel{
			ID:   id,
			Name: name + gutil.GetNowStringInJst(),
		},
	}
}
