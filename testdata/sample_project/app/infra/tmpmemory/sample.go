package tmpmemory

import (
	"github.com/kazdevl/sample_project/app/domain/entity"
	"github.com/kazdevl/sample_project/app/domain/ifrepository"
)

type SampleEntityRepository struct {
	data map[int]*entity.SampleEntity
}

var _ ifrepository.SampleEntityRepository = &SampleEntityRepository{}

func (r *SampleEntityRepository) Create(name string) {
	id := len(r.data) + 1
	r.data[id] = entity.NewSampleEntity(id, name)
}

func (r *SampleEntityRepository) FindById(id int) *entity.SampleEntity {
	if _, ok := r.data[id]; !ok {
		return nil
	}
	return r.data[id]
}
