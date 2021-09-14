package sqlite

const (
	sqlTransactionCreate = `
		INSERT INTO transactions(txid, paymentID, txhex)
		VALUES(:txid, :paymentID, :txhex)
	`

	sqlTransactionByID = `
	SELECT txid, paymentid, txhex, createdat
	FROM transactions
	WHERE txid = :txid
	`

	sqlTxosByTxID = `
	SELECT outpoint, txid, vout, keyname, derivationpath, lockingscript, satoshis, 
				spentat, spendingtxid, createdat, modifiedAt 
	FROM txos
	WHERE txid = :txid
	`
)
