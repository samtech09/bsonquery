package bsonquery

type condition struct {
	operator  string
	value     interface{}
	fieldName string
	options   string
}

//C creates a new Condition
func C() *condition {
	return &condition{}
}

func (c *condition) EQ(fieldname string, val interface{}) condition {
	c.operator = copEQ
	c.fieldName = fieldname
	c.value = val
	return *c
}

func (c *condition) GT(fieldname string, val interface{}) condition {
	c.operator = copGT
	c.fieldName = fieldname
	c.value = val
	return *c
}
func (c *condition) GTE(fieldname string, val interface{}) condition {
	c.operator = copGTE
	c.fieldName = fieldname
	c.value = val
	return *c
}
func (c *condition) IN(fieldname string, val interface{}) condition {
	c.operator = copIN
	c.fieldName = fieldname
	c.value = val
	return *c
}
func (c *condition) LT(fieldname string, val interface{}) condition {
	c.operator = copLT
	c.fieldName = fieldname
	c.value = val
	return *c
}
func (c *condition) LTE(fieldname string, val interface{}) condition {
	c.operator = copLTE
	c.fieldName = fieldname
	c.value = val
	return *c
}
func (c *condition) NE(fieldname string, val interface{}) condition {
	c.operator = copNE
	c.fieldName = fieldname
	c.value = val
	return *c
}
func (c *condition) NIN(fieldname string, val interface{}) condition {
	c.operator = copNIN
	c.fieldName = fieldname
	c.value = val
	return *c
}
func (c *condition) Regex(fieldname, expression string, ignorecase bool) condition {
	c.operator = copRegex
	c.fieldName = fieldname
	c.value = expression
	if ignorecase {
		c.options = "i"
	}
	return *c
}
func (c *condition) Exist(fieldname string, expression bool) condition {
	c.operator = copExists
	c.fieldName = fieldname
	c.value = expression
	return *c
}
