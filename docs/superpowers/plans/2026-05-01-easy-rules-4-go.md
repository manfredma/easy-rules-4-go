# easy-rules-4-go Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 用 Go 重新实现 Java easy-rules 4.x 规则引擎，支持函数式 Builder API、struct 嵌入、表达式引擎和 YAML/JSON 文件加载四种规则定义方式。

**Architecture:** 分为 pkg/api（接口定义）、pkg/core（核心实现）、pkg/composite（组合规则）、pkg/expr（表达式引擎集成）、pkg/reader（文件读取）五个子包，每层依赖方向向上（core 依赖 api，composite/expr/reader 依赖 api+core）。

**Tech Stack:** Go 1.21+，github.com/expr-lang/expr（表达式引擎），gopkg.in/yaml.v3（YAML 解析），标准库 encoding/json（JSON 解析）

---

## 文件清单

| 文件 | 职责 |
|------|------|
| `go.mod` | 模块声明和依赖 |
| `pkg/api/rule.go` | Rule 接口 |
| `pkg/api/facts.go` | Facts + Fact 类型 |
| `pkg/api/facts_test.go` | Facts 单元测试 |
| `pkg/api/rules.go` | Rules 有序集合 |
| `pkg/api/rules_test.go` | Rules 单元测试 |
| `pkg/api/condition.go` | Condition 函数类型 + TRUE/FALSE 常量 |
| `pkg/api/action.go` | Action 函数类型 |
| `pkg/api/engine.go` | RulesEngine 接口 + EngineParameters |
| `pkg/api/listener.go` | RuleListener + RulesEngineListener 接口 |
| `pkg/core/basic_rule.go` | BasicRule struct（可嵌入） |
| `pkg/core/basic_rule_test.go` | BasicRule 单元测试 |
| `pkg/core/default_rule.go` | DefaultRule（Builder 构建产物） |
| `pkg/core/rule_builder.go` | RuleBuilder 流式 API |
| `pkg/core/rule_builder_test.go` | RuleBuilder 单元测试 |
| `pkg/core/default_engine.go` | DefaultRulesEngine |
| `pkg/core/default_engine_test.go` | DefaultRulesEngine 单元测试（含参数组合） |
| `pkg/core/inference_engine.go` | InferenceRulesEngine |
| `pkg/core/inference_engine_test.go` | InferenceRulesEngine 单元测试 |
| `pkg/composite/composite_rule.go` | CompositeRule 基类 |
| `pkg/composite/unit_rule_group.go` | UnitRuleGroup |
| `pkg/composite/unit_rule_group_test.go` | UnitRuleGroup 单元测试 |
| `pkg/composite/activation_rule_group.go` | ActivationRuleGroup |
| `pkg/composite/activation_rule_group_test.go` | ActivationRuleGroup 单元测试 |
| `pkg/composite/conditional_rule_group.go` | ConditionalRuleGroup |
| `pkg/composite/conditional_rule_group_test.go` | ConditionalRuleGroup 单元测试 |
| `pkg/expr/expr_condition.go` | ExprCondition：字符串表达式编译为 Condition |
| `pkg/expr/expr_action.go` | ExprAction：字符串表达式编译为 Action（通过 map 输出副作用） |
| `pkg/expr/expr_rule.go` | ExprRule + ExprRuleBuilder |
| `pkg/expr/expr_rule_test.go` | ExprRule 单元测试 |
| `pkg/reader/rule_definition.go` | RuleDefinition 结构体 |
| `pkg/reader/yaml_reader.go` | YamlRuleFactory |
| `pkg/reader/json_reader.go` | JsonRuleFactory |
| `pkg/reader/reader_test.go` | YAML/JSON 加载集成测试 |
| `examples/hello_world/main.go` | Hello World 示例 |
| `examples/weather/main.go` | 天气规则示例（struct 嵌入方式） |
| `examples/fizzbuzz/main.go` | FizzBuzz 示例（表达式引擎方式） |

---

## Task 1: 初始化 Go 模块

**Files:**
- Create: `go.mod`

- [ ] **Step 1: 初始化模块**

```bash
cd /Users/maxingfang/GolandProjects/easy-rules-4-go
go mod init github.com/manfredma/easy-rules-4-go
```

Expected: 生成 `go.mod`，内容包含 `module github.com/manfredma/easy-rules-4-go`

- [ ] **Step 2: 添加依赖**

```bash
go get github.com/expr-lang/expr@latest
go get gopkg.in/yaml.v3@latest
```

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit --author="maxingfang <maxingfang@126.com>" -m "chore: init go module with dependencies"
```

---

## Task 2: pkg/api — 核心接口定义

**Files:**
- Create: `pkg/api/rule.go`
- Create: `pkg/api/condition.go`
- Create: `pkg/api/action.go`
- Create: `pkg/api/listener.go`
- Create: `pkg/api/engine.go`

- [ ] **Step 1: 创建 pkg/api/rule.go**

```go
package api

const (
    DefaultRuleName        = "rule"
    DefaultRuleDescription = "description"
    DefaultRulePriority    = int(^uint(0)>>1) - 1 // math.MaxInt - 1
)

type Rule interface {
    GetName() string
    GetDescription() string
    GetPriority() int
    Evaluate(facts *Facts) bool
    Execute(facts *Facts) error
}
```

- [ ] **Step 2: 创建 pkg/api/condition.go**

```go
package api

type Condition func(facts *Facts) bool

var ConditionTrue Condition = func(_ *Facts) bool { return true }
var ConditionFalse Condition = func(_ *Facts) bool { return false }
```

- [ ] **Step 3: 创建 pkg/api/action.go**

```go
package api

type Action func(facts *Facts) error
```

- [ ] **Step 4: 创建 pkg/api/listener.go**

```go
package api

type RuleListener interface {
    BeforeEvaluate(rule Rule, facts *Facts) bool
    AfterEvaluate(rule Rule, facts *Facts, evaluationResult bool)
    OnEvaluationError(rule Rule, facts *Facts, err error)
    BeforeExecute(rule Rule, facts *Facts)
    OnSuccess(rule Rule, facts *Facts)
    OnFailure(rule Rule, facts *Facts, err error)
}

type RulesEngineListener interface {
    BeforeEvaluate(rules *Rules, facts *Facts)
    AfterExecute(rules *Rules, facts *Facts)
}

// DefaultRuleListener provides no-op implementations of RuleListener.
type DefaultRuleListener struct{}

func (d *DefaultRuleListener) BeforeEvaluate(_ Rule, _ *Facts) bool          { return true }
func (d *DefaultRuleListener) AfterEvaluate(_ Rule, _ *Facts, _ bool)        {}
func (d *DefaultRuleListener) OnEvaluationError(_ Rule, _ *Facts, _ error)   {}
func (d *DefaultRuleListener) BeforeExecute(_ Rule, _ *Facts)                {}
func (d *DefaultRuleListener) OnSuccess(_ Rule, _ *Facts)                    {}
func (d *DefaultRuleListener) OnFailure(_ Rule, _ *Facts, _ error)           {}

// DefaultRulesEngineListener provides no-op implementations of RulesEngineListener.
type DefaultRulesEngineListener struct{}

func (d *DefaultRulesEngineListener) BeforeEvaluate(_ *Rules, _ *Facts) {}
func (d *DefaultRulesEngineListener) AfterExecute(_ *Rules, _ *Facts)   {}
```

- [ ] **Step 5: 创建 pkg/api/engine.go**

```go
package api

import "math"

type EngineParameters struct {
    SkipOnFirstAppliedRule      bool
    SkipOnFirstNonTriggeredRule bool
    SkipOnFirstFailedRule       bool
    PriorityThreshold           int
}

func NewEngineParameters() *EngineParameters {
    return &EngineParameters{
        PriorityThreshold: math.MaxInt,
    }
}

type RulesEngine interface {
    GetParameters() *EngineParameters
    GetRuleListeners() []RuleListener
    GetEngineListeners() []RulesEngineListener
    Fire(rules *Rules, facts *Facts)
    Check(rules *Rules, facts *Facts) map[Rule]bool
}
```

- [ ] **Step 6: 验证编译**

```bash
go build ./pkg/api/...
```

Expected: 无错误（此时 Facts/Rules 未定义会报错，下一个 task 补充）

- [ ] **Step 7: Commit**

```bash
git add pkg/api/rule.go pkg/api/condition.go pkg/api/action.go pkg/api/listener.go pkg/api/engine.go
git commit --author="maxingfang <maxingfang@126.com>" -m "feat(api): add core interfaces - Rule, Condition, Action, Listener, RulesEngine"
```

---

## Task 3: pkg/api — Facts 和 Rules

**Files:**
- Create: `pkg/api/facts.go`
- Create: `pkg/api/facts_test.go`
- Create: `pkg/api/rules.go`
- Create: `pkg/api/rules_test.go`

- [ ] **Step 1: 写 Facts 测试（先写测试）**

```go
// pkg/api/facts_test.go
package api

import "testing"

