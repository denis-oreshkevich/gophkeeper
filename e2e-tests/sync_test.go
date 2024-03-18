package e2e_tests

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/e2e-tests/process"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"
)

type SyncSuite struct {
	suite.Suite
	serverProcess *process.Process
}

func (suite *SyncSuite) SetupSuite() {
	ctx := context.Background()

	serverBuildCmd := exec.CommandContext(ctx, "go", "build", "-o", "../cmd/server",
		"../cmd/server")
	out, err := serverBuildCmd.CombinedOutput()
	fmt.Println(string(out))
	suite.Require().NoError(err, "ServerBuildCmd command")

	p := process.NewProcess(ctx, "../cmd/server/server")
	suite.serverProcess = p

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	err = p.Start(ctx)
	if err != nil {
		suite.T().Errorf("Невозможно запустить процесс командой %s: %s.", p, err)
		return
	}

	port := "8081"
	err = p.WaitPort(ctx, "tcp", port)
	if err != nil {
		suite.T().Errorf("Не удалось дождаться пока порт %s "+
			"станет доступен для запроса: %s", port, err)
		return
	}
}

func (suite *SyncSuite) TestSync() {
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, 20*time.Second)
	defer cancelFunc()

	clientBuildCmd := exec.CommandContext(ctx, "go", "build",
		"-o", "../cmd/client", "../cmd/client")
	out, err := clientBuildCmd.CombinedOutput()
	fmt.Println(string(out))
	suite.Require().NoError(err, "ClientBuildCmd command")

	fileCmd := exec.CommandContext(ctx, "../cmd/client/client", "file",
		"-ul=Denis", "-up=Denis", "-wd=saved", "-a=save", "-f=./test_data/bom.json", "-in=true")
	out, err = fileCmd.CombinedOutput()
	suite.Assert().NoError(err, "File command")

	fileScanner := bufio.NewScanner(bytes.NewReader(out))
	var fileID string
	for fileScanner.Scan() {
		text := fileScanner.Text()
		if strings.Contains(text, "saved file id = ") {
			split := strings.Split(text, "id = ")
			if len(split) == 2 {
				fileID = strings.Replace(split[1], "\"}", "", 1)
			}
		}
	}

	err = fileScanner.Err()
	suite.Require().NoError(err, "File scanner")
	fmt.Println(string(out))

	textCmd := exec.CommandContext(ctx, "../cmd/client/client", "text",
		"-ul=Denis", "-up=Denis", "-wd=saved", "-a=save", "-t=Denis the best", "-in=true")
	out, err = textCmd.CombinedOutput()
	suite.Assert().NoError(err, "Text command")

	textScanner := bufio.NewScanner(bytes.NewReader(out))
	var textID string
	for textScanner.Scan() {
		text := textScanner.Text()
		if strings.Contains(text, "saved text id = ") {
			split := strings.Split(text, "id = ")
			if len(split) == 2 {
				textID = strings.Replace(split[1], "\"}", "", 1)
			}
		}
	}

	err = textScanner.Err()
	suite.Require().NoError(err, "Text scanner")
	fmt.Println(string(out))

	cardCmd := exec.CommandContext(ctx, "../cmd/client/client", "card",
		"-ul=Denis", "-up=Denis", "-wd=saved", "-a=save", "-hn=\"Denis Denis\"", "-c=111",
		"-n=\"1111 1111 1111 1111\"", "-in=true")
	out, err = cardCmd.CombinedOutput()
	suite.Assert().NoError(err, "Card command")

	cardScanner := bufio.NewScanner(bytes.NewReader(out))
	var cardID string
	for cardScanner.Scan() {
		text := cardScanner.Text()
		if strings.Contains(text, "saved card id = ") {
			split := strings.Split(text, "id = ")
			if len(split) == 2 {
				cardID = strings.Replace(split[1], "\"}", "", 1)
			}
		}
	}

	err = cardScanner.Err()
	suite.Require().NoError(err, "Card scanner")
	fmt.Println(string(out))

	credCmd := exec.CommandContext(ctx, "../cmd/client/client", "cred",
		"-ul=Denis", "-up=Denis", "-wd=saved", "-a=save", "-l=Denis", "-p=Denis",
		"-in=true")
	out, err = credCmd.CombinedOutput()
	suite.Require().NoError(err, "Credentials command")

	credScanner := bufio.NewScanner(bytes.NewReader(out))
	var credID string
	for credScanner.Scan() {
		text := credScanner.Text()
		if strings.Contains(text, "saved credentials id = ") {
			split := strings.Split(text, "id = ")
			if len(split) == 2 {
				credID = strings.Replace(split[1], "\"}", "", 1)
			}
		}
	}
	err = credScanner.Err()
	suite.Require().NoError(err, "Credentials scanner")

	fmt.Println(string(out))

	suite.Run("test sync to db", func() {

		client1Cmd := exec.CommandContext(ctx, "../cmd/client/client", "sync",
			"-ul=Denis", "-up=Denis", "-wd=saved")
		out, err = client1Cmd.CombinedOutput()
		fmt.Println(string(out))
		suite.Assert().NoError(err, "Sync command 1")
	})

	suite.Run("test sync from db", func() {
		client2Cmd := exec.CommandContext(ctx, "../cmd/client/client", "sync",
			"-ul=Denis", "-up=Denis", "-wd=saved2")
		out, err = client2Cmd.CombinedOutput()
		fmt.Println(string(out))
		suite.Assert().NoError(err, "Sync command 2 with new folder")
		suite.Require().DirExists("saved2")

		checkFileCmd := exec.CommandContext(ctx, "../cmd/client/client", "file",
			"-ul=Denis", "-up=Denis", "-wd=saved2", "-a=get", "-id="+fileID)
		out, err = checkFileCmd.CombinedOutput()
		fmt.Println("FILE")
		suite.Assert().NoError(err, "check get file")
		fmt.Println(string(out))

		checkFileScan := bufio.NewScanner(bytes.NewReader(out))
		var fileRes bool
		for checkFileScan.Scan() {
			text := checkFileScan.Text()
			if strings.Contains(text, "Success") {
				fileRes = true
			}
		}
		suite.Assert().True(fileRes)

		checkTextCmd := exec.CommandContext(ctx, "../cmd/client/client", "text",
			"-ul=Denis", "-up=Denis", "-wd=saved2", "-a=get", "-id="+textID)
		out, err = checkTextCmd.CombinedOutput()
		fmt.Println("TEXT")
		suite.Assert().NoError(err, "check get text")
		fmt.Println(string(out))

		checkTextScan := bufio.NewScanner(bytes.NewReader(out))
		var textRes bool
		for checkTextScan.Scan() {
			text := checkTextScan.Text()
			if strings.Contains(text, "Denis the best") {
				textRes = true
			}
		}
		suite.Assert().True(textRes)

		checkCardCmd := exec.CommandContext(ctx, "../cmd/client/client", "card",
			"-ul=Denis", "-up=Denis", "-wd=saved2", "-a=get", "-id="+cardID)
		out, err = checkCardCmd.CombinedOutput()
		fmt.Println("CARD")
		suite.Assert().NoError(err, "check get card")
		fmt.Println(string(out))

		checkCardScan := bufio.NewScanner(bytes.NewReader(out))
		var cardRes bool
		for checkCardScan.Scan() {
			text := checkCardScan.Text()
			if strings.Contains(text, "1111 1111 1111 1111") {
				cardRes = true
			}
		}
		suite.Assert().True(cardRes)

		checkCredCmd := exec.CommandContext(ctx, "../cmd/client/client", "cred",
			"-ul=Denis", "-up=Denis", "-wd=saved2", "-a=get", "-id="+credID)
		out, err = checkCredCmd.CombinedOutput()
		fmt.Println("CRED")
		suite.Assert().NoError(err, "check get cred")
		fmt.Println(string(out))

		checkCredScan := bufio.NewScanner(bytes.NewReader(out))
		var credRes bool
		for checkCredScan.Scan() {
			text := checkCredScan.Text()
			if strings.Contains(text, "Login: Denis, password: Denis") {
				credRes = true
			}
		}
		suite.Assert().True(credRes)
	})

}

