package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseConfigFile(t *testing.T) {

	configurationFile, err := os.CreateTemp("", "configuration*.json")
	require.NoError(t, err)

	configurationString := `{
	"server_address": "localhost:8080",
	"base_url": "http://localhost",
	"file_storage_path": "/path/to/file.db",
	"database_dsn": "",
	"enable_https": true
} `

	_, err = configurationFile.WriteString(configurationString)
	require.NoError(t, err)
	configurationFile.Close()
	defer os.Remove(configurationFile.Name())

	fileConfig := parseConfigFile(configurationFile.Name())

	assert.Equal(t, "localhost:8080", fileConfig.host)
	assert.Equal(t, "http://localhost", fileConfig.baseURL)
	assert.Equal(t, "/path/to/file.db", fileConfig.aliasesFile)
	assert.Equal(t, "", fileConfig.dbConnection)
	assert.Equal(t, true, fileConfig.enableHTTPS)
}

func Test_setConfiguration(t *testing.T) {

	defaultConfig := ConfigurationData{
		host:         "localhost:8080",
		baseURL:      "http://localhost",
		aliasesFile:  "/path/to/file1.db",
		usersFile:    "/path/to/file2.db",
		dbConnection: "",
		enableHTTPS:  true,
		configFile:   "",
	}

	config1 := ConfigurationData{
		aliasesFile:  "/to/file1.db",
		usersFile:    "/to/file2.db",
		dbConnection: "",
		enableHTTPS:  true,
		configFile:   "",
	}
	config2 := ConfigurationData{
		host:         "10.10.10.10:8080",
		aliasesFile:  "/file1.db",
		usersFile:    "/file2.db",
		dbConnection: "",
		enableHTTPS:  false,
		configFile:   "",
	}

	config := &Configuration{}
	config.setConfiguration(defaultConfig, config1, config2)

	assert.Equal(t, "10.10.10.10:8080", config.host)
	assert.Equal(t, "http://localhost", config.baseURL)
	assert.Equal(t, "/to/file1.db", config.aliasesFile)
	assert.Equal(t, "/to/file2.db", config.usersFile)
	assert.Equal(t, "", config.dbConnection)
	assert.Equal(t, true, config.enableHTTPS)
}