func TestFacts_PutAndGet(t *testing.T) {
    facts := NewFacts()
    facts.Put("rain", true)
    val := facts.Get("rain")
    if val != true {
        t.Errorf("expected true, got %v", val)
    }
}

func TestFacts_PutOverwrite(t *testing.T) {
    facts := NewFacts()
    facts.Put("rain", true)
    facts.Put("rain", false)
    if facts.Get("rain") != false {
        t.Error("expected overwrite to false")
    }
}

func TestFacts_Remove(t *testing.T) {
    facts := NewFacts()
    facts.Put("rain", true)
    facts.Remove("rain")
    if facts.Get("rain") != nil {
        t.Error("expected nil after remove")
    }
}

func TestFacts_AsMap(t *testing.T) {
    facts := NewFacts()
    facts.Put("a", 1)
    facts.Put("b", 2)
    m := facts.AsMap()
    if m["a"] != 1 || m["b"] != 2 {
        t.Error("AsMap mismatch")
    }
}

func TestFacts_Clear(t *testing.T) {
    facts := NewFacts()
    facts.Put("a", 1)
    facts.Clear()
    if facts.Get("a") != nil {
        t.Error("expected nil after clear")
    }
}
```

- [ ] **Step 2: 运行测试，确认失败**

```bash
go test ./pkg/api/... -run TestFacts -v
```

Expected: FAIL — `NewFacts undefined`

- [ ] **Step 3: 实现 pkg/api/facts.go**

```go
package api

import "fmt"

type Fact struct {
    Name  string
    Value any
}

func (f *Fact) String() string {
    return fmt.Sprintf("Fact{name=%s, value=%v}", f.Name, f.Value)
}

type Facts struct {
    facts map[string]any
}

func NewFacts() *Facts {
    return &Facts{facts: make(map[string]any)}
}

func (f *Facts) Put(name string, value any) {
    f.facts[name] = value
}

func (f *Facts) Get(name string) any {
    return f.facts[name]
}

func (f *Facts) Remove(name string) {
    delete(f.facts, name)
}

func (f *Facts) AsMap() map[string]any {
    copy := make(map[string]any, len(f.facts))
    for k, v := range f.facts {
        copy[k] = v
    }
    return copy
}

func (f *Facts) Clear() {
    f.facts = make(map[string]any)
}

func (f *Facts) String() string {
    return fmt.Sprintf("%v", f.facts)
}
```

- [ ] **Step 4: 运行测试，确认通过**

```bash
go test ./pkg/api/... -run TestFacts -v
```

Expected: PASS 所有 TestFacts_* 测试

- [ ] **Step 5: 写 Rules 测试**

```go
// pkg/api/rules_test.go
package api

import "testing"

type mockRule struct {
    name     string
    priority int
}

func (r *mockRule) GetName() string        { return r.name }
func (r *mockRule) GetDescription() string { return "" }
func (r *mockRule) GetPriority() int       { return r.priority }
func (r *mockRule) Evaluate(_ *Facts) bool { return false }
func (r *mockRule) Execute(_ *Facts) error { return nil }

func TestRules_RegisterAndSize(t *testing.T) {
    rules := NewRules()
    rules.Register(&mockRule{name: "r1", priority: 1})
    rules.Register(&mockRule{name: "r2", priority: 2})
    if rules.Size() != 2 {
        t.Errorf("expected 2, got %d", rules.Size())
    }
}

func TestRules_SortedByPriority(t *testing.T) {
    rules := NewRules()
    rules.Register(&mockRule{name: "r2", priority: 2})
    rules.Register(&mockRule{name: "r1", priority: 1})
    names := []string{}
    for _, r := range rules.Slice() {
        names = append(names, r.GetName())
    }
    if names[0] != "r1" || names[1] != "r2" {
        t.Errorf("expected [r1, r2], got %v", names)
    }
}

func TestRules_SamePrioritySortedByName(t *testing.T) {
    rules := NewRules()
    rules.Register(&mockRule{name: "z", priority: 1})
    rules.Register(&mockRule{name: "a", priority: 1})
    names := []string{}
    for _, r := range rules.Slice() {
        names = append(names, r.GetName())
    }
    if names[0] != "a" || names[1] != "z" {
        t.Errorf("expected [a, z], got %v", names)
    }
}

func TestRules_UnregisterByName(t *testing.T) {
    rules := NewRules()
    rules.Register(&mockRule{name: "r1", priority: 1})
    rules.Unregister("r1")
    if rules.Size() != 0 {
        t.Error("expected empty after unregister")
    }
}

func TestRules_IsEmpty(t *testing.T) {
    rules := NewRules()
    if !rules.IsEmpty() {
        t.Error("expected empty")
    }
    rules.Register(&mockRule{name: "r1", priority: 1})
    if rules.IsEmpty() {
        t.Error("expected not empty")
    }
}
```

- [ ] **Step 6: 运行测试，确认失败**

```bash
go test ./pkg/api/... -run TestRules -v
```

Expected: FAIL — `NewRules undefined`

- [ ] **Step 7: 实现 pkg/api/rules.go**

```go
package api

import "sort"

type Rules struct {
    rules []Rule
}

func NewRules(rules ...Rule) *Rules {
    r := &Rules{}
    r.Register(rules...)
    return r
}

func (rs *Rules) Register(rules ...Rule) {
    for _, r := range rules {
        rs.rules = append(rs.rules, r)
    }
    rs.sort()
}

func (rs *Rules) Unregister(name string) {
    for i, r := range rs.rules {
        if r.GetName() == name {
            rs.rules = append(rs.rules[:i], rs.rules[i+1:]...)
            return
        }
    }
}

func (rs *Rules) IsEmpty() bool {
    return len(rs.rules) == 0
}

func (rs *Rules) Size() int {
    return len(rs.rules)
}

func (rs *Rules) Slice() []Rule {
    result := make([]Rule, len(rs.rules))
    copy(result, rs.rules)
    return result
}

func (rs *Rules) sort() {
    sort.SliceStable(rs.rules, func(i, j int) bool {
        pi, pj := rs.rules[i].GetPriority(), rs.rules[j].GetPriority()
        if pi != pj {
            return pi < pj
        }
        return rs.rules[i].GetName() < rs.rules[j].GetName()
    })
}
```

- [ ] **Step 8: 运行测试，确认通过**

```bash
go test ./pkg/api/... -v
```

Expected: PASS 所有测试

- [ ] **Step 9: 验证整包编译**

```bash
go build ./pkg/api/...
```

Expected: 无错误

- [ ] **Step 10: Commit**

```bash
git add pkg/api/
git commit --author="maxingfang <maxingfang@126.com>" -m "feat(api): add Facts, Rules with tests"
```

---

## Task 4: pkg/core — BasicRule 和 RuleBuilder

**Files:**
- Create: `pkg/core/basic_rule.go`
- Create: `pkg/core/basic_rule_test.go`
- Create: `pkg/core/default_rule.go`
- Create: `pkg/core/rule_builder.go`
- Create: `pkg/core/rule_builder_test.go`

- [ ] **Step 1: 写 BasicRule 测试**

```go
// pkg/core/basic_rule_test.go
package core

import (
    "testing"
    "github.com/manfredma/easy-rules-4-go/pkg/api"
)

func TestBasicRule_Defaults(t *testing.T) {
    r := &BasicRule{}
    if r.GetName() != api.DefaultRuleName {
        t.Errorf("expected %s, got %s", api.DefaultRuleName, r.GetName())
    }
    if r.GetDescription() != api.DefaultRuleDescription {
        t.Errorf("expected %s, got %s", api.DefaultRuleDescription, r.GetDescription())
    }
    if r.GetPriority() != api.DefaultRulePriority {
        t.Errorf("expected %d, got %d", api.DefaultRulePriority, r.GetPriority())
    }
}

func TestBasicRule_CustomValues(t *testing.T) {
    r := NewBasicRule("my-rule", "desc", 5)
    if r.GetName() != "my-rule" {
        t.Errorf("expected my-rule, got %s", r.GetName())
    }
    if r.GetPriority() != 5 {
        t.Errorf("expected 5, got %d", r.GetPriority())
    }
}

func TestBasicRule_EvaluateReturnsFalse(t *testing.T) {
    r := &BasicRule{}
    if r.Evaluate(api.NewFacts()) {
        t.Error("default Evaluate should return false")
    }
}

func TestBasicRule_ExecuteReturnsNil(t *testing.T) {
    r := &BasicRule{}
    if err := r.Execute(api.NewFacts()); err != nil {
        t.Errorf("default Execute should return nil, got %v", err)
    }
}

// EmbedRule 模拟用户嵌入 BasicRule 并覆盖方法
type EmbedRule struct {
    BasicRule
}

func (e *EmbedRule) Evaluate(_ *api.Facts) bool { return true }
func (e *EmbedRule) Execute(_ *api.Facts) error { return nil }

