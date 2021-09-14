# merkle-proof-standard-example

## Requirements

For development, you will only need Node.js (minimum 10.12.0) and a node global package, NPM, installed in your environment.

### Node

- #### Node installation on Windows

  Just go on [official Node.js website](https://nodejs.org/) and download the installer.
Also, be sure to have `git` available in your PATH, `npm` might need it (You can find git [here](https://git-scm.com/)).

- #### Node installation on Ubuntu

  You can install nodejs and npm easily with apt install, just run the following commands.

    ```sh
    $ sudo apt install nodejs
    $ sudo apt install npm
    ```

- #### Other Operating Systems

  You can find more information about the installation on the [official Node.js website](https://nodejs.org/) and the [official NPM website](https://npmjs.org/).

## Install

```sh
$ git clone https://github.com:bitcoin-sv-specs/merkle-proof-standard-example.git
$ cd merkle-proof-standard-example
$ npm install
```

## Run

```sh
$ npm test
```

## Code Style

The JavaScript code in this repository is written using [JavaScript Standard Style](https://standardjs.com/).

## Example Merkle Proof

### JSON format (default)

```json
  {
    "index": 12,
    "txOrId": "ffeff11c25cde7c06d407490d81ef4d0db64aad6ab3d14393530701561a465ef",
    "target": "75edb0a69eb195cdd81e310553aa4d25e18450e08f168532a2c2e9cf447bf169",
    "nodes": [
      "b9ef07a62553ef8b0898a79c291b92c60f7932260888bde0dab2dd2610d8668e",
      "0fc1c12fb1b57b38140442927fbadb3d1e5a5039a5d6db355ea25486374f104d",
      "60b0e75dd5b8d48f2d069229f20399e07766dd651ceeed55ee3c040aa2812547",
      "c0d8dbda46366c2050b430a05508a3d96dc0ed55aea685bb3d9a993f8b97cc6f",
      "391e62b3419d8a943f7dbc7bddc90e30ec724c033000dc0c8872253c27b03a42"
    ]
  }
```

### Binary format (hex encoding)

```hex
000cef65a4611570303539143dabd6aa64dbd0f41ed89074406dc0e7cd251cf1efff69f17b44cfe9c2a23285168fe05084e1254daa5305311ed8cd95b19ea6b0ed7505008e66d81026ddb2dae0bd88082632790fc6921b299ca798088bef5325a607efb9004d104f378654a25e35dbd6a539505a1e3ddbba7f92420414387bb5b12fc1c10f00472581a20a043cee55edee1c65dd6677e09903f22992062d8fd4b8d55de7b060006fcc978b3f999a3dbb85a6ae55edc06dd9a30855a030b450206c3646dadbd8c000423ab0273c2572880cdc0030034c72ec300ec9dd7bbc7d3f948a9d41b3621e39
```

> By default `txOrId` will contain a txId and `target` will contain a block hash
