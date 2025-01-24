package xerr

// 成功返回
const OK int64 = 200

/**(前3位代表业务,后三位代表具体功能)**/

// 全局错误码
const Fail int64 = 100000
const SERVER_COMMON_ERROR int64 = 100001
const REUQEST_PARAM_ERROR int64 = 100002
const TOKEN_EXPIRE_ERROR int64 = 100003
const TOKEN_GENERATE_ERROR int64 = 100004
const DB_ERROR int64 = 100005
const DB_UPDATE_AFFECTED_ZERO_ERROR int64 = 100006

//用户模块
