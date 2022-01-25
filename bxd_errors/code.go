package bxd_errors

// 特殊定义
const (
	BxdSuccessCode            = 0     // 成功
	BxdDefaultErrorCode = 50000 // 默认错误码
	BxdUnknownErrorCode = 50000 // 未知错误码
)

// 10000~10099
// 参数校验
const (
	BxdParamsDefaultErrorCode = 10000 // 参数校验异常
	BxdParamsNotEnoughCode    = 10001 // 参数不足
	BxdInvalidParamsCode      = 10002 // 非法参数

)

// 50200 ~ 50299
// 系统异常-连接异常
// 前缀：Bxd
const (
	BxdApiErrorCode   = 50201 + iota // API接口连接异常
	BxdMySQLErrorCode                // 数据库MySQL连接异常
	BxdRedisErrorCode                // Redis连接异常
	BxdKafkaErrorCode                // Kafka连接异常
	BxdRmqErrorCode                  // RabbitMQ连接异常
	BxdPanErrorCode                  // Pan连接异常
	BxdRpcErrorCode                  // RPC连接异常
)

// 50400 ~ 50499
// 系统异常-连接异常
// 前缀：Bxd
const (
	BxdApiTimeOutErrorCode   = 50401 + iota // API接口连接超时
	BxdMySQLTimeOutErrorCode                // 数据库MySQL连接超时
	BxdRedisTimeOutErrorCode                // Redis连接超时
	BxdKafkaTimeOutErrorCode                // Kafka连接超时
	BxdRmqTimeOutErrorCode                  // RabbitMQ连接超时
	BxdPanTimeOutErrorCode                  // Pan连接超时
	BxdRpcTimeOutErrorCode                  // RPC连接超时

)

// 50500 ~ 50599
// 系统异常-JSON错误
// 前缀：Bxd
const (
	BxdJsonMarshalErrorCode   = 50501 + iota // JSON Marshal错误
	BxdJsonUnmarshalErrorCode                // JSON Unmarshal错误
)
