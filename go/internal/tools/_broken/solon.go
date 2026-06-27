package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func HandleGetRepoInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	repoURL, _ :=getString(args, "repo_url")
	if repoURL == "" {
		return err("repo_url parameter is required")
}

	// Validate GitHub/Gitee URL
	parsedURL, urlErr := url.Parse(repoURL)
	if urlErr != nil {
		return err("invalid repository URL")
}

	if !strings.Contains(parsedURL.Host, "github.com") && !strings.Contains(parsedURL.Host, "gitee.com") {
		return err("only GitHub and Gitee repositories are supported")
}

	// Extract repo path
	repoPath := strings.TrimPrefix(parsedURL.Path, "/")
	repoPath = strings.TrimSuffix(repoPath, ".git")
	parts := strings.Split(repoPath, "/")
	if len(parts) < 2 {
		return err("invalid repository path")
}

	owner := parts[0]
	repo := parts[1]

	// Get stars count (simplified - in real implementation would need API call)
	stars := 0
	if strings.Contains(parsedURL.Host, "github.com") {
		stars = 1000 // Mock value
	} else {
		stars = 500 // Mock value
	}

	result := map[string]interface{}{
		"owner":      owner,
		"repo":       repo,
		"stars":      stars,
		"platform":   strings.Split(parsedURL.Host, ".")[0],
		"created_at": time.Now().AddDate(-2, 0, 0).Format(time.RFC3339), // Mock value
	}

	return ok(result)
}

func HandleListTemplates(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	platform, _ :=getString(args, "platform")
	if platform == "" {
		platform = "github"
	}

	templates := []map[string]string{
		{
			"name":        "Issue Template (Chinese)",
			"path":        ".gitee/ISSUE_TEMPLATE.zh-CN.md",
			"description": "Chinese issue template for Gitee",
			"platform":    "gitee",
		},
		{
			"name":        "Pull Request Template",
			"path":        ".gitee/PULL_REQUEST_TEMPLATE.md",
			"description": "Pull request template for Gitee",
			"platform":    "gitee",
		},
		{
			"name":        "Issue Template (Chinese)",
			"path":        ".github/ISSUE_TEMPLATE.zh-CN.md",
			"description": "Chinese issue template for GitHub",
			"platform":    "github",
		},
		{
			"name":        "Bug Report Template",
			"path":        ".github/ISSUE_TEMPLATE/bug_report.md",
			"description": "Bug report template for GitHub",
			"platform":    "github",
		},
		{
			"name":        "Feature Request Template",
			"path":        ".github/ISSUE_TEMPLATE/feature_request.md",
			"description": "Feature request template for GitHub",
			"platform":    "github",
		},
		{
			"name":        "Question Template",
			"path":        ".github/ISSUE_TEMPLATE/problem_support.md",
			"description": "Question/support template for GitHub",
			"platform":    "github",
		},
		{
			"name":        "Pull Request Template",
			"path":        ".github/PULL_REQUEST_TEMPLATE.md",
			"description": "Pull request template for GitHub",
			"platform":    "github",
		},
		{
			"name":        "Contributing Guide",
			"path":        "CONTRIBUTING.md",
			"description": "Contribution guidelines",
			"platform":    "both",
		},
	}

	var filtered []map[string]string
	for _, t := range templates {
		if platform == "both" || t["platform"] == platform || t["platform"] == "both" {
			filtered = append(filtered, t)

	}

	return ok(filtered)
}

}

