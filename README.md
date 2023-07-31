# git-transaction

The purpose of this library is to use any git repository of your choice as 
single source of truth, enabling transaction support.

# process

A git transaction follows this workflow (we assume that the repository was already pulled
to some specified path):

1. Create new transaction using the factory method of the package
2. Write to the transaction as much as you want. Each write will trigger a commit.
3. Commit the transaction if the transaction has ended.

If something fails, a rollback is triggered and the branch is reset to its initial state. 
(the state before creating the new transaction)

# concept of modes

Currently two modes are planned. One mode (**Singlebranch Transaction**) is thougt to
be the simplest version of a transaction, which just writes to one single branch. This
mode is sufficient, as long as only one single instance is working on the repository.

The second mode (**Multibranch Transaction**) aims to enable multiple instance working
on the same 

# roadmap

To give you an idea on what is planned for this tool, here is a small roadmap
ordered by time of completion ascending. If you have ideas for features you would 
like to see, don't hesitate to open an issue or even a merge request. The order
of implementation can change or a goal can be omitted completely.

1. Singlebranch Transaction
2. Automatically detect writes (which calls write under the hood)
3. Support repositories which require credentials and allow signing transactions (configure key and sign commits)
4. Support ssh (not ony http)
5. Multibranch Transaction
6. Automatic conflict resolution (maybe some rule based resolution strategy which can be defined?)

For goal three and four the underlying [go-git](https://github.com/go-git/go-git) implementation needs 
to be extended (for example the ability to merge branches must be implemented). A current version 
of supported devices can be found in the [go-git](https://github.com/go-git/go-git) 
[compatibility documentation](https://github.com/go-git/go-git/blob/v5.8.0/COMPATIBILITY.md)

# when you should NOT use this library

You should **definitely not use** this library if:

- you need bleeding edge performance (networking is kind of a bottleneck)
- you can / are willing to use a database
- you need high concurrency guarantees. (since we are using git as underlying technology 
this is theoretically not a problem but maybe some manual action is required from time to time, 
which is not what we're looking for I think)

# contribute

You want to contribute? I'm grateful for every contribution wether it is an issue, bug or even
a merge request. Please take a look at [CONTRIBUTION.md](./CONTRIBUTION.md) for further information.