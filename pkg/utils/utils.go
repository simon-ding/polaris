package utils

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"unicode"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

func IsASCII(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ContainsChineseChar(str string) bool {

	for _, r := range str {
		if unicode.Is(unicode.Han, r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// func IsNameAcceptable(name1, name2 string) bool {
// 	re := regexp.MustCompile(`[^\p{L}\w\s]`)
// 	name1 = re.ReplaceAllString(strings.ToLower(name1), " ")
// 	name2 = re.ReplaceAllString(strings.ToLower(name2), " ")
// 	name1 = strings.Join(strings.Fields(name1), " ")
// 	name2 = strings.Join(strings.Fields(name2), " ")
// 	if strings.Contains(name1, name2) || strings.Contains(name2, name1) {
// 		return true
// 	}
// 	return false
// }

func FindSeasonEpisodeNum(name string) (se int, ep int, err error) {
	seRe := regexp.MustCompile(`S\d+`)
	epRe := regexp.MustCompile(`E\d+`)
	nameUpper := strings.ToUpper(name)
	matchEp := epRe.FindAllString(nameUpper, -1)
	if len(matchEp) == 0 {
		err = errors.New("no episode num")
	}
	matchSe := seRe.FindAllString(nameUpper, -1)
	if len(matchSe) == 0 {
		err = errors.New("no season num")
	}
	if err != nil {
		return 0, 0, err
	}

	epNum := strings.TrimPrefix(matchEp[0], "E")
	epNum1, _ := strconv.Atoi(epNum)
	seNum := strings.TrimPrefix(matchSe[0], "S")
	seNum1, _ := strconv.Atoi(seNum)
	return seNum1, epNum1, nil
}

func FindSeasonPackageInfo(name string) (se int, err error) {
	seRe := regexp.MustCompile(`S\d+`)
	epRe := regexp.MustCompile(`E\d+`)
	nameUpper := strings.ToUpper(name)
	matchEp := epRe.FindAllString(nameUpper, -1)
	if len(matchEp) != 0 {
		err = errors.New("episode number should not exist")
	}
	matchSe := seRe.FindAllString(nameUpper, -1)
	if len(matchSe) == 0 {
		err = errors.New("no season num")
	}
	if err != nil {
		return 0, err
	}

	seNum := strings.TrimPrefix(matchSe[0], "S")
	se, _ = strconv.Atoi(seNum)
	return se, err
}

func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func SeasonId(seasonName string) (int, error) {
	//Season 01
	seRe := regexp.MustCompile(`\d+`)
	matchSe := seRe.FindAllString(seasonName, -1)
	if len(matchSe) == 0 {
		return 0, errors.New("no season number") //no season num
	}
	num, err := strconv.Atoi(matchSe[len(matchSe)-1])
	if err != nil {
		return 0, errors.Wrap(err, "convert")
	}
	return num, nil
}

func ChangeFileHash(name string) error {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0655)
	if err != nil {
		return errors.Wrap(err, "open file")
	}
	defer f.Close()
	_, err = f.Write([]byte("\000"))
	if err != nil {
		return errors.Wrap(err, "write file")
	}
	return nil
}

func TrimFields(v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	var mapSI map[string]interface{}
	if err := json.Unmarshal(bytes, &mapSI); err != nil {
		return err
	}
	mapSI = trimMapStringInterface(mapSI).(map[string]interface{})
	bytes2, err := json.Marshal(mapSI)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes2, v); err != nil {
		return err
	}
	return nil
}

func trimMapStringInterface(data interface{}) interface{} {
	if values, valid := data.([]interface{}); valid {
		for i := range values {
			data.([]interface{})[i] = trimMapStringInterface(values[i])
		}
	} else if values, valid := data.(map[string]interface{}); valid {
		for k, v := range values {
			data.(map[string]interface{})[k] = trimMapStringInterface(v)
		}
	} else if value, valid := data.(string); valid {
		data = strings.TrimSpace(value)
	}
	return data
}

// https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
// openURL opens the specified URL in the default browser of the user.
func OpenURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // "linux", "freebsd", "openbsd", "netbsd"
		// Check if running under WSL
		if isWSL() {
			// Use 'cmd.exe /c start' to open the URL in the default Windows browser
			cmd = "cmd.exe"
			args = []string{"/c", "start", url}
		} else {
			// Use xdg-open on native Linux environments
			cmd = "xdg-open"
			args = []string{url}
		}
	}
	if len(args) > 1 {
		// args[0] is used for 'start' command argument, to prevent issues with URLs starting with a quote
		args = append(args[:1], append([]string{""}, args[1:]...)...)
	}
	return exec.Command(cmd, args...).Start()
}

// isWSL checks if the Go program is running inside Windows Subsystem for Linux
func isWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}

func Link2Magnet(link string) (string, error) {
	if strings.HasPrefix(strings.ToLower(link), "magnet:") {
		return link, nil
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse //do not follow redirects
		},
	}
	
	resp, err := client.Get(link)
	if err != nil {
		return "", errors.Wrap(err, "get link")
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		//redirects
		tourl := resp.Header.Get("Location")
		return Link2Magnet(tourl)
	}
	info, err := metainfo.Load(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "parse response")
	}
	mg, err := info.MagnetV2()
	if err != nil {
		return "", errors.Wrap(err, "convert magnet")
	}
	return mg.String(), nil
}


func MagnetHash(link string) (string, error) {
	if mi, err := metainfo.ParseMagnetV2Uri(link); err != nil {
		return "", errors.Errorf("magnet link is not valid: %v", err)
	} else {
		hash := ""
		if mi.InfoHash.Unwrap().HexString() != "" {
			hash = mi.InfoHash.Unwrap().HexString()
		} else {
			btmh := mi.V2InfoHash.Unwrap()
			if btmh.HexString() != "" {
				hash = btmh.HexString()
			}
		}
		if hash == "" {
			return "", errors.Errorf("magnet has no info hash: %v", link)
		}
		return hash, nil
	}
}