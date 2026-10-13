package objectrepo_test

import (
	"testing"

	"github.com/cvartan/appcontext/pkg/objectrepo"
)

type BaseInterface interface {
	Log(t *testing.T, value string)
}

type BaseStruct struct {
	index string
}

func (s *BaseStruct) Log(t *testing.T, value string) {
	t.Logf("base struct applied %s for %s\n", value, s.index)
}

type SecondStruct struct{}

func (s *SecondStruct) Log(t *testing.T, value string) {
	t.Logf("second struct with value %s", value)
}

type MegaStruct struct {
	L1 BaseInterface `component:""`
	L2 *BaseStruct   `component:""`
}

type MegaStructWithTags struct {
	L1 BaseInterface `component:"BaseStruct"`
	L2 BaseInterface `component:"SecondStruct"`
}

type MegaStructWithSlices struct {
	A1 []BaseInterface `component:""`
	A2 []*BaseStruct   `component:""`
}

func TestSimpleStructBind(t *testing.T) {

	or := objectrepo.New()
	or.Bind(&BaseStruct{index: "test1"})
	m := &MegaStruct{}
	or.Bind(m)
	or.Apply()

	m.L1.Log(t, "L1")
	m.L2.Log(t, "L2")
}

func TestTaggedComponents(t *testing.T) {

	or := objectrepo.New()
	or.Bind(&BaseStruct{index: "L1"})
	or.Bind(&SecondStruct{})

	m := &MegaStructWithTags{}
	or.Bind(m)

	m.L1.Log(t, "T1")
	m.L2.Log(t, "T2")
}

func TestAttributesWithSlices(t *testing.T) {
	or := objectrepo.New()
	or.Bind(&BaseStruct{index: "L1"})
	or.Bind(&SecondStruct{})

	s := &MegaStructWithSlices{}
	or.Bind(s)

	for _, v := range s.A1 {
		v.Log(t, "A1")
	}

	for _, v := range s.A2 {
		v.Log(t, "A2")
	}

}