func (suite *SyncSuite) TearDownSuite() {
	exitCode, err := suite.serverProcess.Stop(syscall.SIGINT, syscall.SIGKILL)
	if err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return
		}
		suite.T().Logf("Не удалось остановить процесс с помощью сигнала ОС: %s", err)
		return
	}

	if exitCode > 0 {
		suite.T().Logf("Процесс завершился с не нулевым статусом %d", exitCode)
	}

	// получаем стандартные выводы (логи) процесса
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	out := suite.serverProcess.Stderr(ctx)
	if len(out) > 0 {
		suite.T().Logf("Получен STDERR лог процесса:\n\n%s", string(out))
	}
	out = suite.serverProcess.Stdout(ctx)
	if len(out) > 0 {
		suite.T().Logf("Получен STDOUT лог процесса:\n\n%s", string(out))
	}

	{
		pool, err := pgxpool.New(ctx, "host=localhost port=5433 user=postgres "+
			"password=postgres dbname=keeper sslmode=disable")
		suite.Assert().NoError(err)
		defer pool.Close()

		querySchema := "drop schema if exists keeper cascade"

		_, err = pool.Exec(ctx, querySchema)
		suite.Assert().NoError(err, "drop schema")

		queryMigr := "drop table if exists public.goose_db_version"

		_, err = pool.Exec(ctx, queryMigr)
		suite.Assert().NoError(err, "drop table")
	}

	{
		err = os.RemoveAll("saved")
		suite.Assert().NoError(err, "remove saved")

		err = os.RemoveAll("saved2")
		suite.Assert().NoError(err, "remove saved2")

		err = os.RemoveAll("certs")
		suite.Assert().NoError(err, "remove certs")

		err = os.Remove("../cmd/client/client")
		suite.Assert().NoError(err, "remove client")

		err = os.RemoveAll("../cmd/server/server")
		suite.Assert().NoError(err, "remove server")
	}
}

func TestSync(t *testing.T) {
	suite.Run(t, new(SyncSuite))
}
