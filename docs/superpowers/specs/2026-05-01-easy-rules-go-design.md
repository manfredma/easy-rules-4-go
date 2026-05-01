# easy-rules-4-go 设计文档

## 概述

将 Java 版 [easy-rules 4.x](https://github.com/j-easy/easy-rules) 用 Go 重新实现。保留原项目的核心概念和 API 语义，同时适配 Go 的惯用风格。

## 目标

- 完整实现 easy-rules 核心功能
- 支持三种规则定义方式：函数式 Builder API、struct 嵌入、表达式引擎
- 支持 YAML/JSON 规则文件加载
- 实现 DefaultRulesEngine 和 InferenceRulesEngine
- 实现三种组合规则：UnitRuleGroup、ActivationRuleGroup、ConditionalRuleGroup
- 全面的测试覆盖（单元测试 + 集成测试）

## 项目结构

```
easy-rules-4-go/
├── go.mod                          # module: github.com/manfredma/easy-rules-4-go
├── go.sum
├── pkg/
│   ├── api/                        # 核心接口和类型定义
│   │   ├── rule.go                 # Rule 接口
│   │   ├── facts.go                # Facts + Fact 类型
│   │   ├── rules.go                # Rules 有序集合
│   │   ├── condition.go            # Condition 函数类型
│   │   ├── action.go               # Action 函数类型
│   │   ├── engine.go               # RulesEngine 接口 + EngineParameters
│   │   └── listener.go             # RuleListener + RulesEngineListener 接口
│   ├── core/                       # 核心实现
│   │   ├── basic_rule.go           # BasicRule（可嵌入扩展的基础规则）
│   │   ├── default_rule.go         # DefaultRule（由 RuleBuilder 构建）
│   │   ├── rule_builder.go         # RuleBuilder 流式 API
│   │   ├── default_engine.go       # DefaultRulesEngine
│   │   ├── default_engine_test.go
│   │   ├── inference_engine.go     # InferenceRulesEngine
│   │   └── inference_engine_test.go
│   ├── composite/                  # 组合规则
│   │   ├── composite_rule.go       # CompositeRule 基类
│   │   ├── unit_rule_group.go
│   │   ├── unit_rule_group_test.go
│   │   ├── activation_rule_group.go
│   │   ├── activation_rule_group_test.go
│   │   ├── conditional_rule_group.go
│   │   └── conditional_rule_group_test.go
│   ├── expr/                       # 表达式引擎集成（expr-lang/expr）
│   │   ├── expr_condition.go       # ExprCondition：字符串表达式 → Condition
│   │   ├── expr_action.go          # ExprAction：字符串表达式 → Action
│   │   ├── expr_rule.go            # ExprRule + ExprRuleFactory
│   │   └── expr_rule_test.go
│   └── reader/                     # 规则文件读取
│       ├── rule_definition.go      # RuleDefinition 结构体（YAML/JSON 映射）
│       ├── yaml_reader.go          # YAML 规则文件加载
│       ├── json_reader.go          # JSON 规则文件加载
│       └── reader_test.go
└── examples/
    ├── hello_world/
    │   └── main.go
    ├── weather/
    │   └── main.go
    └── fizzbuzz/
        └── main.go
```

## 核心 API 设计

### Rule 接口

```go
type Rule interface {
    GetName() string
    GetDescription() string
    GetPriority() int
    Evaluate(facts *Facts) bool
    Execute(facts *Facts) error
}
```

### Facts

```go
type Fact struct {
    Name  string
    Value any
}

type Facts struct {
    facts map[string]any
}
// Put(name, value) / Get(name) / Remove(name) / AsMap() / Clear()
```

### Rules（有序集合）

```go
type Rules struct {
    rules []Rule   // 按 priority ASC、同 priority 按 name 字典序排序
}
// Register(...Rule) / Unregister(name) / IsEmpty() / Size() / Iter()
```

### RulesEngine 接口

```go
type RulesEngine interface {
    GetParameters() *EngineParameters
    GetRuleListeners() []RuleListener
    GetEngineListeners() []RulesEngineListener
    Fire(rules *Rules, facts *Facts)
    Check(rules *Rules, facts *Facts) map[Rule]bool
}
```

### EngineParameters

```go
type EngineParameters struct {
    SkipOnFirstAppliedRule     bool
    SkipOnFirstNonTriggeredRule bool
    SkipOnFirstFailedRule      bool
    PriorityThreshold          int  // 默认 math.MaxInt
}
```

### Condition / Action 函数类型

```go
type Condition func(facts *Facts) bool
type Action    func(facts *Facts) error

var ConditionTrue  Condition = func(_ *Facts) bool { return true }
var ConditionFalse Condition = func(_ *Facts) bool { return false }
```

## 规则定义方式

### 方式一：函数式 Builder API

```go
rule := core.NewRuleBuilder().
    Name("weather rule").
    Description("if it rains then take an umbrella").
    Priority(1).
    When(func(facts *api.Facts) bool {
        rain, _ := facts.Get("rain").(bool)
        return rain
    }).
    Then(func(facts *api.Facts) error {
        fmt.Println("It rains, take an umbrella!")
        return nil
    }).
    Build()
```

### 方式二：struct 嵌入 BasicRule

```go
type WeatherRule struct {
    core.BasicRule
}

func (r *WeatherRule) Evaluate(facts *api.Facts) bool {
    rain, _ := facts.Get("rain").(bool)
    return rain
}

func (r *WeatherRule) Execute(facts *api.Facts) error {
    fmt.Println("It rains, take an umbrella!")
    return nil
}
```

### 方式三：表达式引擎（expr-lang/expr）

```go
rule, err := expr.NewExprRuleBuilder().
    Name("weather rule").
    When(`rain == true`).
    Then(`println("It rains, take an umbrella!")`).
    Build()
```

### 方式四：YAML 文件加载

```yaml
# weather-rules.yml
- name: "weather rule"
  description: "if it rains then take an umbrella"
  priority: 1
  condition: "rain == true"
  actions:
    - 'println("It rains, take an umbrella!")'
```

```go
factory := reader.NewYamlRuleFactory()
rules, err := factory.CreateRulesFrom("weather-rules.yml")
```

## 引擎实现

### DefaultRulesEngine

按 priority 顺序遍历规则，逐条 Evaluate → Execute，受 EngineParameters 控制跳过行为：
- `SkipOnFirstAppliedRule`：首条触发后跳过剩余
- `SkipOnFirstFailedRule`：首条执行失败后跳过剩余
- `SkipOnFirstNonTriggeredRule`：首条条件为 false 后跳过剩余
- `PriorityThreshold`：超过阈值的规则跳过

### InferenceRulesEngine

反复循环：先选出本轮所有条件为 true 的候选规则，交给 DefaultRulesEngine 执行，直到没有候选规则为止。

## 组合规则

| 类型 | 语义 |
|------|------|
| `UnitRuleGroup` | 所有子规则条件均为 true 才执行全部（全或无） |
| `ActivationRuleGroup` | 按优先级找第一个条件为 true 的子规则执行（XOR） |
| `ConditionalRuleGroup` | 优先级最高的规则作为"条件门"，通过后执行其余满足条件的子规则 |

## Java → Go 适配说明

| Java | Go |
|------|-----|
| `interface` + default 方法 | `interface` + 嵌入 struct 提供默认实现 |
| `@FunctionalInterface` | `type Condition func(...)` |
| `throws Exception` | 返回 `error` |
| `TreeSet`（自然排序） | `[]Rule` 插入时保持有序 |
| `Fact<T>` 泛型 | `any` 类型，调用方做类型断言 |
| `Object...` 可变参 | `...Rule` 可变参 |
| Java 注解模型 | 不实现（无对应 Go 惯用法） |

## 测试策略

每个包都需要对应的 `_test.go` 文件，覆盖：

- **单元测试**：每个类型的基本行为、边界条件、错误路径
- **引擎测试**：EngineParameters 各参数组合的行为验证
- **组合规则测试**：三种组合规则的触发语义
- **表达式引擎测试**：条件/动作表达式的编译与执行
- **文件加载测试**：YAML/JSON 规则文件的完整解析流程
- **集成测试**：完整场景（hello world、天气规则、FizzBuzz）

目标覆盖率：≥ 80%。

## 依赖

```
github.com/expr-lang/expr   # 表达式引擎
gopkg.in/yaml.v3            # YAML 解析
```
