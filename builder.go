package bsonquery

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type builder struct {
	conditionGroups map[int]conditionGroup
}

type conditionGroup struct {
	logicalOp  logicalOperator
	conditions []condition
}

type logicalOperator int

const (
	lopAnd logicalOperator = iota
	lopOr
	lopNor
	lopNot

	copEQ     = "$eq"
	copGT     = "$gt"
	copGTE    = "$gte"
	copIN     = "$in"
	copLT     = "$lt"
	copLTE    = "$lte"
	copNE     = "$ne"
	copNIN    = "$nin"
	copRegex  = "$regex"
	copExists = "$exists"
)

func Builder() *builder {
	b := &builder{}
	b.conditionGroups = make(map[int]conditionGroup)
	return b
}

func (b *builder) And(c ...condition) *builder {
	l := len(b.conditionGroups)
	b.conditionGroups[l] = makeConditionGroup(lopAnd, c...)
	return b
}

func (b *builder) Or(c ...condition) *builder {
	if len(c) < 2 {
		panic("OR logical operator require minimum two conditions")
	}
	l := len(b.conditionGroups)
	b.conditionGroups[l] = makeConditionGroup(lopOr, c...)
	return b
}

func (b *builder) Nor(c ...condition) *builder {
	if len(c) < 2 {
		panic("NOR logical operator require minimum two conditions")
	}
	l := len(b.conditionGroups)
	b.conditionGroups[l] = makeConditionGroup(lopNor, c...)
	return b
}

func (b *builder) Not(c condition) *builder {
	l := len(b.conditionGroups)
	b.conditionGroups[l] = makeConditionGroup(lopNor, c)
	return b
}

func (b *builder) Build() bson.M {
	m := bson.M{}
	for _, cg := range b.conditionGroups {
		switch cg.logicalOp {
		case lopAnd:
			for _, c := range cg.conditions {
				m[c.fieldName] = bson.M{c.operator: c.value}
			}

		case lopOr:
			m["$or"] = getArrayOfM(cg.conditions)

		case lopNor:
			m["$nor"] = getArrayOfM(cg.conditions)

		case lopNot:
			c := cg.conditions[0]
			m[c.fieldName] = bson.M{"$not": getM(c)}

		}
	}

	return m
}

func makeConditionGroup(op logicalOperator, c ...condition) conditionGroup {
	cg := conditionGroup{}
	cg.logicalOp = op
	cg.conditions = make([]condition, 0, len(c))
	for _, cd := range c {
		cg.conditions = append(cg.conditions, cd)
	}
	return cg
}

func getArrayOfM(cond []condition) []bson.M {
	ar := make([]bson.M, 0, len(cond))
	for _, c := range cond {
		//ar = append(ar, bson.M{c.fieldName: c.value})
		ar = append(ar, bson.M{c.fieldName: getM(c)})
	}
	return ar
}

func getM(c condition) bson.M {
	if c.operator == copRegex {
		//primitive.Regex{Pattern: "he", Options: ""}
		m := bson.M{c.operator: primitive.Regex{Pattern: c.value.(string), Options: c.options}}
		return m
	}
	m := bson.M{c.operator: c.value}
	return m
}
