package examples

import (
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"testing"
)

const (
	Rule7 = `
rule UserTestRule7 "test 7"  salience 10{
	when
	  User.Age > 1
	then
	  User.SetName("FromRule");
	  Retract("UserTestRule7");
}
`
)

type AUserIssue7 struct {
	Name string
	Age  int
}

func (u *AUserIssue7) GetName() string {
	return u.Name
}

func (u *AUserIssue7) SetName(name interface{}) {
	u.Name = name.(string)
}

func TestMethodCall_Issue7(t *testing.T) {
	user := &AUserIssue7{
		Name: "Watson",
		Age:  7,
	}

	dataContext := ast.NewDataContext()
	err := dataContext.Add("User", user)
	if err != nil {
		t.Fatal(err)
	}

	// Prepare knowledgebase library and load it with our rule.
	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	err = rb.BuildRuleFromResource("Test", "0.1.1", pkg.NewBytesResource([]byte(Rule7)))
	if err != nil {
		t.Log(err)
	} else {
		eng1 := &engine.GruleEngine{MaxCycle: 5}
		kb := lib.NewKnowledgeBaseInstance("Test", "0.1.1")
		err := eng1.Execute(dataContext, kb)
		if err != nil {
			t.Fatal(err)
		}
		if user.GetName() != "FromRule" {
			t.Fatalf("User should be FromRule but %s", user.GetName())
		}
	}
}
