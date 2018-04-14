# pubkey: public key fetcher from GitHub

## What's this?
This is a program which collects GitHub user's public keys and format them as an SSH's authorized keys format.

## How to use?
1. List multiple GitHub user IDs as following YAML format in `settings.yml`:

```
- id: <github id>
```

2. Execute the program and the formatted data will be on the standard output.

3. Copy them into .ssh/authorized_keys as needed.

## License
MIT

## Author
Tatsuya Kobayashi
