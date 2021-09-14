## paymail brfc

List all specs, search by keyword, or generate a new BRFC ID

### Synopsis

```
___.           _____       
\_ |__________/ ____\____  
 | __ \_  __ \   __\/ ___\ 
 | \_\ \  | \/|  | \  \___ 
 |___  /__|   |__|  \___  >
     \/                 \/
```

Use the [list] argument to show all known BRFC protocols.

Use the [generate] argument with required flags to generate a new BRFC ID.

Use the [search] argument to show any matching BRFCs by either ID, Title or Author.

BRFC (Bitcoin SV Request-For-Comments) Specifications describe functionality across the ecosystem. 
"bsvalias" protocols and paymail implementations are described across a series of BRFC documents.

Whilst this is not the authoritative definition of the BRFC process, a summary is included here 
as the BRFC process is the nominated mechanism through which extensions to the paymail system 
are defined and discovered during Service Discovery.

Read more at: http://bsvalias.org/01-brfc-specifications.html

```
paymail brfc [flags]
```

### Examples

```
paymail brfc list
paymail brfc search nChain
paymail brfc generate --title "BRFC Specifications" --author "andy (nChain)" --version 1
```

### Options

```
      --author string     Author(s) new BRFC specification
  -h, --help              help for brfc
      --skip-validation   Skip validating the existing BRFC IDs
      --title string      Title of the new BRFC specification
      --version string    Version of the new BRFC specification
```

### Options inherited from parent commands

```
      --bsvalias string   The bsvalias version (default "1.0")
      --config string     Custom config file (default is $HOME/paymail/config.yaml)
      --docs              Generate docs from all commands (./docs/commands)
      --flush-cache       Flushes ALL cache, empties local database
      --no-cache          Turn off caching for this specific command
  -t, --skip-tracing      Turn off request tracing information
```

### SEE ALSO

* [paymail](paymail.md)	 - Inspect, validate domains or resolve paymail addresses

