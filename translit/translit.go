package translit

var translitMap = map[string]string{
	"а": "a", "б": "b", "в": "v", "г": "g", "д": "d", "е": "e",
	"ё": "yo", "ж": "zh", "з": "z", "и": "i", "й": "y", "к": "k",
	"л": "l", "м": "m", "н": "n", "о": "o", "п": "p", "р": "r",
	"с": "s", "т": "t", "у": "u", "ф": "f", "х": "kh", "ц": "ts",
	"ч": "ch", "ш": "sh", "щ": "sch", "ъ": "", "ы": "i",
	"ь": "", "э": "e", "ю": "yu", "я": "ya",
	"А": "A", "Б": "B", "В": "V", "Г": "G", "Д": "D", "Е": "E",
	"Ё": "YO", "Ж": "ZH", "З": "Z", "И": "I", "Й": "Y", "К": "K",
	"Л": "L", "М": "M", "Н": "N", "О": "O", "П": "P", "Р": "R",
	"С": "S", "Т": "T", "У": "U", "Ф": "F", "Х": "KH", "Ц": "TS",
	"Ч": "CH", "Ш": "SH", "Щ": "SCH", "Ъ": "", "Ы": "I",
	"Ь": "", "Э": "E", "Ю": "YU", "Я": "YA",
}

func Translit(text string) string {
	var result string

	for _, r := range text {
		char, ok := translitMap[string(r)]
		if !ok {
			result += string(r)

			continue
		}

		result += char
	}

	return result
}