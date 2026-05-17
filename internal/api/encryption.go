package api

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "encoding/json"
    "errors"
    "io"
    "os"
)

var encryptionKey = []byte("reconforge-32-byte-key-!!secret!!") // 32 bytes

func Encrypt(plaintext string) (string, error) {
    block, err := aes.NewCipher(encryptionKey)
    if err != nil {
        return "", err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encrypted string) (string, error) {
    data, err := base64.StdEncoding.DecodeString(encrypted)
    if err != nil {
        return "", err
    }
    
    block, err := aes.NewCipher(encryptionKey)
    if err != nil {
        return "", err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    
    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return "", errors.New("ciphertext too short")
    }
    
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }
    
    return string(plaintext), nil
}

func SaveAPIKeys(filename string, keys map[string]string) error {
    encrypted := make(map[string]string)
    for service, key := range keys {
        enc, err := Encrypt(key)
        if err != nil {
            return err
        }
        encrypted[service] = enc
    }
    
    data, err := json.MarshalIndent(encrypted, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(filename, data, 0600)
}

func LoadAPIKeys(filename string) (map[string]string, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    var encrypted map[string]string
    if err := json.Unmarshal(data, &encrypted); err != nil {
        return nil, err
    }
    
    keys := make(map[string]string)
    for service, enc := range encrypted {
        dec, err := Decrypt(enc)
        if err != nil {
            return nil, err
        }
        keys[service] = dec
    }
    
    return keys, nil
}
