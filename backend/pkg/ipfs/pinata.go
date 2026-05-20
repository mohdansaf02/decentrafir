package ipfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type PinataClient struct {
	jwt    string
	apiKey string
	secret string
}

type pinataResponse struct {
	IpfsHash string `json:"IpfsHash"`
	PinSize  int    `json:"PinSize"`
}

func NewPinataClient(jwt, apiKey, secret string) *PinataClient {
	return &PinataClient{jwt: jwt, apiKey: apiKey, secret: secret}
}

func (p *PinataClient) UploadFile(filename string, data []byte) (string, string, error) {
	if p.jwt == "" && p.apiKey == "" {
		// dev fallback: return deterministic mock CID
		cid := fmt.Sprintf("bafydev%s", hashBytes(data))
		return cid, "https://ipfs.io/ipfs/" + cid, nil
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", "", err
	}
	if _, err := part.Write(data); err != nil {
		return "", "", err
	}
	_ = writer.Close()

	req, err := http.NewRequest(http.MethodPost, "https://api.pinata.cloud/pinning/pinFileToIPFS", body)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if p.jwt != "" {
		req.Header.Set("Authorization", "Bearer "+p.jwt)
	} else {
		req.Header.Set("pinata_api_key", p.apiKey)
		req.Header.Set("pinata_secret_api_key", p.secret)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", "", fmt.Errorf("pinata upload failed: %s", string(respBody))
	}

	var result pinataResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", "", err
	}
	url := "https://gateway.pinata.cloud/ipfs/" + result.IpfsHash
	return result.IpfsHash, url, nil
}

func hashBytes(b []byte) string {
	h := uint32(2166136261)
	for _, c := range b {
		h ^= uint32(c)
		h *= 16777619
	}
	return fmt.Sprintf("%08x", h)
}