func HandleGetTemplateContent(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	templatePath, _ :=getString(args, "template_path")
	if templatePath == "" {
		return err("template_path parameter is required")
}

	// Mock template contents based on the Solon project
	var content string
	switch templatePath {
	case ".gitee/ISSUE_TEMPLATE.zh-CN.md":
		content = `### 问题描述

### 我当前使用 Solon 版本是?`
	case ".gitee/PULL_REQUEST_TEMPLATE.md":
		content = `### 这个PR有什么用 / 我们为什么需要它？

### 总结您的更改

#### 请注明您已完成以下工作：
- [ ] 确保测试通过，并在需要时添加测试覆盖率。
- [ ] 确保提交消息遵循 [常规提交规范](https://www.conventionalcommits.org/) 的规则。
- [ ] 考虑文档的影响，如果需要，打开一个新的文档问题或文档更改的PR。`
	case ".github/ISSUE_TEMPLATE.zh-CN.md":
		content = `### 问题描述

### 我当前使用 Solon 版本是?`
	case ".github/ISSUE_TEMPLATE/bug_report.md":
		content = `---
name: BUG 提交
about: 提交问题缺陷帮助我们更好的改进
title: '[BUG]'
labels: 'bug'
assignees: ''

---

#### 关联版本
*您当前正在使用我们框架的哪个版本？*

### 问题描述
*简要描述您碰到的问题。*

### 如何复现
*请详细告诉我们如何复现您遇到的问题，并附上可复现的代码示例*
1.
2.
3.
\`\`\`java
//可在此输入示例代码
\`\`\`

### 预期结果
*请告诉我们您预期会发生什么。*

### 实际结果
*请告诉我们实际发生了什么。*

### 截图或视频
*如果可以的话，上传任何关于 Bug 的截图。*`
	case ".github/ISSUE_TEMPLATE/feature_request.md":
		content = `---
name: 需求建议
about: 提出针对本项目的想法和建议
title: '[FEATURE]'
labels: 'enhancement'
assignees: ''

---

#### 关联版本
*您当前正在使用我们框架的哪个版本？*

### 请描述您的需求或者改进建议
*对您想要需求或建议的清晰简洁的描述。*

### 请描述你建议的实现方案
*对您想要需求或建议的实现方案的详细描述。*

### 描述您考虑过的替代方案
*对您考虑过的任何替代解决方案或功能的描述。*

#### 附加信息
*如果你还有其他需要提供的信息，可以在这里填写（可以提供截图、视频等）。*`
	case ".github/ISSUE_TEMPLATE/problem_support.md":
		content = `---
name: 问题支持
about: 提出针对本项目使用及其他方面的问题
title: '[QUESTION]'
labels: 'question'
assignees: ''

---

#### 关联版本
*您当前正在使用我们框架的哪个版本？*

### 请描述您的问题
*询问有关本项目的使用和其他方面的相关问题。*`
	case ".github/PULL_REQUEST_TEMPLATE.md":
		content = `### 这个PR有什么用 / 我们为什么需要它？

### 总结您的更改

#### 请注明您已完成以下工作：
- [ ] 确保测试通过，并在需要时添加测试覆盖率。
- [ ] 确保提交消息遵循 [常规提交规范](https://www.conventionalcommits.org/) 的规则。
- [ ] 考虑文档的影响，如果需要，打开一个新的文档问题或文档更改的PR。`
	case "CONTRIBUTING.md":
		content = `如果您对开源感兴趣且愿意学习和贡献，欢迎您共建 Solon 生态。

### 1、版权说明
本仓库的源码版权归 noear 开源组织所有。

### 2、贡献分类
代码贡献：
* 修复问题或优化现有的代码
* 新增功能插件
* 添加 Solon AI、Solon Cloud 等适配插件
* 为现有的模块丰富单元测试用例；为官网丰富配套示例。等...

合作贡献：
* 有开源框架的同道，在自己仓库里添加 solon 框架的便利适配（需要帮忙随时联系交流）
* 基于 Solon 开发开源项目或框架。等...

其它贡献：
* 通过 Issue，提交需求、提交问题
* 发博客宣传、录视频界面、在交流群或社区推荐 Solon。等...

### 3、代码贡献说明
1. 提交 Issue ，并与管理员进行确认（避免重复工作）
2. Fork 仓库
3. 在 main 分支上编写代码，并添加对应的单元测试
4. 统一使用 solon-test 做单测（为了批量跑单测）
5. pr 时，选择 main 分支进行合并（提交时需关联一个 Issue）
6. 如果是分布式中间件的适配，优先适配成 solon cloud 规范
7. 注释多些点：）

### 4、代码分支保护规则说明
| 操作 | main |
|----------------------|---------|
| 可推送代码成员 | 禁止任何人 |
| 可合并 Pull Request 成员 | 仓库管理员 |

### 5、代码模块测试目上录结构规范说明
| 目录 | 说明 |
|--------------------|-----------------------|
| src/test/benchmark | 压测目录（可选） |
| src/test/demo | 简单示例目录（必须，只是看看的放这里） |
| src/test/features | 特性测试目录（必选，会进入全项目批量单测） |
| src/test/labs | 实验目录（可选，不能批量跑的单测） |

不要增加别的目录

### 6、代码提前描述的前缀规范
| 前缀 | 示例 | 说明 |
|----|---------------------------------------|------------------|
| 新增 | 新增 solon-xxx 模块 | 表示增加一个全新模块 |
| 添加 | 添加 solon-xxx Yyy 工具类 | 表示在一个模块里增加新的能力 |
| 优化 | 优化 solon-xxx Yyy 延尽订阅处理逻辑 | 表示优化现有代码（没有兼容风险） |
| 修复 | 修复 solon-xxx Yyy 无法读取元数据问题 | 表示修复现有问题（没有兼容风险） |
| 调整 | 调整 solon-xxx Yyy 默认值为 true（之前为 false） | 表示调整现有代码（会有兼容风险） |
| 移除 | 移除 solon-xxx Yyy 注解类（之前已弃用一年） | 表示移除多余的类 |
| 文档 | 文档 solon-xxx Yyy 的注释完善 | 表示文档相关的完善 |
| 测试 | 测试 solon-xxx 补充 Yyy 测试用例 | 表示测试相关的完善 |
| 其它 | 其它 solon-xxx 配置示例变化 | 其它相关内容 |`
	default:
		return err(fmt.Sprintf("template %s not found", templatePath))
}

	return ok(map[string]string{
}
		"path":    templatePath,
		"content": content,
	})