func TestBasicRule_Embed(t *testing.T) {
    r := &EmbedRule{}
    r.BasicRule = *NewBasicRule("embed", "embed rule", 1)
    if r.GetName() != "embed" {
        t.Errorf("expected embed, got %s", r.GetName())
    }
    if !r.Evaluate(api.NewFacts()) {
        t.Error("expected true from embedded evaluate")
    }
}
```

- [ ] **Step 2: 运行测试，确认失败**

```bash
go test ./pkg/core/... -run TestBasicRule -v
```

Expected: FAIL — `BasicRule undefined`

- [ ] **Step 3: 实现 pkg/core/basic_rule.go**

```go
package core

import "github.com/manfredma/easy-rules-4-go/pkg/api"

type BasicRule struct {
    name        string
    description string
    priority    int
}

func NewBasicRule(name, description string, priority int) *BasicRule {
    return &BasicRule{name: name, description: description, priority: priority}
}

func (r *BasicRule) GetName() string        { return r.name }
func (r *BasicRule) GetDescription() string { return r.description }
func (r *BasicRule) GetPriority() int       { return r.priority }
func (r *BasicRule) Evaluate(_ *api.Facts) bool { return false }
func (r *BasicRule) Execute(_ *api.Facts) error { return nil }

func (r *BasicRule) SetName(name string)               { r.name = name }
func (r *BasicRule) SetDescription(description string) { r.description = description }
func (r *BasicRule) SetPriority(priority int)          { r.priority = priority }
```

- [ ] **Step 4: 写 RuleBuilder 测试**

```go
// pkg/core/rule_builder_test.go
package core

import (
    "errors"
    "testing"
    "github.com/manfredma/easy-rules-4-go/pkg/api"
)

func TestRuleBuilder_BuildWithConditionAndAction(t *testing.T) {
    facts := api.NewFacts()
    facts.Put("x", 10)

    rule := NewRuleBuilder().
        Name("test-rule").
        Description("desc").
        Priority(1).
        When(func(f *api.Facts) bool {
            return f.Get("x").(int) > 5
        }).
        Then(func(f *api.Facts) error {
            f.Put("result", "big")
            return nil
        }).
        Build()

    if rule.GetName() != "test-rule" {
        t.Errorf("expected test-rule, got %s", rule.GetName())
    }
    if !rule.Evaluate(facts) {
        t.Error("expected condition to be true")
    }
    _ = rule.Execute(facts)
    if facts.Get("result") != "big" {
        t.Error("expected result=big after execute")
    }
}

func TestRuleBuilder_DefaultsWhenNotSet(t *testing.T) {
    rule := NewRuleBuilder().Build()
    if rule.GetName() != api.DefaultRuleName {
        t.Errorf("expected default name, got %s", rule.GetName())
    }
    if rule.GetPriority() != api.DefaultRulePriority {
        t.Errorf("expected default priority, got %d", rule.GetPriority())
    }
    if rule.Evaluate(api.NewFacts()) {
        t.Error("default condition should return false")
    }
}

func TestRuleBuilder_MultipleActions(t *testing.T) {
    facts := api.NewFacts()
    facts.Put("count", 0)

    rule := NewRuleBuilder().
        When(api.ConditionTrue).
        Then(func(f *api.Facts) error { f.Put("count", 1); return nil }).
        Then(func(f *api.Facts) error { f.Put("count", f.Get("count").(int)+1); return nil }).
        Build()

    _ = rule.Execute(facts)
    if facts.Get("count") != 2 {
        t.Errorf("expected count=2, got %v", facts.Get("count"))
    }
}

func TestRuleBuilder_ActionError(t *testing.T) {
    rule := NewRuleBuilder().
        When(api.ConditionTrue).
        Then(func(_ *api.Facts) error { return errors.New("boom") }).
        Build()

    err := rule.Execute(api.NewFacts())
    if err == nil || err.Error() != "boom" {
        t.Errorf("expected error boom, got %v", err)
    }
}
```

- [ ] **Step 5: 运行测试，确认失败**

```bash
go test ./pkg/core/... -run TestRuleBuilder -v
```

Expected: FAIL — `NewRuleBuilder undefined`

- [ ] **Step 6: 实现 pkg/core/default_rule.go**

```go
package core

import "github.com/manfredma/easy-rules-4-go/pkg/api"

type DefaultRule struct {
    BasicRule
    condition api.Condition
    actions   []api.Action
}

func (r *DefaultRule) Evaluate(facts *api.Facts) bool {
    return r.condition(facts)
}

func (r *DefaultRule) Execute(facts *api.Facts) error {
    for _, action := range r.actions {
        if err := action(facts); err != nil {
            return err
        }
    }
    return nil
}
```

- [ ] **Step 7: 实现 pkg/core/rule_builder.go**

```go
package core

import "github.com/manfredma/easy-rules-4-go/pkg/api"

type RuleBuilder struct {
    name        string
    description string
    priority    int
    condition   api.Condition
    actions     []api.Action
}

func NewRuleBuilder() *RuleBuilder {
    return &RuleBuilder{
        name:        api.DefaultRuleName,
        description: api.DefaultRuleDescription,
        priority:    api.DefaultRulePriority,
        condition:   api.ConditionFalse,
    }
}

func (b *RuleBuilder) Name(name string) *RuleBuilder {
    b.name = name
    return b
}

func (b *RuleBuilder) Description(desc string) *RuleBuilder {
    b.description = desc
    return b
}

func (b *RuleBuilder) Priority(priority int) *RuleBuilder {
    b.priority = priority
    return b
}

func (b *RuleBuilder) When(condition api.Condition) *RuleBuilder {
    b.condition = condition
    return b
}

func (b *RuleBuilder) Then(action api.Action) *RuleBuilder {
    b.actions = append(b.actions, action)
    return b
}

func (b *RuleBuilder) Build() api.Rule {
    r := &DefaultRule{
        condition: b.condition,
        actions:   b.actions,
    }
    r.BasicRule = *NewBasicRule(b.name, b.description, b.priority)
    return r
}
```

- [ ] **Step 8: 运行所有 core 测试，确认通过**

```bash
go test ./pkg/core/... -v
```

Expected: PASS 所有测试

- [ ] **Step 9: Commit**

```bash
git add pkg/core/
git commit --author="maxingfang <maxingfang@126.com>" -m "feat(core): add BasicRule, DefaultRule, RuleBuilder with tests"
```

---

## Task 5: pkg/core — DefaultRulesEngine

**Files:**
- Create: `pkg/core/default_engine.go`
- Create: `pkg/core/default_engine_test.go`

- [ ] **Step 1: 写 DefaultRulesEngine 测试**

```go
// pkg/core/default_engine_test.go
package core

import (
    "errors"
    "testing"
    "github.com/manfredma/easy-rules-4-go/pkg/api"
)

func makeCountRule(name string, priority int, result *int, increment int) api.Rule {
    return NewRuleBuilder().
        Name(name).
        Priority(priority).
        When(api.ConditionTrue).
        Then(func(f *api.Facts) error { *result += increment; return nil }).
        Build()
}

func TestDefaultEngine_FiresMatchingRules(t *testing.T) {
    count := 0
    rules := api.NewRules(
        makeCountRule("r1", 1, &count, 1),
        makeCountRule("r2", 2, &count, 10),
    )
    engine := NewDefaultRulesEngine(api.NewEngineParameters())
    engine.Fire(rules, api.NewFacts())
    if count != 11 {
        t.Errorf("expected 11, got %d", count)
    }
}

func TestDefaultEngine_SkipsNonMatchingRules(t *testing.T) {
    count := 0
    rules := api.NewRules(
        NewRuleBuilder().Name("r1").When(api.ConditionFalse).Then(func(_ *api.Facts) error { count++; return nil }).Build(),
    )
    engine := NewDefaultRulesEngine(api.NewEngineParameters())
    engine.Fire(rules, api.NewFacts())
    if count != 0 {
        t.Errorf("expected 0, got %d", count)
    }
}

func TestDefaultEngine_SkipOnFirstAppliedRule(t *testing.T) {
    count := 0
    rules := api.NewRules(
        makeCountRule("r1", 1, &count, 1),
        makeCountRule("r2", 2, &count, 10),
    )
    params := api.NewEngineParameters()
    params.SkipOnFirstAppliedRule = true
    engine := NewDefaultRulesEngine(params)
    engine.Fire(rules, api.NewFacts())
    if count != 1 {
        t.Errorf("expected 1 (stop after first applied), got %d", count)
    }
}

func TestDefaultEngine_SkipOnFirstNonTriggeredRule(t *testing.T) {
    count := 0
    rules := api.NewRules(
        NewRuleBuilder().Name("r1").Priority(1).When(api.ConditionFalse).Then(func(_ *api.Facts) error { count++; return nil }).Build(),
        makeCountRule("r2", 2, &count, 10),
    )
    params := api.NewEngineParameters()
    params.SkipOnFirstNonTriggeredRule = true
    engine := NewDefaultRulesEngine(params)
    engine.Fire(rules, api.NewFacts())
    if count != 0 {
        t.Errorf("expected 0 (stop after first non-triggered), got %d", count)
    }
}

