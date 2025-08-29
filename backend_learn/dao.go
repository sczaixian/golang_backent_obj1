package backendlearn






type Dao struct {
	ctx context.Context

	DB      *gorm.DB
	KvStore *xkv.Store
}