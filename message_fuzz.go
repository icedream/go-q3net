// +build gofuzz

package quake

func Fuzz(data []byte) int {
	parsedMsg, err := UnmarshalMessage(data)
	if err != nil {
		if parsedMsg != nil {
			panic("parsedMsg was not nil in error case")
		}
		return 0
	}
	return 1
}