func TestDefaultEngine_SkipOnFirstFailedRule(t *testing.T) {
    count := 0
    rules := api.NewRules(
        NewRuleBuilder().Name("r1").Priority(1).When(api.ConditionTrue).Then(func(_ *api.Facts) error { return errors.New("fail") }).Build(),
        makeCountRule("r2", 2, &count, 10),
    )
    params := api.NewEngineParameters()
    params.SkipOnFirstFailedRule = true
    engine := NewDefaultRulesEngine(params)
    engine.Fire(rules, api.NewFacts())
    if count != 0 {
        t.Errorf("expected 0 (stop after first failed), got %d", count)
    }
}

func TestDefaultEngine_PriorityThreshold(t *testing.T) {
    count := 0
    rules := api.NewRules(
        makeCountRule("r1", 1, &count, 1),
        makeCountRule("r2", 10, &count, 100),
    )
    params := api.NewEngineParameters()
    params.PriorityThreshold = 5
    engine := NewDefaultRulesEngine(params)
    engine.Fire(rules, api.NewFacts())
    if count != 1 {
        t.Errorf("expected 1 (r2 skipped by threshold), got %d", count)
    }
}

func TestDefaultEngine_Check(t *testing.T) {
    rules := api.NewRules(
        NewRuleBuilder().Name("yes").When(api.ConditionTrue).Build(),
        NewRuleBuilder().Name("no").When(api.ConditionFalse).Build(),
    )
    engine := NewDefaultRulesEngine(api.NewEngineParameters())
    result := engine.Check(rules, api.NewFacts())
    if len(result) != 2 {
        t.Errorf("expected 2 entries, got %d", len(result))
    }
    for r, v := range result {
        if r.GetName() == "yes" && !v {
            t.Error("expected yes=true")
        }
        if r.GetName() == "no" && v {
            t.Error("expected no=false")
        }
    }
}

func TestDefaultEngine_RuleListener(t *testing.T) {
    var beforeEvalCalled, onSuccessCalled bool

    listener := &testRuleListener{
        beforeEvaluate: func(_ api.Rule, _ *api.Facts) bool { beforeEvalCalled = true; return true },
        onSuccess:      func(_ api.Rule, _ *api.Facts) { onSuccessCalled = true },
    }

    rules := api.NewRules(NewRuleBuilder().Name("r1").When(api.ConditionTrue).Then(func(_ *api.Facts) error { return nil }).Build())
    engine := NewDefaultRulesEngine(api.NewEngineParameters())
    engine.RegisterRuleListener(listener)
    engine.Fire(rules, api.NewFacts())

    if !beforeEvalCalled {
        t.Error("expected BeforeEvaluate to be called")
    }
    if !onSuccessCalled {
        t.Error("expected OnSuccess to be called")
    }
}

// testRuleListener is a test helper
type testRuleListener struct {
    api.DefaultRuleListener
    beforeEvaluate func(api.Rule, *api.Facts) bool
    onSuccess      func(api.Rule, *api.Facts)
    onFailure      func(api.Rule, *api.Facts, error)
}

func (l *testRuleListener) BeforeEvaluate(r api.Rule, f *api.Facts) bool {
    if l.beforeEvaluate != nil { return l.beforeEvaluate(r, f) }
    return true
}
func (l *testRuleListener) OnSuccess(r api.Rule, f *api.Facts) {
    if l.onSuccess != nil { l.onSuccess(r, f) }
}
func (l *testRuleListener) OnFailure(r api.Rule, f *api.Facts, err error) {
    if l.onFailure != nil { l.onFailure(r, f, err) }
}
```

- [ ] **Step 2: 运行测试，确认失败**

```bash
go test ./pkg/core/... -run TestDefaultEngine -v
```

Expected: FAIL — `NewDefaultRulesEngine undefined`

- [ ] **Step 3: 实现 pkg/core/default_engine.go**

```go
package core

import (
    "fmt"
    "github.com/manfredma/easy-rules-4-go/pkg/api"
)

type DefaultRulesEngine struct {
    parameters          *api.EngineParameters
    ruleListeners       []api.RuleListener
    engineListeners     []api.RulesEngineListener
}

func NewDefaultRulesEngine(params *api.EngineParameters) *DefaultRulesEngine {
    if params == nil {
        params = api.NewEngineParameters()
    }
    return &DefaultRulesEngine{parameters: params}
}

func (e *DefaultRulesEngine) GetParameters() *api.EngineParameters { return e.parameters }
func (e *DefaultRulesEngine) GetRuleListeners() []api.RuleListener { return e.ruleListeners }
func (e *DefaultRulesEngine) GetEngineListeners() []api.RulesEngineListener { return e.engineListeners }

func (e *DefaultRulesEngine) RegisterRuleListener(l api.RuleListener) {
    e.ruleListeners = append(e.ruleListeners, l)
}

func (e *DefaultRulesEngine) RegisterEngineListener(l api.RulesEngineListener) {
    e.engineListeners = append(e.engineListeners, l)
}

func (e *DefaultRulesEngine) Fire(rules *api.Rules, facts *api.Facts) {
    e.triggerEngineListenersBefore(rules, facts)
    e.doFire(rules, facts)
    e.triggerEngineListenersAfter(rules, facts)
}

func (e *DefaultRulesEngine) doFire(rules *api.Rules, facts *api.Facts) {
    if rules.IsEmpty() {
        return
    }
    for _, rule := range rules.Slice() {
        if rule.GetPriority() > e.parameters.PriorityThreshold {
            break
        }
        if !e.shouldBeEvaluated(rule, facts) {
            continue
        }
        evaluationResult := false
        var evalErr error
        func() {
            defer func() {
                if rec := recover(); rec != nil {
                    evalErr = fmt.Errorf("panic: %v", rec)
                }
            }()
            evaluationResult = rule.Evaluate(facts)
        }()
        if evalErr != nil {
            e.triggerOnEvaluationError(rule, facts, evalErr)
            if e.parameters.SkipOnFirstNonTriggeredRule {
                break
            }
            continue
        }
        if evaluationResult {
            e.triggerAfterEvaluate(rule, facts, true)
            e.triggerBeforeExecute(rule, facts)
            err := rule.Execute(facts)
            if err != nil {
                e.triggerOnFailure(rule, facts, err)
                if e.parameters.SkipOnFirstFailedRule {
                    break
                }
            } else {
                e.triggerOnSuccess(rule, facts)
                if e.parameters.SkipOnFirstAppliedRule {
                    break
                }
            }
        } else {
            e.triggerAfterEvaluate(rule, facts, false)
            if e.parameters.SkipOnFirstNonTriggeredRule {
                break
            }
        }
    }
}

func (e *DefaultRulesEngine) Check(rules *api.Rules, facts *api.Facts) map[api.Rule]bool {
    e.triggerEngineListenersBefore(rules, facts)
    result := make(map[api.Rule]bool)
    for _, rule := range rules.Slice() {
        if e.shouldBeEvaluated(rule, facts) {
            result[rule] = rule.Evaluate(facts)
        }
    }
    e.triggerEngineListenersAfter(rules, facts)
    return result
}

func (e *DefaultRulesEngine) shouldBeEvaluated(rule api.Rule, facts *api.Facts) bool {
    for _, l := range e.ruleListeners {
        if !l.BeforeEvaluate(rule, facts) {
            return false
        }
    }
    return true
}

func (e *DefaultRulesEngine) triggerAfterEvaluate(rule api.Rule, facts *api.Facts, result bool) {
    for _, l := range e.ruleListeners { l.AfterEvaluate(rule, facts, result) }
}

func (e *DefaultRulesEngine) triggerOnEvaluationError(rule api.Rule, facts *api.Facts, err error) {
    for _, l := range e.ruleListeners { l.OnEvaluationError(rule, facts, err) }
}

func (e *DefaultRulesEngine) triggerBeforeExecute(rule api.Rule, facts *api.Facts) {
    for _, l := range e.ruleListeners { l.BeforeExecute(rule, facts) }
}

func (e *DefaultRulesEngine) triggerOnSuccess(rule api.Rule, facts *api.Facts) {
    for _, l := range e.ruleListeners { l.OnSuccess(rule, facts) }
}

func (e *DefaultRulesEngine) triggerOnFailure(rule api.Rule, facts *api.Facts, err error) {
    for _, l := range e.ruleListeners { l.OnFailure(rule, facts, err) }
}

func (e *DefaultRulesEngine) triggerEngineListenersBefore(rules *api.Rules, facts *api.Facts) {
    for _, l := range e.engineListeners { l.BeforeEvaluate(rules, facts) }
}

