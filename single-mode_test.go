package gittransaction_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	gittransaction "github.com/Joel-Haeberli/git-transaction"
	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
)

const testPath = "./tests"

const testTokenName = ""
const testTokenEmail = ""
const testToken = ""
const testRepoUrl = ""

const testFileContent = "this is some test content"

var testFilePath = filepath.Join(testPath, "testfile-"+uuid.NewString()+".txt")

func TestWrite(t *testing.T) {

	cloneRepoTo(testRepoUrl, "", "", testPath)

	ctx, transaction, err := gittransaction.New(gittransaction.SINGLEBRANCH, testPath, testTokenName, testToken, testTokenEmail)
	printAndFailNowOnError(err, t)

	transaction2 := gittransaction.FindTransaction(ctx)
	if transaction != transaction2 {
		fmt.Println("expected same transaction")
		t.FailNow()
	}

	err = os.WriteFile(testFilePath, []byte(testFileContent), 0666)
	printAndFailNowOnError(err, t)

	err = transaction2.Write(ctx)
	printAndFailNowOnError(err, t)

	os.RemoveAll(testPath)
}

func TestCommit(t *testing.T) {

	cloneRepoTo(testRepoUrl, "", "", testPath)

	ctx, transaction, err := gittransaction.New(gittransaction.SINGLEBRANCH, testPath, testTokenName, testToken, testTokenEmail)
	printAndFailNowOnError(err, t)

	transaction2 := gittransaction.FindTransaction(ctx)
	if transaction != transaction2 {
		fmt.Println("expected same transaction")
		t.FailNow()
	}

	err = os.WriteFile(testFilePath, []byte(testFileContent), 0666)
	printAndFailNowOnError(err, t)

	err = transaction2.Write(ctx)
	printAndFailNowOnError(err, t)

	err = transaction2.Commit(ctx)
	printAndFailNowOnError(err, t)

	os.RemoveAll(testPath)
}

func TestRollback(t *testing.T) {

	cloneRepoTo(testRepoUrl, "", "", testPath)

	ctx, transaction, err := gittransaction.New(gittransaction.SINGLEBRANCH, testPath, testTokenName, testToken, testTokenEmail)
	printAndFailNowOnError(err, t)

	transaction2 := gittransaction.FindTransaction(ctx)
	if transaction != transaction2 {
		fmt.Println("expected same transaction")
		t.FailNow()
	}

	err = os.WriteFile(testFilePath, []byte(testFileContent), 0666)
	printAndFailNowOnError(err, t)

	err = transaction2.Write(ctx)
	printAndFailNowOnError(err, t)

	err = transaction2.Rollback(ctx)
	printAndFailNowOnError(err, t)

	repo, err := git.PlainOpen(testPath)
	printAndFailNowOnError(err, t)

	wt, err := repo.Worktree()
	printAndFailNowOnError(err, t)

	status, err := wt.Status()
	printAndFailNowOnError(err, t)

	if !status.IsClean() {
		fmt.Println("expected clean worktree")
		t.FailNow()
	}

	os.RemoveAll(testPath)
}

func printAndFailNowOnError(e error, t *testing.T) {

	if e != nil {
		fmt.Println(e.Error(), e)
		t.FailNow()
	}
}

func cloneRepoTo(url string, username string, password string, location string) error {
	_, err := git.PlainClone(location, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	return err
}
