package backendlearn



type ServerCtx struct {
	C  *config.Config
	DB *gorm.DB
	//ImageMgr image.ImageManager
	Dao      *dao.Dao
	KvStore  *xkv.Store
	RankKey  string
	NodeSrvs map[int64]*nftchainservice.Service
}




//********************** 函数式选项模式 ********************************//
type CtxConfig struct {  // 配置结构体，用于存储服务器上下文所需的各种依赖
	db *gorm.DB   // 是小写的（如 db, dao），表示它们是包内私有的，这提供了封装性
	//imageMgr image.ImageManager
	dao     *dao.Dao
	kvStore *xkv.Store
	evm     erc.Erc  <----------可选
}

type CtxOption func(conf *CtxConfig)

func NewServerCtx(options ...CtxOption) *ServerCtx {
	c := &CtxConfig{}
	for _, opt := range options {
		opt(c)
	}
	return &ServerCtx{
		DB: c.db,
		//ImageMgr: c.imageMgr,
		KvStore: c.KvStore,
		Dao:     c.dao,
	}
}

func WithKv(kv *xkv.Store) CtxOption {
	return func(conf *CtxConfig) {
		conf.KvStore = kv
	}
}

func WithDB(db *gorm.DB) CtxOption {
	return func(conf *CtxConfig) {
		conf.db = db
	}
}

func WithDao(dao *dao.Dao) CtxOption {
	return func(conf *CtxConfig) {
		conf.dao = dao
	}
}

serverCtx := NewServerCtx(
		WithDB(db),   <--------
		WithKv(store),
		//WithImageMgr(imageMgr),
		WithDao(dao),
	)

	//  解决参数过多，顺序， 可选，可向后兼容 