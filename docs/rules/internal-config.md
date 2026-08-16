# internal/config 层规则

- 修改添加配置字段
  - 只能在 internal/config/dev的目录下修改配置。
  - 配置字段为空，必须占位，可用类型缺省值占位。
  - _开头的字段名为可重复类型，其余字段名必须包内唯一。
  - 修改完成后，执行make config命令，生成相应结构体。
  - configs目录为当前服务执行的配置目录，可通过make initConfig填充，其余环境可自行调整命令。
- configs 目录下修改数值。当前运行服务会自动跟新。
