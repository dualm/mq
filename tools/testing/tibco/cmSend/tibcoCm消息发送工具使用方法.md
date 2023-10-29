1. 配置conf.toml文件并保存。具体选项如下：
   1. [Option]节点
      1. LedgerName: Ledger文件名
      2. RequestOld: 按照tibco需求填写
      3. SyncLedger: 按照tibco需求填写
      4. RelayAgent: 填写RelayAgent地址
      5. LimitTime: 填写消息等待时候，默认为0.0，即一直等待
      6. FieldName: 为空即可，后续Reports设置中填写
      7. Service: 按实际通信配置填写
      8. Network: 按实际通信配置填写
      9. Daemon: 按实际通信配置填写
      10. TargetSubjectName: 填写消息监听端的subject
      11. PooledMessage: 是否启用池化消息，建议为true
      12. CmName: 本地CmName，需要全网唯一
   2. [[Reports]]节点, 表示一条消息中的多个字段, 此时消息仅发送，不支持接收reply. 同一个[[Reports]]节点下的[[Reports.Fields]]为同一条消息的多个字段. 不同的[[Reports]]为不同的消息. 
      1. FieldName: 字段名
      2. Type: 字段类型，目前共有i16、i32、string、opaque可选
      3. Message: 字段内容
2. 首先进入命令行，cmd、powershell、terminal均可. 输入"cmSend + 配置文件名"(如cmSend conf.toml)即可运行工具，如果发送错误程序会输出错误信息. 如果发送正确则工具会在3秒后退出.