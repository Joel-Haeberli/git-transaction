package gittransaction

import (
	"errors"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	object "github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/uuid"
)

// the mode defines how the transaction is executed
//   - SINGLEBRANCH: the transaction will be done in the currently active branch
//
// currently only the SINGLEBRANCH option s supported.
type MODE int

const (
	SINGLEBRANCH MODE = 1 << iota
	// MULTIBRANCH
)

// currently active transactions. In SINGLEBRANCH mode, only one transaction
// per repository can be active to prevent conflicts.
var ongoingTransactions = make(map[string]Transaction, 0)

// the default id generation strategy
var idGenerationStrategy IdGenerationStrategy = new(PathIdGenerationStrategy)

// The idGenerationStrategy interface defines the abstraction for creating an id.
type IdGenerationStrategy interface {
	GenerateId(string) string
}

// The Transaction interface defines the abstraction of a transaction
type Transaction interface {
	Write(*Context) error
	Commit(*Context) error
	Rollback(*Context) error
}

// The transaction context holds information about a transaction
type Context struct {
	Id                    string
	path                  string
	headBeforeTransaction plumbing.Hash
	credentials           *credentials
}

// A token represents the properties of an acccess token which are needed
// to allow its usage for
type Token struct {
	Username string
	Token    string
	Email    string
}

// credentials wrap the the token in git library manner
type credentials struct {
	object.Signature

	accessToken string
}

// the UUID as id generation strategy
type UUIDIdGenerationStrategy struct {
	IdGenerationStrategy
}

// generates an UUID as transaction id. This id will be written
// in each write allowing to trace what was done inside one transaction
func (strategy *UUIDIdGenerationStrategy) GenerateId(seed string) string {

	uuid, err := uuid.NewRandomFromReader(strings.NewReader(seed))

	if err != nil {
		panic("unable to generate uuid")
	}

	return uuid.String()
}

// the path id generation strategy
type PathIdGenerationStrategy struct {
	IdGenerationStrategy
}

// the id simply represents the path of the underlying repository
func (strategy *PathIdGenerationStrategy) GenerateId(path string) string {

	return path
}

// set the id generation strategy which shall be used to generate the transaction ids.
func SetIdGenerationStrategy(strategy IdGenerationStrategy) {

	idGenerationStrategy = strategy
}

// setup credentials for the given transaction
func SetupCredentials(ctx *Context, token *Token) {

	creds := new(credentials)
	creds.Email = token.Email
	creds.Name = token.Username
	creds.accessToken = token.Token

	ctx.credentials = creds
}

// creates a new transaction in given mode and path (must be git repo)
// returns a transaction context and the transaction or an error
func New(m MODE, path string, token *Token) (*Context, Transaction, error) {

	for id := range ongoingTransactions {
		if id == idGenerationStrategy.GenerateId(path) {
			return nil, nil, errors.New("only one active transaction per repository allowed")
		}
	}

	wt, err := repository(path)
	if err != nil {
		return nil, nil, err
	}

	ref, err := wt.Head()
	if err != nil {
		return nil, nil, err
	}

	var transaction Transaction = new(SinglebranchTransaction)
	ctx := new(Context)
	ctx.path = path
	ctx.Id = idGenerationStrategy.GenerateId(path)
	ctx.headBeforeTransaction = ref.Hash()

	SetupCredentials(ctx, token)

	addTransaction(ctx.Id, transaction)

	return ctx, transaction, nil
}

// returns a transaction if matching id is found, otherwise nil
func FindTransaction(ctx *Context) Transaction {

	for id, trs := range ongoingTransactions {
		if id == ctx.Id {
			return trs
		}
	}

	return nil
}

func addTransaction(id string, trs Transaction) {

	ongoingTransactions[id] = trs
}

func removeTransaction(id string) {

	delete(ongoingTransactions, id)
}