func (e *DefaultRulesEngine) triggerEngineListenersAfter(rules *api.Rules, facts *api.Facts) {
    for _, l := range e.engineListeners { l.AfterExecute(rules, facts) }
}
```

- [ ] **Step 4: 运行测试，确认通过**

```bash
go test ./pkg/core/... -run TestDefaultEngine -v
```

Expected: PASS 所有 TestDefaultEngine_* 测试

- [ ] **Step 5: Commit**

```bash
git add pkg/core/default_engine.go pkg/core/default_engine_test.go
git commit --author="maxingfang <maxingfang@126.com>" -m "feat(core): add DefaultRulesEngine with tests"
```

---

## Task 6: pkg/core — InferenceRulesEngine

**Files:**
- Create: `pkg/core/inference_engine.go`
- Create: `pkg/core/inference_engine_test.go`

- [ ] **Step 1: 写 InferenceRulesEngine 测试**

```go
// pkg/core/inference_engine_test.go
package core

import (
    "testing"
    "github.com/manfredma/easy-rules-4-go/pkg/api"
)

func TestInferenceEngine_LoopsUntilNoCandidate(t *testing.T) {
    // counter 从 0 增加到 3，每次循环 +1，条件是 counter < 3
    facts := api.NewFacts()
    facts.Put("counter", 0)

    rule := NewRuleBuilder().
        Name("increment").
        When(func(f *api.Facts) bool { return f.Get("counter").(int) < 3 }).
        Then(func(f *api.Facts) error {
            f.Put("counter", f.Get("counter").(int)+1)
            return nil
        }).
        Build()

    rules := api.NewRules(rule)
    engine := NewInferenceRulesEngine(api.NewEngineParameters())
    engine.Fire(rules, facts)

    if facts.Get("counter").(int) != 3 {
        t.Errorf("expected counter=3, got %v", facts.Get("counter"))
    }
}

func TestInferenceEngine_StopsWhenNoCandidates(t *testing.T) {
    count := 0
    rules := api.NewRules(
        NewRuleBuilder().Name("r1").When(api.ConditionFalse).Then(func(_ *api.Facts) error { count++; return nil }).Build(),
    )
    engine := NewInferenceRulesEngine(api.NewEngineParameters())
    engine.Fire(rules, api.NewFacts())
    if count != 0 {
        t.Errorf("expected 0 executions, got %d", count)
    }
}
```

- [ ] **Step 2: 运行测试，确认失败**

```bash
go test ./pkg/core/... -run TestInferenceEngine -v
```

Expected: FAIL — `NewInferenceRulesEngine undefined`

- [ ] **Step 3: 实现 pkg/core/inference_engine.go**

```go
package core

import "github.com/manfredma/easy-rules-4-go/pkg/api"

type InferenceRulesEngine struct {
    delegate *DefaultRulesEngine
}

func NewInferenceRulesEngine(params *api.EngineParameters) *InferenceRulesEngine {
    return &InferenceRulesEngine{delegate: NewDefaultRulesEngine(params)}
}

func (e *InferenceRulesEngine) GetParameters() *api.EngineParameters {
    return e.delegate.GetParameters()
}

func (e *InferenceRulesEngine) GetRuleListeners() []api.RuleListener {
    return e.delegate.GetRuleListeners()
}

func (e *InferenceRulesEngine) GetEngineListeners() []api.RulesEngineListener {
    return e.delegate.GetEngineListeners()
}

func (e *InferenceRulesEngine) RegisterRuleListener(l api.RuleListener) {
    e.delegate.RegisterRuleListener(l)
}

func (e *InferenceRulesEngine) RegisterEngineListener(l api.RulesEngineListener) {
    e.delegate.RegisterEngineListener(l)
}

func (e *InferenceRulesEngine) Fire(rules *api.Rules, facts *api.Facts) {
    for {
        candidates := e.selectCandidates(rules, facts)
        if candidates.IsEmpty() {
            break
        }
        e.delegate.Fire(candidates, facts)
    }
}

func (e *InferenceRulesEngine) Check(rules *api.Rules, facts *api.Facts) map[api.Rule]bool {
    return e.delegate.Check(rules, facts)
}

func (e *InferenceRulesEngine) selectCandidates(rules *api.Rules, facts *api.Facts) *api.Rules {
    candidates := api.NewRules()
    for _, rule := range rules.Slice() {
        if rule.Evaluate(facts) {
            candidates.Register(rule)
        }
    }
    return candidates
}
```

- [ ] **Step 4: 运行所有 core 测试，确认通过**

```bash
go test ./pkg/core/... -v
```

Expected: PASS 全部测试

- [ ] **Step 5: Commit**

```bash
git add pkg/core/inference_engine.go pkg/core/inference_engine_test.go
git commit --author="maxingfang <maxingfang@126.com>" -m "feat(core): add InferenceRulesEngine with tests"
```

---

## Task 7: pkg/composite — 组合规则

**Files:**
- Create: `pkg/composite/composite_rule.go`
- Create: `pkg/composite/unit_rule_group.go`
- Create: `pkg/composite/unit_rule_group_test.go`
- Create: `pkg/composite/activation_rule_group.go`
- Create: `pkg/composite/activation_rule_group_test.go`
- Create: `pkg/composite/conditional_rule_group.go`
- Create: `pkg/composite/conditional_rule_group_test.go`

- [ ] **Step 1: 实现 pkg/composite/composite_rule.go**

```go
package composite

import (
    "github.com/manfredma/easy-rules-4-go/pkg/api"
    "github.com/manfredma/easy-rules-4-go/pkg/core"
    "sort"
)

type CompositeRule struct {
    core.BasicRule
    rules []api.Rule
}

func NewCompositeRule(name, description string, priority int) *CompositeRule {
    c := &CompositeRule{}
    c.BasicRule = *core.NewBasicRule(name, description, priority)
    return c
}

func (c *CompositeRule) AddRule(rule api.Rule) {
    c.rules = append(c.rules, rule)
    sort.SliceStable(c.rules, func(i, j int) bool {
        pi, pj := c.rules[i].GetPriority(), c.rules[j].GetPriority()
        if pi != pj { return pi < pj }
        return c.rules[i].GetName() < c.rules[j].GetName()
    })
}

func (c *CompositeRule) RemoveRule(name string) {
    for i, r := range c.rules {
        if r.GetName() == name {
            c.rules = append(c.rules[:i], c.rules[i+1:]...)
            return
        }
    }
}

func (c *CompositeRule) GetRules() []api.Rule {
    result := make([]api.Rule, len(c.rules))
    copy(result, c.rules)
    return result
}
```

- [ ] **Step 2: 写 UnitRuleGroup 测试**

```go
// pkg/composite/unit_rule_group_test.go
package composite

import (
    "fmt"
    "testing"
    "github.com/manfredma/easy-rules-4-go/pkg/api"
    "github.com/manfredma/easy-rules-4-go/pkg/core"
)

func trueRule(name string) api.Rule {
    return core.NewRuleBuilder().Name(name).When(api.ConditionTrue).Then(func(_ *api.Facts) error { return nil }).Build()
}

func falseRule(name string) api.Rule {
    return core.NewRuleBuilder().Name(name).When(api.ConditionFalse).Then(func(_ *api.Facts) error { return nil }).Build()
}

func TestUnitRuleGroup_AllTrueEvaluatesTrue(t *testing.T) {
    g := NewUnitRuleGroup("g", "", 1)
    g.AddRule(trueRule("r1"))
    g.AddRule(trueRule("r2"))
    if !g.Evaluate(api.NewFacts()) {
        t.Error("expected true when all rules are true")
    }
}

func TestUnitRuleGroup_OneFalseEvaluatesFalse(t *testing.T) {
    g := NewUnitRuleGroup("g", "", 1)
    g.AddRule(trueRule("r1"))
    g.AddRule(falseRule("r2"))
    if g.Evaluate(api.NewFacts()) {
        t.Error("expected false when one rule is false")
    }
}

func TestUnitRuleGroup_ExecutesAllRules(t *testing.T) {
    count := 0
    g := NewUnitRuleGroup("g", "", 1)
    for i := 0; i < 3; i++ {
        g.AddRule(core.NewRuleBuilder().Name(fmt.Sprintf("r%d", i)).When(api.ConditionTrue).Then(func(_ *api.Facts) error { count++; return nil }).Build())
    }
    _ = g.Execute(api.NewFacts())
    if count != 3 {
        t.Errorf("expected 3 executions, got %d", count)
    }
}

func TestUnitRuleGroup_EmptyEvaluatesFalse(t *testing.T) {
    g := NewUnitRuleGroup("g", "", 1)
    if g.Evaluate(api.NewFacts()) {
        t.Error("expected false for empty group")
    }
}
```

- [ ] **Step 3: 运行测试，确认失败**

```bash
go test ./pkg/composite/... -run TestUnitRuleGroup -v
```

Expected: FAIL — `NewUnitRuleGroup undefined`

- [ ] **Step 4: 实现 pkg/composite/unit_rule_group.go**

```go
package composite

import "github.com/manfredma/easy-rules-4-go/pkg/api"

type UnitRuleGroup struct {
    CompositeRule
}

func NewUnitRuleGroup(name, description string, priority int) *UnitRuleGroup {
    g := &UnitRuleGroup{}
    g.CompositeRule = *NewCompositeRule(name, description, priority)
    return g
}

