package config

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testEnvFile string = ".testenv"
)

func checkErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("could not execute test preparation. Error: %s", err))
	}
}

func writeTestEnv(fileName string) {
	f, err := os.Create(fileName)
	checkErr(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	_, err = w.WriteString("GIN_MODE=\"debug\"\n")
	checkErr(err)
	_, err = w.WriteString("DB_USER=\"user\"\n")
	checkErr(err)
	_, err = w.WriteString("DB_PASSWD=\"passwd\"\n")
	checkErr(err)
	_, err = w.WriteString("DB_ADDR=\"localhost\"\n")
	checkErr(err)
	_, err = w.WriteString("DB_PORT=\"3306\"\n")
	checkErr(err)
	_, err = w.WriteString("DB_NAME=\"name\"\n")
	checkErr(err)
	_, err = w.WriteString("SERVER_ADDR=\"127.0.0.1\"\n")
	checkErr(err)
	_, err = w.WriteString("SERVER_PORT=\"9999\"\n")
	checkErr(err)
	w.Flush()
}

func deleteEnvFile(fileName string) {
	err := os.Remove(fileName)
	checkErr(err)
}

func unsetEnvVars() {
	os.Unsetenv("GIN_MODE")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWD")
	os.Unsetenv("DB_ADDR")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("SERVER_ADDR")
	os.Unsetenv("SERVER_PORT")
}

func Test_loadConfig_NoEnvFile_Returns_Error(t *testing.T) {
	err := loadConfig("file_does_not_exist.txt")
	assert.NotNil(t, err)
	assert.EqualValues(t, "open file_does_not_exist.txt: The system cannot find the file specified.", err.Error())
}

func Test_loadConfig_WithEnvFile_Returns_NoError(t *testing.T) {
	writeTestEnv(testEnvFile)
	defer deleteEnvFile(testEnvFile)
	err := loadConfig(testEnvFile)
	defer unsetEnvVars()
	assert.Nil(t, err)
	assert.EqualValues(t, "debug", os.Getenv("GIN_MODE"))

}

func Test_configDb_NoEnvVars_Returns_Error(t *testing.T) {
	err := configDb()
	assert.NotNil(t, err)
	assert.EqualValues(t, "DB environment not defined", err.Error())
}

func Test_configDb_WithEnvVars_Returns_NoError(t *testing.T) {
	writeTestEnv(testEnvFile)
	defer deleteEnvFile(testEnvFile)
	loadConfig(testEnvFile)
	defer unsetEnvVars()
	err := configDb()
	assert.Nil(t, err)
	assert.EqualValues(t, "user", os.Getenv("DB_USER"))
	assert.EqualValues(t, "passwd", os.Getenv("DB_PASSWD"))
	assert.EqualValues(t, "localhost", os.Getenv("DB_ADDR"))
	assert.EqualValues(t, "3306", os.Getenv("DB_PORT"))
	assert.EqualValues(t, "name", os.Getenv("DB_NAME"))
}

func Test_configGin_NoEnvVars_SetsReleaseMode(t *testing.T) {
	assert.EqualValues(t, "", GinMode)
	configGin()
	assert.EqualValues(t, "release", GinMode)
}

func Test_configGin_WrongEnvVar_SetsReleaseMode(t *testing.T) {
	os.Setenv("GIN_MODE", "bogus")
	configGin()
	assert.EqualValues(t, "release", GinMode)
	os.Unsetenv("GIN_MODE")
}

func Test_configGin_WithEnvVar_SetsMode(t *testing.T) {
	os.Setenv("GIN_MODE", "debug")
	configGin()
	assert.EqualValues(t, "debug", GinMode)
	os.Unsetenv("GIN_MODE")
}

func Test_configServer_NoEnvVars_SetsDefaults(t *testing.T) {
	configServer()
	assert.EqualValues(t, "", ServerAddr)
	assert.EqualValues(t, "8080", ServerPort)
}

func Test_configServer_WithEnvVars_SetsValues(t *testing.T) {
	writeTestEnv(testEnvFile)
	defer deleteEnvFile(testEnvFile)
	loadConfig(testEnvFile)
	defer unsetEnvVars()
	configServer()
	assert.EqualValues(t, "127.0.0.1", ServerAddr)
	assert.EqualValues(t, "9999", ServerPort)
}

func Test_InitConfig_NoEnvVars_Returns_Error(t *testing.T) {
	err := InitConfig("")
	assert.NotNil(t, err)
	assert.EqualValues(t, "DB environment not defined", err.Error())
}

func Test_InitConfig_WithEnvVars_Returns_NoError(t *testing.T) {
	writeTestEnv(testEnvFile)
	defer deleteEnvFile(testEnvFile)
	err := InitConfig(testEnvFile)
	assert.Nil(t, err)
	unsetEnvVars()
}
