# Hostly - Multi-Domain File Server for Caddy

Hostly is a Caddy extension for hosting multiple domains, each with its own directory.

## Features

- Host multiple domains from a single Caddy server
- Automatic TLS certificate management with on-demand issuance
- Domain validation via filesystem check or explicit whitelist
- Optimized for Single Page Applications (SPAs)
- Security headers preconfigured

## Installation

### Pre-built Binaries

Download ready-to-run executables from the [Releases](https://github.com/andrifix/hostly/releases) page and run it:
```bash
hostly run multi-file-server
```
By default, Hostly serves files from the current directory. Create subdirectories named after your domains (e.g., ./example.com/) to host multiple sites.


### From Source

```bash
go install github.com/andrifix/hostly@latest
```

### Using Docker

You can pull the pre-built Docker image from GitHub Container Registry:

```bash
docker pull ghcr.io/andrifix/hostly
docker run -p 80:80 -p 443:443 -v /path/to/websites:/var/www/html -v data:/data ghcr.io/andrifix/hostly
```

Or build it yourself:

```bash
docker build -t hostly .
docker run -p 80:80 -p 443:443 -v /path/to/websites:/var/www/html hostly
```

## Usage

1. Create a directory structure for your domains:

```
/var/www/html/
├── example.com/
│   └── index.html
├── another-domain.com/
│   └── index.html
```

2. Start the server:

```bash
hostly run --config /path/to/Caddyfile
# OR
hostly multi-file-server --root /var/www/html/
```

## Caddyfile Variants

The project includes three Caddyfile variants:

- **Caddyfile**: Standard configuration for production use
- **Caddyfile_full**: Extended configuration with rate limiting and fail2ban support
- **Caddyfile_local**: Development configuration with local paths and stdout logging

## Domain Validation

A domain is considered valid if either:
1. It's explicitly listed in the `domains` directive
2. A corresponding directory exists in the specified path

## Custom Domain Rules

You can add custom rules for specific domains:

```caddyfile
@example host example.com
handle @example {
    root * /var/www/special/example
}
```

## Dependencies

- [Caddy v2](https://github.com/caddyserver/caddy)
- [transform-encoder](https://github.com/caddyserver/transform-encoder)
- [caddy-ratelimit](https://github.com/mholt/caddy-ratelimit)

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