func (g *UnitRuleGroup) Evaluate(facts *api.Facts) bool {
    if len(g.rules) == 0 {
        return false
    }
    for _, r := range g.rules {
        if !r.Evaluate(facts) {
            return false
        }
    }
    return true
}

func (g *UnitRuleGroup) Execute(facts *api.Facts) error {
    for _, r := range g.rules {
        if err := r.Execute(facts); err != nil {
            return err
        }
    }
    return nil
}
```

- [ ] **Step 5: 写 ActivationRuleGroup 测试**

```go
// pkg/composite/activation_rule_group_test.go
package composite

import (
    "fmt"
    "testing"
    "github.com/manfredma/easy-rules-4-go/pkg/api"
    "github.com/manfredma/easy-rules-4-go/pkg/core"
)

func TestActivationRuleGroup_FiresFirstMatchingRule(t *testing.T) {
    fired := ""
    g := NewActivationRuleGroup("g", "", 1)
    g.AddRule(core.NewRuleBuilder().Name("r1").Priority(1).When(api.ConditionTrue).Then(func(_ *api.Facts) error { fired = "r1"; return nil }).Build())
    g.AddRule(core.NewRuleBuilder().Name("r2").Priority(2).When(api.ConditionTrue).Then(func(_ *api.Facts) error { fired = "r2"; return nil }).Build())

    if !g.Evaluate(api.NewFacts()) {
        t.Error("expected true")
    }
    _ = g.Execute(api.NewFacts())
    if fired != "r1" {
        t.Errorf("expected r1 (highest priority), got %s", fired)
    }
}

func TestActivationRuleGroup_ReturnsFalseWhenNoneMatch(t *testing.T) {
    g := NewActivationRuleGroup("g", "", 1)
    g.AddRule(falseRule("r1"))
    if g.Evaluate(api.NewFacts()) {
        t.Error("expected false")
    }
}

func TestActivationRuleGroup_SkipsLowerPriorityRules(t *testing.T) {
    count := 0
    g := NewActivationRuleGroup("g", "", 1)
    g.AddRule(core.NewRuleBuilder().Name("r1").Priority(1).When(api.ConditionTrue).Then(func(_ *api.Facts) error { count++; return nil }).Build())
    g.AddRule(core.NewRuleBuilder().Name("r2").Priority(2).When(api.ConditionTrue).Then(func(_ *api.Facts) error { count++; return nil }).Build())

    g.Evaluate(api.NewFacts())
    _ = g.Execute(api.NewFacts())
    if count != 1 {
        t.Errorf("expected only 1 rule fired (XOR), got %d", count)
    }
}
```

- [ ] **Step 6: 实现 pkg/composite/activation_rule_group.go**

```go
package composite

import "github.com/manfredma/easy-rules-4-go/pkg/api"

type ActivationRuleGroup struct {
    CompositeRule
    selectedRule api.Rule
}

func NewActivationRuleGroup(name, description string, priority int) *ActivationRuleGroup {
    g := &ActivationRuleGroup{}
    g.CompositeRule = *NewCompositeRule(name, description, priority)
    return g
}

func (g *ActivationRuleGroup) Evaluate(facts *api.Facts) bool {
    for _, r := range g.rules {
        if r.Evaluate(facts) {
            g.selectedRule = r
            return true
        }
    }
    return false
}

func (g *ActivationRuleGroup) Execute(facts *api.Facts) error {
    if g.selectedRule != nil {
        return g.selectedRule.Execute(facts)
    }
    return nil
}
```

- [ ] **Step 7: 写 ConditionalRuleGroup 测试**

```go
// pkg/composite/conditional_rule_group_test.go
package composite

import (
    "testing"
    "github.com/manfredma/easy-rules-4-go/pkg/api"
    "github.com/manfredma/easy-rules-4-go/pkg/core"
)

func TestConditionalRuleGroup_GatePassesExecutesOthers(t *testing.T) {
    count := 0
    g := NewConditionalRuleGroup("g", "", 1)
    // priority=0 是条件门（最高优先级）
    g.AddRule(core.NewRuleBuilder().Name("gate").Priority(0).When(api.ConditionTrue).Then(func(_ *api.Facts) error { count++; return nil }).Build())
    g.AddRule(core.NewRuleBuilder().Name("r1").Priority(1).When(api.ConditionTrue).Then(func(_ *api.Facts) error { count++; return nil }).Build())
    g.AddRule(core.NewRuleBuilder().Name("r2").Priority(2).When(api.ConditionTrue).Then(func(_ *api.Facts) error { count++; return nil }).Build())

    if !g.Evaluate(api.NewFacts()) {
        t.Error("expected true when gate passes")
    }
    _ = g.Execute(api.NewFacts())
    if count != 3 {
        t.Errorf("expected 3 (gate + r1 + r2), got %d", count)
    }
}

func TestConditionalRuleGroup_GateFailsSkipsAll(t *testing.T) {
    count := 0
    g := NewConditionalRuleGroup("g", "", 1)
    g.AddRule(core.NewRuleBuilder().Name("gate").Priority(0).When(api.ConditionFalse).Then(func(_ *api.Facts) error { count++; return nil }).Build())
    g.AddRule(core.NewRuleBuilder().Name("r1").Priority(1).When(api.ConditionTrue).Then(func(_ *api.Facts) error { count++; return nil }).Build())

    if g.Evaluate(api.NewFacts()) {
        t.Error("expected false when gate fails")
    }
    if count != 0 {
        t.Errorf("expected 0, got %d", count)
    }
}

func TestConditionalRuleGroup_OnlyGateFired_WhenOthersFalse(t *testing.T) {
    count := 0
    g := NewConditionalRuleGroup("g", "", 1)
    g.AddRule(core.NewRuleBuilder().Name("gate").Priority(0).When(api.ConditionTrue).Then(func(_ *api.Facts) error { count++; return nil }).Build())
    g.AddRule(core.NewRuleBuilder().Name("r1").Priority(1).When(api.ConditionFalse).Then(func(_ *api.Facts) error { count++; return nil }).Build())

    g.Evaluate(api.NewFacts())
    _ = g.Execute(api.NewFacts())
    if count != 1 {
        t.Errorf("expected only gate to fire (count=1), got %d", count)
    }
}

func TestConditionalRuleGroup_PanicsOnTwoRulesWithSameHighestPriority(t *testing.T) {
    g := NewConditionalRuleGroup("g", "", 1)
    g.AddRule(core.NewRuleBuilder().Name("r1").Priority(0).When(api.ConditionTrue).Build())
    g.AddRule(core.NewRuleBuilder().Name("r2").Priority(0).When(api.ConditionTrue).Build())

    defer func() {
        if r := recover(); r == nil {
            t.Error("expected panic for duplicate highest priority")
        }
    }()
    g.Evaluate(api.NewFacts())
}
```

- [ ] **Step 8: 实现 pkg/composite/conditional_rule_group.go**

```go
package composite

import (
    "github.com/manfredma/easy-rules-4-go/pkg/api"
    "sort"
)

type ConditionalRuleGroup struct {
    CompositeRule
    conditionalRule      api.Rule
    successfulEvaluations []api.Rule
}

func NewConditionalRuleGroup(name, description string, priority int) *ConditionalRuleGroup {
    g := &ConditionalRuleGroup{}
    g.CompositeRule = *NewCompositeRule(name, description, priority)
    return g
}

func (g *ConditionalRuleGroup) Evaluate(facts *api.Facts) bool {
    g.successfulEvaluations = nil
    g.conditionalRule = g.getRuleWithHighestPriority()
    if g.conditionalRule.Evaluate(facts) {
        for _, r := range g.rules {
            if r != g.conditionalRule && r.Evaluate(facts) {
                g.successfulEvaluations = append(g.successfulEvaluations, r)
            }
        }
        return true
    }
    return false
}

func (g *ConditionalRuleGroup) Execute(facts *api.Facts) error {
    if err := g.conditionalRule.Execute(facts); err != nil {
        return err
    }
    sorted := make([]api.Rule, len(g.successfulEvaluations))
    copy(sorted, g.successfulEvaluations)
    sort.SliceStable(sorted, func(i, j int) bool {
        pi, pj := sorted[i].GetPriority(), sorted[j].GetPriority()
        if pi != pj { return pi < pj }
        return sorted[i].GetName() < sorted[j].GetName()
    })
    for _, r := range sorted {
        if err := r.Execute(facts); err != nil {
            return err
        }
    }
    return nil
}

