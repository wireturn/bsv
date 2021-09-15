let bsv = require('bsv')

function specToBrfcId(spec) {
    let hash = bsv.crypto.Hash.sha256sha256(Buffer.from(
        spec.title.trim() +
        (spec.author || '').trim() +
        (spec.version || '').trim()
    ));

    let bitcoinDisplayHash = hash
        .reverse()
        .toString('hex');

    return bitcoinDisplayHash.substring(0, 12);

}

let spec = {
    title: "BRFC Specifications",
    author: "andy (nChain)",
    version: "1"
}

console.log(specToBrfcId(spec));