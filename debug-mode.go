package gittransaction

// The debug transaction abstracts a dummy transaction and might
// be useful for testing or deactivating some transaction while
// not needing to remove the implementation
type DebugTransaction struct {
	Transaction
}

// On the debug transaction nothing will happen when calling Write
func (sbt *DebugTransaction) Write(ctx *Context) error {
	return nil
}

// On the debug transaction nothing will happen when calling Commit
func (sbt *DebugTransaction) Commit(ctx *Context) error {
	return nil
}

// On the debug transaction nothing will happen when calling Rollback
func (sbt *DebugTransaction) Rollback(ctx *Context) error {
	return nil
}