func (g *ConditionalRuleGroup) getRuleWithHighestPriority() api.Rule {
    if len(g.rules) == 0 {
        panic("ConditionalRuleGroup has no rules")
    }
    sorted := make([]api.Rule, len(g.rules))
    copy(sorted, g.rules)
    sort.SliceStable(sorted, func(i, j int) bool {
        return sorted[i].GetPriority() < sorted[j].GetPriority()
    })
    highest := sorted[0]
    if len(sorted) > 1 && sorted[1].GetPriority() == highest.GetPriority() {
        panic("ConditionalRuleGroup: only one rule can have the highest priority")
    }
    return highest
}
```

- [ ] **Step 9: 运行所有 composite 测试**

```bash
go test ./pkg/composite/... -v
```

Expected: PASS 全部测试

- [ ] **Step 10: Commit**

```bash
git add pkg/composite/
git commit --author="maxingfang <maxingfang@126.com>" -m "feat(composite): add UnitRuleGroup, ActivationRuleGroup, ConditionalRuleGroup with tests"
```

---

## Task 8: pkg/expr — 表达式引擎集成

**Files:**
- Create: `pkg/expr/expr_condition.go`
- Create: `pkg/expr/expr_action.go`
- Create: `pkg/expr/expr_rule.go`
- Create: `pkg/expr/expr_rule_test.go`

- [ ] **Step 1: 写表达式引擎测试**

```go
// pkg/expr/expr_rule_test.go
package expr

import (
    "testing"
    "github.com/manfredma/easy-rules-4-go/pkg/api"
)

func TestExprCondition_TrueExpression(t *testing.T) {
    facts := api.NewFacts()
    facts.Put("rain", true)

    cond, err := NewExprCondition("rain == true")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !cond(facts) {
        t.Error("expected condition to be true")
    }
}

func TestExprCondition_FalseExpression(t *testing.T) {
    facts := api.NewFacts()
    facts.Put("rain", false)

    cond, err := NewExprCondition("rain == true")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if cond(facts) {
        t.Error("expected condition to be false")
    }
}

func TestExprCondition_InvalidExpression(t *testing.T) {
    _, err := NewExprCondition("!!!invalid(((")
    if err == nil {
        t.Error("expected error for invalid expression")
    }
}

func TestExprAction_ModifyFacts(t *testing.T) {
    facts := api.NewFacts()
    facts.Put("count", 0)

    action, err := NewExprAction("count + 1", "count")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    _ = action(facts)
    if facts.Get("count").(int) != 1 {
        t.Errorf("expected count=1, got %v", facts.Get("count"))
    }
}

func TestExprRule_EvaluateAndExecute(t *testing.T) {
    facts := api.NewFacts()
    facts.Put("temperature", 35)
    facts.Put("cooled", false)

    rule, err := NewExprRuleBuilder().
        Name("heat-rule").
        When("temperature > 30").
        Then("!cooled", "cooled").
        Build()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if !rule.Evaluate(facts) {
        t.Error("expected evaluate to be true")
    }
    _ = rule.Execute(facts)
    if facts.Get("cooled") != true {
        t.Errorf("expected cooled=true, got %v", facts.Get("cooled"))
    }
}

func TestExprRule_InvalidCondition(t *testing.T) {
    _, err := NewExprRuleBuilder().
        Name("bad").
        When("!!!bad").
        Build()
    if err == nil {
        t.Error("expected error for invalid condition")
    }
}
```

- [ ] **Step 2: 运行测试，确认失败**

```bash
go test ./pkg/expr/... -v
```

Expected: FAIL — `NewExprCondition undefined`

- [ ] **Step 3: 实现 pkg/expr/expr_condition.go**

```go
package expr

import (
    "github.com/expr-lang/expr"
    "github.com/manfredma/easy-rules-4-go/pkg/api"
)

func NewExprCondition(expression string) (api.Condition, error) {
    program, err := expr.Compile(expression, expr.Env(map[string]any{}), expr.AsBool())
    if err != nil {
        return nil, err
    }
    return func(facts *api.Facts) bool {
        result, err := expr.Run(program, facts.AsMap())
        if err != nil {
            return false
        }
        v, _ := result.(bool)
        return v
    }, nil
}
```

- [ ] **Step 4: 实现 pkg/expr/expr_action.go**

```go
package expr

import (
    "github.com/expr-lang/expr"
    "github.com/manfredma/easy-rules-4-go/pkg/api"
)

// NewExprAction compiles an expression and stores its result back into facts under outputKey.
func NewExprAction(expression string, outputKey string) (api.Action, error) {
    program, err := expr.Compile(expression, expr.Env(map[string]any{}))
    if err != nil {
        return nil, err
    }
    return func(facts *api.Facts) error {
        result, err := expr.Run(program, facts.AsMap())
        if err != nil {
            return err
        }
        facts.Put(outputKey, result)
        return nil
    }, nil
}
```

- [ ] **Step 5: 实现 pkg/expr/expr_rule.go**

```go
package expr

import (
    "github.com/manfredma/easy-rules-4-go/pkg/api"
    "github.com/manfredma/easy-rules-4-go/pkg/core"
)

type ExprRule struct {
    core.BasicRule
    condition api.Condition
    actions   []api.Action
}

func (r *ExprRule) Evaluate(facts *api.Facts) bool {
    return r.condition(facts)
}

func (r *ExprRule) Execute(facts *api.Facts) error {
    for _, action := range r.actions {
        if err := action(facts); err != nil {
            return err
        }
    }
    return nil
}

type ExprRuleBuilder struct {
    name        string
    description string
    priority    int
    condition   string
    actions     []struct{ expr, key string }
}

func NewExprRuleBuilder() *ExprRuleBuilder {
    return &ExprRuleBuilder{
        name:        api.DefaultRuleName,
        description: api.DefaultRuleDescription,
        priority:    api.DefaultRulePriority,
    }
}

func (b *ExprRuleBuilder) Name(name string) *ExprRuleBuilder {
    b.name = name
    return b
}

func (b *ExprRuleBuilder) Description(desc string) *ExprRuleBuilder {
    b.description = desc
    return b
}

func (b *ExprRuleBuilder) Priority(priority int) *ExprRuleBuilder {
    b.priority = priority
    return b
}

func (b *ExprRuleBuilder) When(expression string) *ExprRuleBuilder {
    b.condition = expression
    return b
}

// Then adds an action: expression result is stored in outputKey of facts.
func (b *ExprRuleBuilder) Then(expression, outputKey string) *ExprRuleBuilder {
    b.actions = append(b.actions, struct{ expr, key string }{expression, outputKey})
    return b
}

func (b *ExprRuleBuilder) Build() (*ExprRule, error) {
    cond, err := NewExprCondition(b.condition)
    if err != nil {
        return nil, err
    }
    var actions []api.Action
    for _, a := range b.actions {
        action, err := NewExprAction(a.expr, a.key)
        if err != nil {
            return nil, err
        }
        actions = append(actions, action)
    }
    r := &ExprRule{condition: cond, actions: actions}
    r.BasicRule = *core.NewBasicRule(b.name, b.description, b.priority)
    return r, nil
}
```

- [ ] **Step 6: 运行测试，确认通过**

```bash
go test ./pkg/expr/... -v
```

Expected: PASS 全部测试

- [ ] **Step 7: Commit**

```bash
git add pkg/expr/
git commit --author="maxingfang <maxingfang@126.com>" -m "feat(expr): add ExprCondition, ExprAction, ExprRule with tests"
```

---

## Task 9: pkg/reader — YAML/JSON 文件加载

**Files:**
- Create: `pkg/reader/rule_definition.go`
- Create: `pkg/reader/yaml_reader.go`
- Create: `pkg/reader/json_reader.go`
- Create: `pkg/reader/reader_test.go`
- Create: `pkg/reader/testdata/weather-rules.yml`
- Create: `pkg/reader/testdata/weather-rules.json`

- [ ] **Step 1: 创建测试数据文件**

```yaml
# pkg/reader/testdata/weather-rules.yml
- name: "weather rule"
  description: "if it rains then take an umbrella"
  priority: 1
  condition: "rain == true"
  actions:
    - expr: "!umbrella"
      outputKey: "umbrella"
```

```json
// pkg/reader/testdata/weather-rules.json
[
  {
    "name": "weather rule",
    "description": "if it rains then take an umbrella",
    "priority": 1,
    "condition": "rain == true",
    "actions": [
      {"expr": "!umbrella", "outputKey": "umbrella"}
    ]
  }
]
```

- [ ] **Step 2: 写 reader 测试**

```go
// pkg/reader/reader_test.go
package reader

import (
    "testing"
    "github.com/manfredma/easy-rules-4-go/pkg/api"
)

func TestYamlRuleFactory_LoadAndFire(t *testing.T) {
    factory := NewYamlRuleFactory()
    rules, err := factory.CreateRulesFrom("testdata/weather-rules.yml")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if rules.Size() != 1 {
        t.Errorf("expected 1 rule, got %d", rules.Size())
    }

    facts := api.NewFacts()
    facts.Put("rain", true)
    facts.Put("umbrella", false)

    rule := rules.Slice()[0]
    if rule.GetName() != "weather rule" {
        t.Errorf("expected 'weather rule', got %s", rule.GetName())
    }
    if !rule.Evaluate(facts) {
        t.Error("expected condition to be true")
    }
    _ = rule.Execute(facts)
    if facts.Get("umbrella") != true {
        t.Errorf("expected umbrella=true after execute, got %v", facts.Get("umbrella"))
    }
}

