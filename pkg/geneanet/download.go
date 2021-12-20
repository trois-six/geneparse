package geneanet

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	userAgent       = "GeneaNet v2.15 (Android 11 1080x2009@440)"
	loginURL        = "https://www.geneanet.org/connexion/verify.php?ctype=id"
	accountInfosURL = "https://www.geneanet.org/app/arbre/index.php?action=accountInfos&k=%s"
	loggedURL       = "https://www.geneanet.org/app/arbre/index.php?action=logged"
	importURL       = "https://www.geneanet.org/app/arbre/index.php?action=import"
	errNewRequest   = "newRequestWithContext %q: %w"
	errDoRequest    = "doing %q: %w"
	errReadBody     = "reading body: %w"
	errJSONMarshall = "json marshall: %w"
	errIOCopy       = "io copy: %w"
	randKLength     = 5
)

var (
	errStatusCode           = errors.New("authentication status code")
	errLoginNoSessionCookie = errors.New("no session Cookie")
)

type Geneanet struct {
	client    http.Client
	username  string
	password  string
	outputDir string
	timeout   time.Duration
	session   string
	nbPersons uint32
	sosa      uint32
	timestamp int64
}

// New initialize a Geneanet.
func New(username, password, outputDir string, timeout time.Duration) *Geneanet {
	return &Geneanet{
		client:    http.Client{},
		username:  username,
		password:  password,
		outputDir: outputDir,
		timeout:   timeout,
	}
}

func (g *Geneanet) Login(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()

	data := url.Values{
		"persistent": {"1"},
		"login":      {g.username},
		"password":   {g.password},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf(errNewRequest, loginURL, err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf(errDoRequest, loginURL, err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(errReadBody, err)
	}

	if string(respBody) != "1" {
		return fmt.Errorf("%w %s", errStatusCode, string(respBody))
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "gntsess" {
			g.session = cookie.Value
			log.Printf("Session cookie (gntsess) value: %s\n", g.session)

			return nil
		}
	}

	return errLoginNoSessionCookie
}

func (g *Geneanet) GetAccountInfos(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()

	randomBytes := make([]byte, randKLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return fmt.Errorf("key creation error: %w", err)
	}

	url := fmt.Sprintf(accountInfosURL, hex.EncodeToString(randomBytes))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(""))
	if err != nil {
		return fmt.Errorf(errNewRequest, url, err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "gntsess", Value: g.session})
	req.AddCookie(&http.Cookie{Name: "$Version", Value: "1"})

	// dumpReq, _ := httputil.DumpRequestOut(req, true)
	// log.Printf("%s\n\n", dumpReq)

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf(errDoRequest, url, err)
	}

	defer resp.Body.Close()

	// dumpResp, _ := httputil.DumpResponse(resp, true)
	// log.Printf("%s\n\n", dumpResp)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(errReadBody, err)
	}

	var out bytes.Buffer

	err = json.Indent(&out, respBody, "", "  ")
	if err != nil {
		return fmt.Errorf(errJSONMarshall, err)
	}

	log.Printf("Account infos:\n%s", out.String())

	return nil
}

func (g *Geneanet) SetLogged(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, loggedURL, strings.NewReader(""))
	if err != nil {
		return fmt.Errorf(errNewRequest, loggedURL, err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "gntsess", Value: g.session})
	req.AddCookie(&http.Cookie{Name: "$Version", Value: "1"})

	// dumpReq, _ := httputil.DumpRequestOut(req, true)
	// log.Printf("%s\n\n", dumpReq)

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf(errDoRequest, loggedURL, err)
	}

	defer resp.Body.Close()

	// dumpResp, _ := httputil.DumpResponse(resp, true)
	// log.Printf("%s\n\n", dumpResp)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(errReadBody, err)
	}

	if string(respBody) != "1" {
		return fmt.Errorf("%w %s", errStatusCode, string(respBody))
	}

	return nil
}

func (g *Geneanet) GetBase(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()

	data := url.Values{
		"st": {""},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, importURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf(errNewRequest, importURL, err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "gntsess", Value: g.session})
	req.AddCookie(&http.Cookie{Name: "$Version", Value: "1"})

	// dumpReq, _ := httputil.DumpRequestOut(req, true)
	// log.Printf("%s\n\n", dumpReq)

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf(errDoRequest, importURL, err)
	}

	defer resp.Body.Close()

	// dumpResp, _ := httputil.DumpResponse(resp, true)
	// log.Printf("%s\n\n", dumpResp)

	buf := bytes.NewBuffer([]byte{})

	size, err := io.Copy(buf, resp.Body)
	if err != nil {
		return fmt.Errorf(errIOCopy, err)
	}

	r := bytes.NewReader(buf.Bytes())

	return unzip(r, size, g.outputDir)
}
