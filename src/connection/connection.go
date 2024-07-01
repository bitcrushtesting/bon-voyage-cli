// Copyright 2024 Bitcrush Testing

package connection

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"bon-voyage-cli/models"

	"github.com/gorilla/websocket"
)

var (
	host  string
	port  string
	token string
)

func Init(h string, p string) {
	host = h
	port = p
}

func serverUrl() string {
	return fmt.Sprintf("http://%s:%s", host, port)
}

func DeviceList() ([]models.DeviceBase, error) {

	var devices []models.DeviceBase
	u := serverUrl() + models.ClientsApiPath + "/list"
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return devices, err
	}
	req.Header.Add("Authorization", "Bearer "+Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return devices, fmt.Errorf("failed to get device list: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return devices, fmt.Errorf("failed to get device list: %s", string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return devices, fmt.Errorf("failed to read response body: %v", err)
	}

	err = json.Unmarshal(body, &devices)
	if err != nil {
		return devices, fmt.Errorf("failed to decode JSON response: %v", err)
	}
	return devices, nil
}

func DeviceSocket(deviceUuid string) error {

	u := url.URL{
		Scheme: "ws",
		Host:   host + ":" + port,
		Path:   models.ClientsApiPath + "/socket/" + deviceUuid,
	}

	h := http.Header{}
	h.Add("Authorization", "Bearer "+Token)

	fmt.Println("Connecting to:", deviceUuid)

	d := websocket.DefaultDialer
	socket, resp, err := d.Dial(u.String(), h)
	if resp.StatusCode != http.StatusSwitchingProtocols {
		return fmt.Errorf("dial: %s", resp.Status)
	}
	if err != nil {
		return fmt.Errorf("dial: %v", err)
	}
	defer socket.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := socket.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					return
				}
				fmt.Println("read:", err)
				return
			}
			fmt.Printf("Received: %s\n", message)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			err := socket.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				fmt.Println("write:", err)
				return
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading from stdin:", err)
		}
	}()

	select {
	case <-done:
		fmt.Println("WebSocket connection closed")
	case <-interrupt:
		fmt.Println("\nInterrupt signal received, closing connection...")
		err := socket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			fmt.Println("write close:", err)
			return err
		}
		select {
		case <-done:
		case <-time.After(time.Second):
		}
	}
	return nil
}

func SessionList() ([]models.SessionListReply, error) {

	var sessions []models.SessionListReply
	u := serverUrl() + models.ClientsApiPath + "/session/list"
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return sessions, err
	}
	req.Header.Add("Authorization", "Bearer "+Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return sessions, fmt.Errorf("failed to get session list: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return sessions, fmt.Errorf("failed to get session list: %s", string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return sessions, fmt.Errorf("failed to read response body: %v", err)
	}

	err = json.Unmarshal(body, &sessions)
	if err != nil {
		return sessions, fmt.Errorf("failed to decode JSON response: %v", err)
	}
	return sessions, nil
}

func SessionCreate(deviceID string) (string, error) {

	u := serverUrl() + models.ClientsApiPath + "/session"

	payload := models.SessionCreatePayload{
		Command:  "create",
		DeviceId: deviceID,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to create session: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	var session models.SessionCreateReply
	err = json.Unmarshal(body, &session)
	if err != nil {
		return "", fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return session.SessionId, nil
}

func SessionRead(sessionID string, filter []string) ([]string, error) {

	var sessions []string
	u := serverUrl() + models.ClientsApiPath + "/session/" + sessionID
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return sessions, err
	}
	req.URL.Query()
	req.Header.Add("Authorization", "Bearer "+Token)

	return nil, nil
}

func SessionUpdate(sessionID, command string) error {

	u := serverUrl() + models.ClientsApiPath + "/session/" + sessionID

	jsonBody := fmt.Sprintf(`{"command": "%s"}`, command)
	b := bytes.NewReader([]byte(jsonBody))
	req, err := http.NewRequest("POST", u, b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update session: %s", string(body))
	}

	return nil
}

func SessionDelete(sessionID string) error {

	u := serverUrl() + models.ClientsApiPath + "/session/" + sessionID
	req, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete session: %s", string(body))
	}

	return nil
}

func Information() (models.Information, error) {

	var info models.Information
	resp, err := http.Get(serverUrl() + "/info")
	if err != nil {
		return info, fmt.Errorf("failed to make info request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return info, fmt.Errorf("received non-200 response status: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return info, fmt.Errorf("failed to read response body: %v", err)
	}

	err = json.Unmarshal(body, &info)
	if err != nil {
		return info, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return info, nil
}

func Login(username, password string) error {

	if username == "" {
		return fmt.Errorf("username is required")
	}

	if password == "" {
		return fmt.Errorf("password is required")
	}

	loginData := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return fmt.Errorf("failed to marshal login data: %v", err)
	}
	url := serverUrl() + "/login"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to login: %s", string(body))
	}

	Token, err = extractTokenFromCookies(resp.Cookies())
	if err != nil {
		return fmt.Errorf("failed to extract token from cookies: %v", err)
	}
	return SaveToken()
}

func Register(username, password string) error {

	if username == "" {
		return fmt.Errorf("username is required")
	}
	if password == "" {
		return fmt.Errorf("password is required")
	}

	registerData := models.UserRegisterCredentials{
		Username: username,
		Password: password,
	}
	jsonData, err := json.Marshal(registerData)
	if err != nil {
		return fmt.Errorf("failed to marshal register data: %v", err)
	}
	url := serverUrl() + "/register"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to register: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to register: %s", string(body))
	}
	return nil
}

func ChangeUsername(newUsername string) error {

	if newUsername == "" {
		return fmt.Errorf("new username is required")
	}

	data := url.Values{}
	data.Set("newUsername", newUsername)
	url := serverUrl() + "/change_username"
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to change username: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to change username: %s", string(body))
	}
	return nil
}

func ChangePassword(newPassword string) error {

	if newPassword == "" {
		return fmt.Errorf("new password is required")
	}

	data := url.Values{}
	data.Set("newPassword", newPassword)
	u := serverUrl() + "/change_password"
	req, err := http.NewRequest("POST", u, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to change password: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to change password: %s", string(body))
	}

	return nil
}

func Logout() error {
	return nil
}

func DeleteAccount(username, password string) error {
	return nil
}

func extractTokenFromCookies(cookies []*http.Cookie) (string, error) {

	for _, cookie := range cookies {
		if cookie.Name == "token" {
			return cookie.Value, nil
		}
	}
	return "", fmt.Errorf("token not found in cookies")
}