func TestJsonRuleFactory_LoadAndFire(t *testing.T) {
    factory := NewJsonRuleFactory()
    rules, err := factory.CreateRulesFrom("testdata/weather-rules.json")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if rules.Size() != 1 {
        t.Errorf("expected 1 rule, got %d", rules.Size())
    }

    facts := api.NewFacts()
    facts.Put("rain", true)
    facts.Put("umbrella", false)

    rule := rules.Slice()[0]
    if !rule.Evaluate(facts) {
        t.Error("expected condition to be true")
    }
    _ = rule.Execute(facts)
    if facts.Get("umbrella") != true {
        t.Errorf("expected umbrella=true, got %v", facts.Get("umbrella"))
    }
}

func TestYamlRuleFactory_FileNotFound(t *testing.T) {
    factory := NewYamlRuleFactory()
    _, err := factory.CreateRulesFrom("testdata/nonexistent.yml")
    if err == nil {
        t.Error("expected error for missing file")
    }
}
```

- [ ] **Step 3: 运行测试，确认失败**

```bash
go test ./pkg/reader/... -v
```

Expected: FAIL — `NewYamlRuleFactory undefined`

- [ ] **Step 4: 实现 pkg/reader/rule_definition.go**

```go
package reader

type ActionDefinition struct {
    Expr      string `yaml:"expr" json:"expr"`
    OutputKey string `yaml:"outputKey" json:"outputKey"`
}

type RuleDefinition struct {
    Name        string             `yaml:"name" json:"name"`
    Description string             `yaml:"description" json:"description"`
    Priority    int                `yaml:"priority" json:"priority"`
    Condition   string             `yaml:"condition" json:"condition"`
    Actions     []ActionDefinition `yaml:"actions" json:"actions"`
}
```

- [ ] **Step 5: 实现 pkg/reader/yaml_reader.go**

```go
package reader

import (
    "os"
    "gopkg.in/yaml.v3"
    "github.com/manfredma/easy-rules-4-go/pkg/api"
    exprpkg "github.com/manfredma/easy-rules-4-go/pkg/expr"
)

type YamlRuleFactory struct{}

func NewYamlRuleFactory() *YamlRuleFactory { return &YamlRuleFactory{} }

func (f *YamlRuleFactory) CreateRulesFrom(path string) (*api.Rules, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var defs []RuleDefinition
    if err := yaml.Unmarshal(data, &defs); err != nil {
        return nil, err
    }
    return buildRules(defs)
}

func buildRules(defs []RuleDefinition) (*api.Rules, error) {
    rules := api.NewRules()
    for _, def := range defs {
        b := exprpkg.NewExprRuleBuilder().
            Name(def.Name).
            Description(def.Description).
            Priority(def.Priority).
            When(def.Condition)
        for _, a := range def.Actions {
            b.Then(a.Expr, a.OutputKey)
        }
        rule, err := b.Build()
        if err != nil {
            return nil, err
        }
        rules.Register(rule)
    }
    return rules, nil
}
```

- [ ] **Step 6: 实现 pkg/reader/json_reader.go**

```go
package reader

import (
    "encoding/json"
    "os"
    "github.com/manfredma/easy-rules-4-go/pkg/api"
)

type JsonRuleFactory struct{}

func NewJsonRuleFactory() *JsonRuleFactory { return &JsonRuleFactory{} }

func (f *JsonRuleFactory) CreateRulesFrom(path string) (*api.Rules, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var defs []RuleDefinition
    if err := json.Unmarshal(data, &defs); err != nil {
        return nil, err
    }
    return buildRules(defs)
}
```

- [ ] **Step 7: 运行测试，确认通过**

```bash
go test ./pkg/reader/... -v
```

Expected: PASS 全部测试

- [ ] **Step 8: Commit**

```bash
git add pkg/reader/
git commit --author="maxingfang <maxingfang@126.com>" -m "feat(reader): add YAML/JSON rule factory with tests"
```

---

## Task 10: examples — 示例程序

**Files:**
- Create: `examples/hello_world/main.go`
- Create: `examples/weather/main.go`
- Create: `examples/fizzbuzz/main.go`

- [ ] **Step 1: 实现 examples/hello_world/main.go**

```go
package main

import (
    "fmt"
    "github.com/manfredma/easy-rules-4-go/pkg/api"
    "github.com/manfredma/easy-rules-4-go/pkg/core"
)

func main() {
    rule := core.NewRuleBuilder().
        Name("hello world rule").
        Description("say hello world").
        When(api.ConditionTrue).
        Then(func(_ *api.Facts) error {
            fmt.Println("Hello World!")
            return nil
        }).
        Build()

    rules := api.NewRules(rule)
    engine := core.NewDefaultRulesEngine(api.NewEngineParameters())
    engine.Fire(rules, api.NewFacts())
}
```

- [ ] **Step 2: 实现 examples/weather/main.go（struct 嵌入方式）**

```go
package main

import (
    "fmt"
    "github.com/manfredma/easy-rules-4-go/pkg/api"
    "github.com/manfredma/easy-rules-4-go/pkg/core"
)

type WeatherRule struct {
    core.BasicRule
}

func (r *WeatherRule) Evaluate(facts *api.Facts) bool {
    rain, _ := facts.Get("rain").(bool)
    return rain
}

func (r *WeatherRule) Execute(_ *api.Facts) error {
    fmt.Println("It rains, take an umbrella!")
    return nil
}

func main() {
    rule := &WeatherRule{}
    rule.BasicRule = *core.NewBasicRule("weather rule", "if it rains then take an umbrella", 1)

    facts := api.NewFacts()
    facts.Put("rain", true)

    rules := api.NewRules(rule)
    engine := core.NewDefaultRulesEngine(api.NewEngineParameters())
    engine.Fire(rules, facts)
}
```

- [ ] **Step 3: 实现 examples/fizzbuzz/main.go（表达式引擎方式）**

```go
package main

import (
    "fmt"
    "github.com/manfredma/easy-rules-4-go/pkg/api"
    "github.com/manfredma/easy-rules-4-go/pkg/core"
    easyrules_expr "github.com/manfredma/easy-rules-4-go/pkg/expr"
)

func main() {
    fizzRule, _ := easyrules_expr.NewExprRuleBuilder().
        Name("fizz rule").
        Priority(1).
        When(`number % 3 == 0`).
        Then(`"fizz"`, "print").
        Build()

    buzzRule, _ := easyrules_expr.NewExprRuleBuilder().
        Name("buzz rule").
        Priority(2).
        When(`number % 5 == 0`).
        Then(`"buzz"`, "print").
        Build()

    fizzBuzzRule, _ := easyrules_expr.NewExprRuleBuilder().
        Name("fizzbuzz rule").
        Priority(0).
        When(`number % 3 == 0 && number % 5 == 0`).
        Then(`"fizzbuzz"`, "print").
        Build()

    numberRule := core.NewRuleBuilder().
        Name("number rule").
        Priority(3).
        When(api.ConditionTrue).
        Then(func(facts *api.Facts) error {
            if facts.Get("print") == nil {
                fmt.Println(facts.Get("number"))
            } else {
                fmt.Println(facts.Get("print"))
                facts.Remove("print")
            }
            return nil
        }).
        Build()

    rules := api.NewRules(fizzBuzzRule, fizzRule, buzzRule, numberRule)
    engine := core.NewDefaultRulesEngine(api.NewEngineParameters())

    for i := 1; i <= 20; i++ {
        facts := api.NewFacts()
        facts.Put("number", i)
        engine.Fire(rules, facts)
    }
}
```

- [ ] **Step 4: 验证示例编译和运行**

```bash
go run examples/hello_world/main.go
```
Expected: 输出 `Hello World!`

```bash
go run examples/weather/main.go
```
Expected: 输出 `It rains, take an umbrella!`

```bash
go run examples/fizzbuzz/main.go
```
Expected: 输出 1~20 的 fizzbuzz 结果

- [ ] **Step 5: Commit**

```bash
git add examples/
git commit --author="maxingfang <maxingfang@126.com>" -m "feat(examples): add hello_world, weather, fizzbuzz examples"
```

---

## Task 11: 全量测试 + 覆盖率检查

- [ ] **Step 1: 运行全量测试**

```bash
go test ./... -v
```

Expected: PASS 全部测试，无编译错误

- [ ] **Step 2: 检查覆盖率**

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total
```

Expected: 总覆盖率 ≥ 80%

- [ ] **Step 3: 如覆盖率不足，查看未覆盖行**

```bash
go tool cover -html=coverage.out -o coverage.html
open coverage.html
```

补充缺少的测试用例直到覆盖率 ≥ 80%。

- [ ] **Step 4: 最终 Commit**

```bash
git add -A
git commit --author="maxingfang <maxingfang@126.com>" -m "test: ensure ≥80% test coverage across all packages"
```
