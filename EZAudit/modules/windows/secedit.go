package windows

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"bytes"
	"errors"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
	"strings"
	"unicode/utf16"
)

func init() {
	registry.Register("secedit", &secedit{}, false, registry.Windows)
}

type secedit struct {
	pattern string
}

func (sc secedit) New(p map[string]interface{}) (global.Module, error) {
	// Parse arguments
	ok := true

	//Check if there are only allowed keys
	for k, _ := range p {
		switch k {
		case "pattern":
			sc.pattern, ok = p["pattern"].(string)
			if !ok {
				return nil, errors.New("the key \"pattern\" is set for the module. The value must be a \"string\"")
			}
		default:
			if k != "pattern" {
				return nil, errors.New("there is no key called: \"" + k + "\" in the module secedit")
			}
		}
	}

	return &sc, nil
}

func (sc *secedit) Execute(s *global.Step) (output string, err error) {
	//Lock because of the mainResource
	global.SeceditLock.Lock()
	defer global.SeceditLock.Unlock()

	//Check if Secedit file artifact already exists
	if _, errSeceditFile := os.Stat(global.SeceditPath); os.IsNotExist(errSeceditFile) {

		// files does not exist
		artifact.CreateFolderIfNotExist(global.ArtifactsDir)
		artifact.CreateFolderIfNotExist(global.MainResourcesDir)
		artifact.CreateFolderIfNotExist(global.SeceditDir)

		//Replacing slashes with backslashes
		global.SeceditPath = strings.Replace(global.SeceditPath, "/", "\\", -1)

		//Get the secedit file
		s, errS := interact.Program(s, "secedit.exe", "/export", "/cfg", global.SeceditPath)
		if errS != nil {
			return s, errS
		}
	}

	// artifact exists
	seceditFile, err := ReadFileUTF16(global.SeceditPath)
	output = string(seceditFile)
	if err != nil {
		return output, err
	}

	seceditMap := make(map[string]string)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		pair := strings.Split(line, " = ")
		if len(pair) > 1 {
			seceditMap[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
		}
	}

	output = seceditMap[sc.pattern]

	return
}

func ReadFileUTF16(filename string) ([]byte, error) {

	// Read the file into a []byte:
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Make an tranformer that converts MS-Win default to UTF8:
	win16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	// Make a transformer that is like win16be, but abides by BOM:
	utf16bom := unicode.BOMOverride(win16be.NewDecoder())

	// Make a Reader that uses utf16bom:
	unicodeReader := transform.NewReader(bytes.NewReader(raw), utf16bom)

	// decode and print:
	decoded, err := ioutil.ReadAll(unicodeReader)
	return decoded, err
}

func WriteFileUTF16(data string) error {
	var bytes [2]byte
	const BOM = '\ufffe' //LE. for BE '\ufeff'

	file, err := os.Create(global.SeceditPath)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes[0] = BOM >> 8
	bytes[1] = BOM & 255

	file.Write(bytes[0:])
	runes := utf16.Encode([]rune(data))
	for _, r := range runes {
		bytes[1] = byte(r >> 8)
		bytes[0] = byte(r & 255)
		file.Write(bytes[0:])
	}
	return nil
}
