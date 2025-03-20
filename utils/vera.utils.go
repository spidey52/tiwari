package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

func encryptAES(plainText string, key []byte) (string, string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	nonce := make([]byte, 12) // 12-byte nonce for GCM
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	cipherText := aesGCM.Seal(nil, nonce, []byte(plainText), nil)
	return hex.EncodeToString(cipherText), hex.EncodeToString(nonce), nil
}
func generateKey() ([]byte, error) {
	seed_key := time.Now().UTC().Format("2006-01-02-mgdh")
	hash := sha256.Sum256([]byte(seed_key))

	substr := hash[:16]
	fmt.Println("seed key", seed_key)

	return substr, nil

}
// Decrypt text using AES-GCM
func decryptAES(cipherTextHex string, nonceHex string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	cipherText, _ := hex.DecodeString(cipherTextHex)
	nonce, _ := hex.DecodeString(nonceHex)

	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}



func EncryptPassword(password string) (string, error) {

	key, err := generateKey()

	if err != nil {
		return "", err
	}

	encryptedPassword, nonce, err := encryptAES(password, key)

	if err != nil {
		return "", err
	}

	return encryptedPassword + "." + nonce, nil
}

func DecryptPassword(encryptedPassword string) (string, error) {

	key, err := generateKey()

	if err != nil {
		return "", err
	}

	password := strings.Split(encryptedPassword, ".")

	if len(password) != 2 {
		return "", errors.New("invalid password")
	}

	decryptedPassword, err := decryptAES(password[0], password[1], key)

	if err != nil {
		return "", err
	}

	return decryptedPassword, nil
}

// read json file
type VeraVolume struct {
	VolumeName  string `json:"volume_name"`
	DriveLetter string `json:"drive_letter"`
	VolumePath  string `json:"volume_path"`
}

type BackupVolume struct {
	VolumeName   string `json:"volume_name"`
	DriveLetter  string `json:"drive_letter"`
	VolumePath   string `json:"volume_path"`
	BackupFolder string `json:"backup_folder"`
}

type VeraCryptPath struct {
	VeraCryptPath string         `json:"vera_crypt_path"`
	Volumes       []VeraVolume   `json:"volumes"`
	BackupVolume  []BackupVolume `json:"backup_volume"`
}

func GetVeraDetails() (VeraCryptPath, error) {
	homedir, err := os.UserHomeDir()

	if err != nil {
		return VeraCryptPath{}, err
	}

	filePath := path.Join(homedir, ".vera-config.json")

	file, err := os.Open(filePath)

	if err != nil {
		return VeraCryptPath{}, err
	}

	defer file.Close()

	var veraCryptPath VeraCryptPath

	err = json.NewDecoder(file).Decode(&veraCryptPath)

	if err != nil {
		return VeraCryptPath{}, err
	}

	if veraCryptPath.VeraCryptPath == "" {
		return VeraCryptPath{}, errors.New("vera crypt path not found")
	}

	if len(veraCryptPath.Volumes) == 0 {
		return VeraCryptPath{}, errors.New("no volumes found")
	}

	return veraCryptPath, nil

}
