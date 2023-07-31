package gittransaction

import (
	"errors"
	"fmt"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// The singlebranch transaction abstracts a git-transaction inside one branch
// and represents the most simple version of a git transaction.
type SinglebranchTransaction struct {
	Transaction

	idGenStrategy IdGenerationStrategy
}

// Write will write all current changes to the transaction by creating a commit.
func (sbt *SinglebranchTransaction) Write(ctx *TransactionContext) error {

	if ctx.credentials == nil {
		return errors.New("need to define credentials in order to write to transaction")
	}

	wt, err := update(ctx.path)
	if err != nil {
		return err
	}

	addOpts := new(git.AddOptions)
	addOpts.All = true

	err = wt.AddWithOptions(addOpts)
	if err != nil {
		return err
	}

	opts := new(git.CommitOptions)
	opts.Author = &ctx.credentials.Signature
	_, err = wt.Commit("transaction: "+ctx.Id, opts)
	if err != nil {
		return err
	}

	return nil
}

// Commit will write commit the transaction by pushing to the defined repository.
func (sbt *SinglebranchTransaction) Commit(ctx *TransactionContext) error {

	err := push(ctx)
	if err != nil {
		return err
	}

	removeTransaction(ctx.Id)

	return nil
}

// Rollback will reset to the state where the branch transaction started.
// Rollback is called if an error happens during writing or commiting a transaction
func (sbt *SinglebranchTransaction) Rollback(ctx *TransactionContext) error {

	err := reset(ctx.path, ctx.headBeforeTransaction)
	if err != nil {
		return err
	}

	return nil
}

func push(ctx *TransactionContext) error {

	repo, err := repository(ctx.path)
	if err != nil {
		return err
	}

	err = repo.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: ctx.credentials.Name,
			Password: ctx.credentials.accessToken,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func update(location string) (*git.Worktree, error) {

	worktree, err := worktree(location)
	if err != nil {
		return nil, err
	}

	err = worktree.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err.Error() != "already up-to-date" {
		return nil, err
	}
	return worktree, nil
}

func reset(location string, to plumbing.Hash) error {

	wt, err := worktree(location)
	if err != nil {
		return err
	}

	resetOpts := new(git.ResetOptions)
	resetOpts.Mode = git.HardReset
	resetOpts.Commit = to

	err = wt.Reset(resetOpts)
	if err != nil {
		return err
	}

	wt, err = update(location)
	if err != nil {
		return err
	}

	return nil
}

func worktree(location string) (*git.Worktree, error) {

	repo, err := repository(location)
	if err != nil {
		return nil, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Println("reading worktree failed git repo failed:", err.Error())
		return nil, err
	}

	return worktree, nil
}

func repository(location string) (*git.Repository, error) {

	repo, err := git.PlainOpen(location)
	if err != nil {
		fmt.Println("opening git repo at", location, "failed:", err.Error())
		return nil, err
	}

	return repo, nil
}
