package utils

type CheckNames struct {
	OrgName string
	Name    string
}

var checkNames []CheckNames = []CheckNames{
	{
		OrgName: `АКЦИОНЕРНОЕ ОБЩЕСТВО "ТАНДЕР"`,
		Name:    "магнит",
	},
	{
		OrgName: `ОБЩЕСТВО С ОГРАНИЧЕННОЙ ОТВЕТСТВЕННОСТЬЮ "ДНС РИТЕЙЛ"`,
		Name:    "ДНС",
	},
	{
		OrgName: `ООО "О'КЕЙ"`,
		Name:    "ОКЕЙ",
	},
	{
		OrgName: `ООО "Камелот-А"`,
		Name:    "Ярче",
	},
}

func AssumeCheckName(orgName string) string {
	names := make(map[string]string)
	for _, pair := range checkNames {
		names[pair.OrgName] = pair.Name
	}

	if name, ok := names[orgName]; ok {
		return name
	}

	return orgName
}