func HandleValidateCommitMessage(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	message, _ :=getString(args, "message")
	if message == "" {
		return err("commit message is required")
}

	// Check conventional commits format
	pattern := `^(新增|添加|优化|修复|调整|移除|文档|测试|其它|feat|fix|docs|style|refactor|perf|test|chore)\(?[a-zA-Z0-9\-_]+\)?: .+$`
	re := regexp.MustCompile(pattern)

	if !re.MatchString(message) {
		return err("commit message doesn't follow conventional commits format")
}

	return ok("commit message is valid")
}

func HandleListRepoDirectories(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	repoType, _ :=getString(args, "repo_type")
	if repoType == "" {
		repoType = "main"
	}

	var dirs []map[string]string
	switch repoType {
	case "main":
		dirs = []map[string]string{
			{"path": "src/test/benchmark", "description": "压测目录（可选）"},
			{"path": "src/test/demo", "description": "简单示例目录（必须）"},
			{"path": "src/test/features", "description": "特性测试目录（必选）"},
			{"path": "src/test/labs", "description": "实验目录（可选）"},
		}
	case "plugins":
		dirs = []map[string]string{
			{"path": "src/main/java", "description": "主代码目录"},
			{"path": "src/test/java", "description": "测试代码目录"},
		}
	default:
		return err("invalid repo_type. Must be 'main' or 'plugins'")
}

	return ok(dirs)
}

func HandleGetRepoStats(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	repoPath, _ :=getString(args, "repo_path")
	if repoPath == "" {
		return err("repo_path parameter is required")
}

	// Mock git stats - in real implementation would use git commands
	stats := map[string]interface{}{
		"total_commits":   1250,
		"contributors":    42,
		"stars":           3200,
		"forks":           850,
		"open_issues":     45,
		"open_prs":        12,
		"last_commit":     time.Now().AddDate(0, 0, -2).Format(time.RFC3339),
		"default_branch":  "main",
		"license":         "Apache-2.0",
		"primary_language": "Java",
	}

	return ok(stats)
}