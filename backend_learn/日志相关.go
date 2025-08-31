



_, err = xzap.SetUp(*cfg.Log) // 初始化日志模块
if err != nil {
	xzap.WithContext(ctx).Error("Failed to set up logger", zap.Error(err))
	onSyncExit <- err
	return
}
/*
记录启动信息: 使用Info级别记录服务器启动日志
携带上下文: WithContext(ctx)将上下文信息（如traceID等）注入日志
输出配置信息: zap.Any("config", cfg)将整个配置对象以任意格式记录到日志中
*/
xzap.WithContext(ctx).Info("sync server start", zap.Any("config", cfg))
