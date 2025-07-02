package clause

type DropTag struct {
	IfNotExist bool
	TagName    string
}

const DropTagName = "DROP_TAG"

func (dt DropTag) Name() string {
	return DropTagName
}

func (dt DropTag) MergeIn(clause *Clause) {
	clause.Expression = dt
}

func (dt DropTag) Build(nGQL Builder) error {
	nGQL.WriteString("DROP TAG ")
	if dt.IfNotExist {
		nGQL.WriteString("IF NOT EXISTS ")
	}
	nGQL.WriteString(dt.TagName)
	return nil
}
