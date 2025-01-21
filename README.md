# ACMEgo [![ACMEgo Badge](https://acmegobadge.shishir.dev:8445/badge.svg)](https://acmegobadge.shishir.dev:8445/badge.svg)  
**A Lightweight ACME Client Written in Go for getting **`TLS`** certificates from the Let's Encrypt CA**

ACMEgo is a fully compliant [RFC 8555](https://tools.ietf.org/html/rfc8555) (ACME) implementation written in pure Go. It is lightweight and uses Go’s standard library for simplicity and performance.

---

## Features

- Tested with **Let's Encrypt** ACME CA.  
- Implements advanced and niche aspects of RFC 8555.  
- Supports manual addition of APIs for DNS providers for automated DNS-01 challenge solving.

---

## Installation

### Requirements:

- Go 1.22+  
- Ensure the environment variable `GO111MODULE=on` is set.

### From Source

```bash
git clone https://github.com/Shishircipher/acmego.git
cd acmego
make build
```
## Usage
Navigate to the dist directory:
```
cd dist
```
**BE SURE TO:**
Create the data.json configuration file in the following format:
```
{
	"email" : "mailto:help@example.com",
	"domain" : "example.com"
}
```
Run the ACMEgo client:
```

./letsacme

```

## Warning ⚠️
For security, restrict access to the .acmego folder.

Ensure the following:

- Backup your account private key (e.g., account.key).
- Prevent unauthorized access to your account private key and domain private key.




