package backendlearn








Clauses 的底层原理
GORM 执行操作（如 Find）时，会创建 *gorm.Statement 对象，
所有链式调用（如 Where、Limit）本质是向该对象添加子句。
Clauses() 方法允许注入自定义子句生成器，最终拼接为完整 SQL 





处理数据冲突（UPSERT）
当插入数据时，若主键或唯一索引冲突，自动转为更新操作。
示例：批量更新用户分数，冲突时更新 score 字段
db.Clauses(clause.OnConflict{
    Columns:   []clause.Column{{Name: "id"}}, // 冲突字段
    DoUpdates: clause.AssignmentColumns([]string{"score"}), // 冲突时更新的字段
}).Create(&users)
适用场景：避免循环逐条更新，一次性完成批量插入或更新


		
实现行级锁定（并发控制）
高并发下保证数据一致性，如悲观锁	
示例：查询用户时加排他锁（FOR UPDATE）	
db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&users)
适用场景：订单支付时锁定库存，防止超卖


// 优化软删除查询性能     TODO:
GORM 默认软删除使用 deleted_at IS NULL，可能导致索引失效。通过重写 QueryClauses 优化
func (DeletedAt) QueryClauses(f *schema.Field) []clause.Interface {
    return []clause.Interface{BeautifulSoftDeleteQueryClause{}}
}
// 自定义生成 IFNULL(deleted_at,0)=0 条件




实际开发中的经典案例

1. 动态分页封装
type Pagination struct{ Page, PageSize int }
func (p *Pagination) ModifyStatement(stmt *gorm.Statement) {
    stmt.DB.Offset((p.Page-1)*p.PageSize).Limit(p.PageSize)
}
// 使用
db.Clauses(&Pagination{Page: 2, PageSize: 10}).Find(&users)

生成 LIMIT 10 OFFSET 10，避免重复代码



子查询嵌套
subQuery := db.Model(&Order{}).Select("AVG(amount)")
db.Where("amount > (?)", subQuery).Find(&orders)

生成 SELECT * FROM orders WHERE amount > (SELECT AVG(amount) FROM orders) 
















DoNothing: true这个选项。根据搜索结果，这相当于SQL中的ON CONFLICT DO NOTHING。
当插入的数据与表中现有数据的主键或唯一索引冲突时，数据库会静默跳过这条插入，既不会更新现有记录，也不会报错
s.db.WithContext(s.ctx).Table(multi.OrderTableName(s.chain)).
	Clauses(clause.OnConflict{DoNothing: true,}).
	Create(&newOrder).Error

















------------------  Scan  -------------------------------

当返回结果只有一个值时  scan 默认接收， 当有多个值时 按照名称对应接收


当查询的字段与模型结构不匹配时，使用Scan将结果映射到自定义结构
type OrderResult struct {
    OrderID   string
    Total     float64
    CreatedAt time.Time
}
var result OrderResult
db.Table("orders").Select("order_id, total, created_at").
    Where("id = ?", 1).
    Scan(&result)


只查询需要的字段，提高性能  <---------   虽然model 中可能有更多字段，但是查询出的之后需要的字段
var orderInfo struct {
    ID    uint
    Price float64
    // ....
}
db.Model(&Order{}).Where("id = ?", 1).Scan(&orderInfo)


           Scan vs Find 的区别

特性               Find                        Scan
映射目标         模型结构体                任意结构体/变量
字段匹配       严格按模型字段                 按名称匹配
性能            查询所有字段                只查询指定字段
使用场景          CRUD操作                 自定义查询、报表


// 2. 查询特定字段
var orderSummary struct {
    ID     uint
    Number string
    Status string
}
db.Model(&Order{}).Where("id = ?", 1).Scan(&orderSummary)

// 3. 扫描到map
var result map[string]interface{}
db.Model(&Order{}).First(&Order{}).Scan(&result)


字段顺序不重要: Scan 按名称匹配，不是按顺序
大小写不敏感:   SQL 的列名通常不区分大小写
缺失字段：      如果结构体没有对应的字段，该列数据会被忽略
多余字段：      如果查询结果有结构体没有的列，不会报错