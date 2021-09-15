# Host Discovery

Host discovery is the process through which a domain owner optionally specifies which web host to interrogate during capability discovery. The host discovery process involves the creation and subsequent query of `SRV` DNS records.

The use of an `SRV` record was chosen for the following reasons:

* Domain owners may choose to use a third-party paymail service provider. Delegating authority to this provider is a one-time activity (the creation of a DNS `SRV` record).
* `TXT` records were considered, however should the paymail service be reconfigured (for example, the root URI change from `.../api/v1/...` to `.../api/v2/...`), the domain owner would have to coordinate with the service provider to ensure the `TXT` record (containing a full endpoint URI) was updated. With `SRV` records (plus the [Capability Discovery](02-02-capability-discovery.md) protocol), the Host Discovery phase is _set-and-forget_.
* As an optional step, if the canonical Capability Discovery host is the same as the domain found in the paymail alias, the DNS record can be omitted entirely.
* DNS records (including `SRV`) feature a TTL, which clients can use for caching responses. All common DNS clients implement this caching _out-of-the-box_, meaning implementers do not have to roll this themselves.

## Setup

A domain owner may create an `SRV` record with the following parameters:

| Parameter | Value                       |
|-----------|-----------------------------|
| Service   | `_bsvalias`                 |
| Proto     | `_tcp`                      |
| Name      | `<domain>.<tld>.`           |
| TTL       | `3600` (see notes)          |
| Class     | `IN`                        |
| Priority  | `10`                        |
| Weight    | `10`                        |
| Port      | `443`                       |
| Target    | `<endpoint-discovery-host>` |

The `TTL` parameter should be set very low for test configurations (a few seconds), whereas for production deployments this should be set higher, to allow caching to work. A value of `3600` is suggested for production deployments.

Although the DNS system allows for multiple records with a variety of priorities and weights, which allows for some level of traffic management, resilience, and load-balancing via DNS records, it is recommended by this specification that these considerations be handled by more modern infrastructure and only a single `SRV` record be created.

See [https://en.wikipedia.org/wiki/SRV_record](https://en.wikipedia.org/wiki/SRV_record) for more information on `SRV` DNS records.

`SRV` records **must** be served with a valid [DNSSEC](https://en.wikipedia.org/wiki/Domain_Name_System_Security_Extensions) signature chain with the exception that if the `SRV` record points to `<domain>.<tld>` or `www.<domain>.<tld>`, then the `SRV` records **should** be served with a valid DNSSEC signature chain.

## Client Queries

Given a paymail alias `<alias>@<domain>.<tld>`, a paymail client would perform a DNS lookup for an `SRV` record matching `_bsvalias._tcp.<domain>.<tld>`. The combination of `Target` and `Port` fields are then used for Capability Discovery. Should no record be returned, a paymail client should assume a host of `<domain>.<tld>` and a port of `443`.

### Security and SRV Records

In its legacy form, that is, without DNSSEC, DNS is not a secure mechanism. It is susceptible to a range of attacks, the most serious for the purposes of the BSV Alias protocol being an intercept or man-in-the-middle (MITM) attack. In this scenario an attacker intercepts DNS queries and responds with their own data. This would allow an attacker to direct a client to a paymail implementation of their choosing, which would further allow for them to control all further communications with a client.

Note: whilst an exception to the requirement for DNSSEC is made for cases where the `SRV` record points to the same target domain as the `SRV` record itself it is best practice to enable DNSSEC in all cases.

Clients **must**, therefore, resolve hosts with the following procedure.

1. Query for an `SRV` endpoint at `_bsvalias._tcp.<domain>.<tld>.`.
1. If an `SRV` record is found and the reponse target points to `<domain>.<tld>` or `www.<domain>.<tld>`, but the response is _not_ served with a valid DNSSEC signature chain, proceed to capability discovery using the target/port combination as specified by the `SRV` record and rely on the SSL certificate for the prevention of man in the middle attacks. Note that in this scenario it is possible to still receive malicious DNS information, however SSL certificates prevent further MITM attacks.
1. If an `SRV` record is found, but the response is _not_ served with a valid DNSSEC signature chain and the response target does NOT point to `<domain>.<tld>` or `www.<domain>.<tld>`, ignore the `SRV` record completely and continue with an `A` record query as in step 5, relying on the SSL certificate for the prevention of man in the middle attacks. Note that in this scenario it is possible to still receive malicious DNS information, however SSL certificates prevent further MITM attacks. The worst that may happen here is that, for a domain that has delegated BSV Alias services out to a service provider, further lookups fail and the client cannot proceed at this time. However the option of falling back to an `A` record and relying on SSL prevents a possible denial of service vector where a malicious actor has compromised DNS.
1. If an `SRV` record is found, and the response _is_ served with a valid DNSSEC signature chain, proceed to capability discovery using the target/port combination as specified by the `SRV` record.
1. If no `SRV` record is found, instead query for an `A` record for `<domain>.<tld>.` and proceed to capability discovery, verifying that the SSL certificate presented by the remote host is valid for `<domain>.<tld>.`.

A client **must** fail the request in the following cases:
* The SSL certificate for the target domain does not verify as valid for that target domain regardless of whether the target domain is discovered via the `SRV` record or via the `A` record.
* An `SRV` record is not served; AND AND `<domain>.<tld>` does not host a valid capability discovery service.
* The `SRV` response target does NOT point `<domain>.<tld>` or `www.<domain>.<tld>`; AND the `SRV` record is not served with a valid DNSSEC signature chain; AND `<domain>.<tld>` does not host a valid capability discovery service.