package libv2ray

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const APP_KEY = "573d23e9ff036761bfb8d179717453173d8191834b73cb9388d12a00228cb62e"
const SERVER_KEY = "afa2d1267fa1609aa58e268312ca2940d7793c5688eb6a2fe9e00b5f9f68e014"
const APP_STORE_FINGERPRINT = "C0:6F:83:53:83:12:07:D2:D6:76:C7:4C:E6:89:57:BB:B4:18:C6:23:EA:91:9C:9F:AE:D5:B4:F9:C1:89:22:9F"

type KeyService interface {
	GetAppKey() string
	GetServerKey() string
}

func GetAppKey() string {
	return APP_KEY
}

func GetServerKey() string {
	return SERVER_KEY
}

func GetFingerPrint() string {
	return APP_STORE_FINGERPRINT
}

func RevertConfigBack(config string) (string, error) {
	// Check if the config string starts with the expected scheme
	if !strings.HasPrefix(config, "vless://") {
		return "", errors.New("invalid config format")
	}

	// Split the config string to extract UUID, IP, and port
	configParts := strings.Split(config, "@")
	if len(configParts) != 2 {
		return "", errors.New("invalid config format")
	}

	// Extract and modify UUID
	uuid := configParts[0][8:] // skip "vless://"
	uuidParts := strings.Split(uuid, "-")
	if len(uuidParts) != 5 {
		return "", errors.New("invalid UUID format")
	}
	uuidFourthPart, err := strconv.ParseUint(uuidParts[3], 16, 64)
	if err != nil {
		return "", errors.New("invalid UUID part")
	}
	uuidFourthPart--
	uuidParts[3] = fmt.Sprintf("%04x", uuidFourthPart)
	originalUUID := strings.Join(uuidParts, "-")

	// Split to extract IP:port and query
	ipPortAndQuery := configParts[1]
	ipPortParts := strings.Split(ipPortAndQuery, ":")
	if len(ipPortParts) < 2 {
		return "", errors.New("invalid IP:port format")
	}

	// Modify IP
	ip := ipPortParts[0]
	ipParts := strings.Split(ip, ".")
	if len(ipParts) != 4 {
		return "", errors.New("invalid IP address")
	}
	lastByte, err := strconv.Atoi(ipParts[3])
	if err != nil {
		return "", errors.New("invalid IP address part")
	}
	lastByte--
	ipParts[3] = strconv.Itoa(lastByte)
	originalIP := strings.Join(ipParts, ".")

	// Extract port and query
	portAndQuery := strings.SplitN(ipPortParts[1], "?", 2)
	port := portAndQuery[0]
	query := ""
	if len(portAndQuery) > 1 {
		query = portAndQuery[1]
	}

	// Modify port
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return "", errors.New("invalid port number")
	}
	portInt -= 5
	originalPort := strconv.Itoa(portInt)

	// Construct the original config string
	originalconfig := fmt.Sprintf("vless://%s@%s:%s?%s", originalUUID, originalIP, originalPort, query)
	return originalconfig, nil
}

func RevertComplexProxy(proxy string) (string, error) {
	// Revert UUID
	uuidRegex := regexp.MustCompile(`([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-)([0-9a-fA-F]{4})(-[0-9a-fA-F]{12})`)
	proxy = uuidRegex.ReplaceAllStringFunc(proxy, func(match string) string {
		parts := uuidRegex.FindStringSubmatch(match)
		uuidFourthPart, _ := strconv.ParseUint(parts[2], 16, 64)
		uuidFourthPart--
		return parts[1] + fmt.Sprintf("%04x", uuidFourthPart) + parts[3]
	})

	// Revert IP address
	ipRegex := regexp.MustCompile(`(\d+\.\d+\.\d+\.)(\d+)`)
	proxy = ipRegex.ReplaceAllStringFunc(proxy, func(match string) string {
		parts := ipRegex.FindStringSubmatch(match)
		lastByte, _ := strconv.Atoi(parts[2])
		lastByte--
		return parts[1] + strconv.Itoa(lastByte)
	})

	// Revert port
	portRegex := regexp.MustCompile(`'port':(\d+)`)
	proxy = portRegex.ReplaceAllStringFunc(proxy, func(match string) string {
		parts := portRegex.FindStringSubmatch(match)
		port, _ := strconv.Atoi(parts[1])
		port -= 5
		return fmt.Sprintf(`'port':%d`, port)
	})

	return proxy, nil
}
