# als-go

> 跨 shell 别名和定时任务管理工具

## 安装方法

**构建**

```shell
go build -o als
```

![image-20250611142056506](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20250611142039994.png)

**安装到系统路径**

```shell
# 方法1：手动复制到系统路径 (/usr/local/bin  存放用户自己安装的可执行文件)
sudo cp als /usr/local/bin/

# 方法2：添加到 PATH
export PATH="$(pwd):$PATH"
echo 'export PATH="'$(pwd)':$PATH"' >> ~/.bashrc
```

![image-20250611142039994](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20250611142056506.png)

## 使用方法

### 配置文件

vim ~/.config/alias.toml 

```toml
[aliases]
# 通用别名（所有shell通用）
ll = "ls -la"
gs = "git status"
ga = "git add"
gc = "git commit"
gp = "git push"

# 复杂命令别名
test = "echo Hello, World from Go!"

# cron 任务管理别名       
cronlist = "als cron list"        # 查看 cron 任务
cronmon = "als cron monitor"       # 监控 cron 任务
cronedit = "als cron edit"        # 编辑 cron 任务

# 为不同shell定义不同别名
[aliases.clear_screen]
bash = "clear"
zsh = "clear"
fish = "clear"
powershell = "Clear-Host"

[aliases.list_files]
bash = "ls -la"
zsh = "ls -la --color=auto"
fish = "ls -la"
powershell = "Get-ChildItem"

[aliases.remove_recursive]
bash = "rm -rf"
zsh = "rm -rf"
fish = "rm -rf"
powershell = "Remove-Item -Recurse -Force"

[aliases.cron_status]
bash = "systemctl status cron"
zsh = "systemctl status cron"
fish = "systemctl status cron"
powershell = "Get-Service -Name 'Task Scheduler'"

[aliases.alias_name]
bash = "bash_specific_command"
zsh = "zsh_specific_command"  
fish = "fish_specific_command"
powershell = "PowerShell_Specific_Command"
```

### 生成shell初始化脚本

根据您使用的shell，运行相应的命令：

```shell
# 为 bash 生成初始化脚本
als init bash

# 为 zsh 生成初始化脚本  
./als init zsh

# 为 fish 生成初始化脚本
./als init fish

# 为 PowerShell 生成初始化脚本
./als init powershell
```

![image-20250611142522969](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20250611225304905.png)

### Shell 集成

**Bash**

```bash
# 将初始化代码添加到 .bashrc
echo 'eval "$(als init bash)"' >> ~/.bashrc
# 重新加载配置 # 或直接重启终端
source ~/.bashrc
```

**Zsh**

```bash
# 将初始化代码添加到 .zshrc
echo 'eval "$(als init zsh)"' >> ~/.zshrc
# 重新加载配置 # 或直接重启终端
source ~/.zshrc
```

**Fish**

```shell
# 创建fish配置目录
mkdir -p ~/.config/fish
# 将初始化代码添加到 config.fish
echo 'als init fish | source' >> ~/.config/fish/config.fish
# 重新加载配置或重启fish 
```

**PowerShell**

将以下内容添加到 PowerShell 配置文件（`$PROFILE`）：

```powershell
Invoke-Expression (& als init powershell | Out-String)
# 直接重启终端
```

## 命令参考

### 基本命令

```shell
# 查看帮助
als --help

# 为 bash 生成初始化脚本
als init bash

# 为 zsh 生成初始化脚本  
als init zsh

# 为 fish 生成初始化脚本
als init fish

# 为 PowerShell 生成初始化脚本
als init powershell
```

![image-20250611225700222](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20250611142522969.png)

### Cron 任务管理

> 支持跨平台使用

#### 定时任务

vim /opt/software/test.sh

```bash
#!/bin/bash
echo "hello world" >>/opt/software/a.txt
```

添加新的 cron 任务

```shell
als cron add "* * * * *" "/opt/software/test.sh"    #每分钟执行
als cron add "0 2 * * *" "/opt/script/clean_xiaoduo_logs.sh"  # 每天凌晨 2 点执行
```

![image-20250611225304905](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20250611225700222.png)

查看当前用户的 cron 任务

```shell
als cron list
```

![image-20250611225321477](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20250611225321477.png)

编辑 cron 任务

```shell
als cron edit
```

实时监控 cron 任务执行情况

```shell
als cron monitor
```

删除指定行号的 cron 任务

```shell
als cron delete 3    # 删除第3行的任务
```

查看 cron 服务状态（使用别名）

```shell
cron_status
```

![image-20250611225432822](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20250611225432822.png)

####  例子

```bash
/opt # cd script/
opt/script # ls
clean_xiaoduo_logs.sh  connect_dev.sh  connect_mini.sh
opt/script # vim clean_xiaoduo_logs.sh
opt/script # chmod +x clean_xiaoduo_logs.sh 
opt/script # pwd                                                                
/opt/script
opt/script # als cron add "0 2 * * *" "/opt/script/clean_xiaoduo_logs.sh"                           
成功添加 cron 任务: 0 2 * * * /opt/script/clean_xiaoduo_logs.sh
opt/script # als cron list
当前 cron 任务:
------------------------------------------------------------
1: 0 2 * * * /opt/script/clean_xiaoduo_logs.sh
总共 1 个活跃任务
opt/script #     
```

![image-20250611232715033](https://github.com/Auroraol/Drawing-bed/raw/main/img/image-20250611232715033.png)

#### Cron 表达式格式

使用: [Crontab.guru - The cron schedule expression generator](https://crontab.guru/#01_)

```
分 时 日 月 周
*  *  *  *  *
│  │  │  │  │
│  │  │  │  └─── 周几 (0-7, 0和7都代表周日)
│  │  │  └────── 月份 (1-12)
│  │  └───────── 日期 (1-31)
│  └──────────── 小时 (0-23)
└─────────────── 分钟 (0-59)
```

常用 Cron 表达式示例

```bash
"0 0 * * *"     # 每天午夜
"0 2 * * *"     # 每天凌晨2点
"*/15 * * * *"  # 每15分钟
"0 9-17 * * 1-5" # 工作日每小时（9-17点）
"0 0 1 * *"     # 每月1号午夜
"0 0 * * 0"     # 每周日午夜
```

#### Windows 支持

在 Windows 系统上，该工具会自动使用 Windows 任务计划程序：

- `als cron list` - 显示计划任务
- `als cron add` - 创建新的计划任务（简化版本）
- `als cron edit` - 打开任务计划程序界面
- `als cron monitor` - 打开事件查看器查看任务日志

## 开发说明

### 技术架构

+ go
+ cobra(命令行程序库，可以用来编写命令行程序)

### 依赖说明
- `github.com/spf13/cobra`: 命令行接口框架
- `github.com/BurntSushi/toml`: TOML 配置文件解析

### 构建和测试

```bash
# 查看所有可用命令
make help

# 构建
make build

# 运行测试
make test

# 运行测试（带覆盖率）
make test-coverage

# 格式化代码
make fmt

# 代码检查
make vet

# 清理构建文件
make clean

# 构建多平台版本
make build-all
```

## 许可证

与原项目保持一致的许可证。

## 贡献

欢迎提交 Issues 和 Pull Requests！

### 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request
