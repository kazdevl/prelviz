package service

import (
	"fmt"

	"github.com/kazdevl/sample_project/app/domain/entity"
	"github.com/kazdevl/sample_project/app/domain/ifrepository"
	"github.com/kazdevl/sample_project/app/domain/model"
)

type SampleService struct {
	repo ifrepository.SampleEntityRepository
}

func NewSampleService(repo ifrepository.SampleEntityRepository) *SampleService {
	return &SampleService{
		repo: repo,
	}
}

func (s *SampleService) Create(name string) {
	s.repo.Create(fmt.Sprintf("test_%s", name))
}

func (s *SampleService) FindById(id int) *entity.SampleEntity {
	return s.repo.FindById(id)
}

func (s *SampleService) NgSetting(name string) *model.SampleModel {
	return nil
}
