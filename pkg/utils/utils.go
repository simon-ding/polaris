package utils

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
	"golang.org/x/sys/unix"
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

func IsNameAcceptable(name1, name2 string) bool {
	re := regexp.MustCompile(`[^\p{L}\w\s]`)
	name1 = re.ReplaceAllString(strings.ToLower(name1), " ")
	name2 = re.ReplaceAllString(strings.ToLower(name2), " ")
	name1 = strings.Join(strings.Fields(name1), " ")
	name2 = strings.Join(strings.Fields(name2), " ")
	if strings.Contains(name1, name2) || strings.Contains(name2, name1) {
		return true
	}
	return false
}

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

func AvailableSpace(dir string) uint64 {
	var stat unix.Statfs_t

	unix.Statfs(dir, &stat)
	return stat.Bavail * uint64(stat.Bsize)
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
