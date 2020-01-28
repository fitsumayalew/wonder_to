package util

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	. "xCut/constants"
	"xCut/rtoken"
)




func ParseForm(w http.ResponseWriter,r *http.Request) bool{
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	}
	return true
}

func HashPassword(password string) ([]byte,error){
	return bcrypt.GenerateFromPassword([]byte(password), 12)
}
func ArePasswordsSame(hashedPassword string,rawPassword string)bool{
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
	return err != bcrypt.ErrMismatchedHashAndPassword
}


func IsParsableFormPost(w http.ResponseWriter, r *http.Request,csrfSignKey []byte) bool {
	return r.Method == http.MethodPost &&
		ParseForm(w, r) &&
		rtoken.IsCSRFValid(r.FormValue(CsrfKey), csrfSignKey)
}


func GenerateFileName(mf *multipart.File,filename string)(string,error){
	var returnMD5String string
	hash := md5.New()
	if _, err := io.Copy(hash, *mf); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	now := time.Now()
	year := strconv.Itoa(now.Year())
	month := now.Month().String()
	day := strconv.Itoa(now.Day())
	hour := strconv.Itoa(now.Hour())
	minute := strconv.Itoa(now.Minute())
	second := strconv.Itoa(now.Second())
	a := strings.Split(filename,".")
	filename = a[len(a)-1]
	return returnMD5String + "_" + year + month + day + "_" + hour + minute + second + "."+  filename, nil
}


func WriteFile(mf *multipart.File, fname string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := filepath.Join(wd, "ui", "assets", "img", fname)
	image, err := os.Create(path)
	if err != nil {
		return err
	}
	defer image.Close()
	_,err = io.Copy(image, *mf)
	return nil
}