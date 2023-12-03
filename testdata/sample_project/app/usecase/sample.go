package usecase

import (
	"github.com/kazdevl/sample_project/app/domain/entity"
	"github.com/kazdevl/sample_project/app/domain/model"
	"github.com/kazdevl/sample_project/app/domain/service"
)

type SampleUsecase struct {
	service *service.SampleService
}

func NewSampleUsecase(service *service.SampleService) *SampleUsecase {
	return &SampleUsecase{
		service: service,
	}
}

func (u *SampleUsecase) CreateSample(name string) {
	u.service.Create(name)
}

func (u *SampleUsecase) FindSampleById(id int) *entity.SampleEntity {
	return u.service.FindById(id)
}

func (u *SampleUsecase) NgSetting(name string) *model.SampleModel {
	return nil
}
