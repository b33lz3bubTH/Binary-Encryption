### encrypt any bin first using a password
> go run encryptor.go /usr/bin/ls encrypted_binary abcd1
`abcd1` -> is the password

# use the loader to decrypt and run the binary
> ./loader encrypted_binary abcd1
