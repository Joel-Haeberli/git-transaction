package gittransaction_test

import (
	"testing"

	gittransaction "github.com/Joel-Haeberli/git-transaction"
)

func TestDebugWrite(t *testing.T) {

	ctx, transaction, err := gittransaction.New(gittransaction.DEBUG, testPath, token)
	printAndFailNowOnError(err, t)

	err = transaction.Write(ctx)
	printAndFailNowOnError(err, t)
}

func TestDebugCommit(t *testing.T) {

	cloneRepoTo(testRepoUrl, "", "", testPath)

	ctx, transaction, err := gittransaction.New(gittransaction.DEBUG, testPath, token)
	printAndFailNowOnError(err, t)

	err = transaction.Write(ctx)
	printAndFailNowOnError(err, t)

	err = transaction.Commit(ctx)
	printAndFailNowOnError(err, t)
}

func TestDebugRollback(t *testing.T) {

	ctx, transaction, err := gittransaction.New(gittransaction.DEBUG, testPath, token)
	printAndFailNowOnError(err, t)

	err = transaction.Write(ctx)
	printAndFailNowOnError(err, t)

	err = transaction.Commit(ctx)
	printAndFailNowOnError(err, t)

	err = transaction.Rollback(ctx)
	printAndFailNowOnError(err, t)
}
