# Spy Node

The package provides a simple node, that provides raw Bitcoin network data.
This node does not verify signatures or consensus rules. It simply assumes the chain given by
  the trusted external full node is valid. The other nodes are only for more tx propagation
  visibility for double spend detection.

While there is a binary, it is only used for testing.

This package is meant to be used as a building block (data source), rather than a node.

BroadcastTx will send a tx message to all connected peers, trusted and untrusted. Use this to send
  txs that are created and need to be put on chain.

Specify a "start" block hash in the config.
All transactions and blocks in the blockchain (unfiltered) from that point will be fed through the
  tx and block listeners.
This will send a full tx message to the tx listener when a tx is announced or included in a block.
Whenever it is seen first.
This will send a "new block" message to the block listener when it sees a new block followed by
sending every tx hash included in the block to the block listener.
When it comes online it will "sync" with the external full node and process all blocks not yet
  processed, then process any in the external nodes mempool.
From then on it processes tx and blocks as they are announced/propagated.

Life of a tx:
	The `HandleTx` callback function notifies you of a new tx.
	The `HandleTxUpdate` callback provides information about the state of a transaction as it changes.
	The `HandleHeaders` callback provides headers when subscribed.
	The `HandleInSync` callback notifies when data provided is update to date with the network.


Double spends:
We currently have two ways of detecting double spends.
	If a conflicting tx is mined, then spynode detects this by checking for conflicts in its mempool.
	Spynode listens to several nodes for tx announcements and if one of the nodes sees a
	  conflicting tx first it will announce to us and we will know of a double spend attempt.
It is still possible for someone to send double spends directly to miners and if the miners don't
  propagate the tx (because of fee or other reasons), then we will not know about it until confirm.

Untrusted Nodes:
Peer connection that just listens for transactions.
Used to monitor for double spend attempts.
Life cycle:
	When we are synced with the trusted peer we start connecting to untrusted nodes.
	Initially we request some headers to make sure we are on the same chain.
	Then we just listen for txs and add them to the mempool.
	As they are added we check for other txs with conflicting inputs (double spends).
	If conflicts are detected, then a ListenerMsgTxUnsafe message is sent for all txs involved.
Malicious:
	Technically if we connect to a malicious untrusted node, they could send us invalid txs. They
	couldn't fake confirms, but we might see pending tx that aren't valid or double spends that
	aren't valid and we wouldn't actually care about them. Spynode won't mark a tx as safe until
	the trusted node has verified it, so it will just sit in the mempool.


### Makefile

If you're interested in working with the source code, this is your best option.

```
mkdir -p $GOPATH/src/github.com/tokenized
cd $GOPATH/src/github.com/tokenized
git clone git@github.com:tokenized/spynode.git
cd spynode
make deps
make build
```

On Windows use build-win and run-win.


## Configuration

Configuration is supplied via environment variables. See the
[example](conf/dev.env.example) file that would export environment variables
on Linux and OS X.

Make a copy the example file, and edit the values to suit.

The file can be placed anywhere you prefer.


## Running

This example shows the config file containing the environment variables
residing at `./tmp/dev.env`

```
source ./tmp/dev.env && go get cmd/spynode/main.go
```
