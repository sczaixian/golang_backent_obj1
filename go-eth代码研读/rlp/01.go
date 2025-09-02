



decode.go                        解码器，把RLP数据解码为go的数据结构
doc.go                                文档代码
encode.go                        编码器，把GO的数据结构序列化为字节数组
raw.go                                未解码的RLP数据
typecache.go                        类型缓存， 类型缓存记录了类型->(编码器|解码器)的内容。