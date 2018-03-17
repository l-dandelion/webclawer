package module

type Type string

const (
	TYPE_DOWNLOADER Type = "downloader"
	TYPE_ANALYZER   Type = "analyzer"
	TYPE_PIPELINE   Type = "pipeline"
)

var legalLetterTypeMap = map[string]Type{
	"D": TYPE_DOWNLOADER,
	"A": TYPE_ANALYZER,
	"P": TYPE_PIPELINE,
}

var legalTypeLetterMap = map[Type]string{
	TYPE_DOWNLOADER: "D",
	TYPE_ANALYZER:   "A",
	TYPE_PIPELINE:   "P",
}

func GetType(mid MID) (bool, Type) {
	parts, err := SplitMID(mid)
	if err != nil {
		return false, ""
	}
	mt, ok := legalLetterTypeMap[parts[0]]
	return ok, mt
}

func getLetter(moduleType Type) (bool, string) {
	var letter string
	var found bool
	for l, t := range legalLetterTypeMap {
		if t == moduleType {
			letter = l
			found = true
			break
		}
	}
	return found, letter
}

func CheckType(moduleType Type, module Module) bool {
	if moduleType == "" || module == nil {
		return false
	}
	switch moduleType {
	case TYPE_DOWNLOADER:
		if _, ok := module.(Downloader); ok {
			return true
		}
	case TYPE_ANALYZER:
		if _, ok := module.(Analyzer); ok {
			return true
		}
	case TYPE_PIPELINE:
		if _, ok := module.(Pipeline); ok {
			return true
		}
	}
	return false
}

func LegalType(moduleType Type) bool {
	_, ok := legalTypeLetterMap[moduleType]
	return ok
}

func LegalLetter(letter string) bool {
	_, ok := legalLetterTypeMap[letter]
	return ok
}

func typeToLetter(moduleType Type) (bool, string) {
	letter, ok := legalTypeLetterMap[moduleType]
	return ok, letter
}

func letterToType(letter string) (bool, Type) {
	moduleType, ok := legalLetterTypeMap[letter]
	return ok, moduleType
}
