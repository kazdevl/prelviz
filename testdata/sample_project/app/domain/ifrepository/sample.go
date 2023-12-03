package ifrepository

import "github.com/kazdevl/sample_project/app/domain/entity"

type SampleEntityRepository interface {
	Create(name string)
	FindById(id int) *entity.SampleEntity
}
