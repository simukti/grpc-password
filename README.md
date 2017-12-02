## GRPC-PASSWORD

gRPC service which only do one thing: **password hashing**.

### PASSWORD HASHING

- Password hashing using [Argon2id](https://en.wikipedia.org/wiki/Argon2)/[Bcrypt](https://en.wikipedia.org/wiki/Bcrypt) with pre-hash password HMAC-ed by secret string (default is argon2id).
- Pre-hash process will sum hash of ascii-85 plaintext password and secret key using blake2b-512.
- Final sum for bcrpyt using blake2b-384 to [fit 60 chars](https://security.stackexchange.com/questions/39849/does-bcrypt-have-a-maximum-password-length) password ascii-85 encoded.
- FInal sum for sodium will use blake2b-512 to create 80 chars ascii-85 encoded.
- In hashing create operation, only `id` will be returned for further reference (for validation/update/delete).
- Credential `id` returned by this service is globally unique id using [xid](https://github.com/rs/xid).
- Secret string and hash result will never returned from this service.
- Secret string for HMAC key is random based on `crypto/rand` byte seed.
- Every password had their own HMAC secret string (no reusable salt), so even if the original plain password is same, pre-hash result will be different.

Take a look at this sample:

#### SAMPLE DATA

First sodium hash:
```
PLAIN_PASSWORD_ORIGINAL:        12345678
PLAIN_PASSWORD_AFTER_PRE-HASH:  oUpnga;gl9$,FTtc<lediR(H>s'f]Z/2$f8Dk/`[W4Yt,KrCG<2MB'8>72<1E6=%^@)Ak1HIoIH9'I73
```

Second sodium hash:
```
PLAIN_PASSWORD_ORIGINAL:        12345678
PLAIN_PASSWORD_AFTER_PRE-HASH:  !p%tLWQG<f_Jai+a1e'J0,mM`C4kV"pNHP1^49-3LT]JBA.^Z4DD+1e_e'Z>&<0\&!&`Gk$h5OH.1(#O
```

First bcrypt hash:
```
PLAIN_PASSWORD_ORIGINAL:        12345678
PLAIN_PASSWORD_AFTER_PRE-HASH:  2o#h+A6o&XMV[f5)Kp4rnkFQ7<VLYq1:`m#qG3kd8G(,;5P+/4WG`+q]]M3G
```

Second bcrypt hash:
```
PLAIN_PASSWORD_ORIGINAL:        12345678
PLAIN_PASSWORD_AFTER_PRE-HASH:  '/^$8?DFUL0o")"R)Z@%[0>Hs[!aNaXF4g9GG9of%u1_A*s*WfBDLr>e73h6
```

### USAGE

Protocol buffer definition defined in [`password.proto`](proto/password.proto) file as follow :

```proto
service Password {
    rpc PasswordCreate (PasswordCreateRequest) returns (PasswordCreateResponse) {}
    rpc PasswordUpdate (PasswordUpdateRequest) returns (PasswordUpdateResponse) {}
    rpc PasswordDelete (PasswordDeleteRequest) returns (PasswordDeleteResponse) {}
    rpc PasswordValidate (PasswordValidateRequest) returns (PasswordValidateResponse) {}
}
```

#### CLIENT EXAMPLE

Run gRPC server with `go run main.go`.

Open another console tab and run password sample client `go run ./_examples/password_client/main.go`.
It should output log like these :

```
$ go run ./_examples/password_client/main.go
2017/12/02 15:24:41 ######### CREATE NEW PASSWORD CREDENTIAL #########
2017/12/02 15:24:42 &PasswordCreateResponse{Success:true,Id:b8h66i86f5s22tk9cmn0,}
2017/12/02 15:24:42 ######### VALIDATE CURRENT PASSWORD #########
2017/12/02 15:24:42 &PasswordValidateResponse{Success:true,Id:b8h66i86f5s22tk9cmn0,}
2017/12/02 15:24:42 ######### VALIDATE NEW PASSWORD #########
2017/12/02 15:24:42 rpc error: code = Unknown desc = ERROR_WRONG_PASSWORD
2017/12/02 15:24:42 nil
2017/12/02 15:24:42 ######### UPDATE TO NEW PASSWORD #########
2017/12/02 15:24:42 &PasswordUpdateResponse{Success:true,Id:b8h66i86f5s22tk9cmn0,}
2017/12/02 15:24:42 ######### VALIDATE OLD PASSWORD #########
2017/12/02 15:24:42 rpc error: code = Unknown desc = ERROR_WRONG_PASSWORD
2017/12/02 15:24:42 nil
2017/12/02 15:24:42 ######### VALIDATE NEW PASSWORD #########
2017/12/02 15:24:43 &PasswordValidateResponse{Success:true,Id:b8h66i86f5s22tk9cmn0,}
2017/12/02 15:24:43 ######### DELETE CREDENTIAL #########
2017/12/02 15:24:43 &PasswordDeleteResponse{Success:true,}
2017/12/02 15:24:43 ######### VALIDATE NEW PASSWORD #########
2017/12/02 15:24:43 rpc error: code = Unknown desc = NOT_FOUND
2017/12/02 15:24:43 nil
```

#### STORAGE

Default storage is `in-memory` storage.

For mysql storage:

0. Create database `YOURDBNAME`
1. Set `APP_DB_TYPE` environment variable to `mysql`
2. Set `APP_DB_DSN` to `"username:password@tcp(127.0.0.1:3306)/YOURDBNAME?charset=utf8mb4&collation=utf8mb4_unicode_ci"`.

See [.env file](.env.dist).

### FEATURES

- [x] Password hashing service
- [x] Argon2 hashing
- [x] Bcrypt hashing
- [x] InMemory storage
- [x] MySQL storage

### TODO

- [ ] Public docker image
- [ ] Unit Tests

### IS PASSWORD HASHING SECURE ??

Nothing is 100% absolute secure. At least we have to make it harder to break.
Please create issue if you found bug, or have any idea for improvement.

### REFERENCES

- https://password-hashing.net/
- https://crackstation.net/hashing-security.htm
- https://github.com/riverrun/comeonin/wiki/Choosing-the-password-hashing-algorithm
- https://pthree.org/2016/06/28/lets-talk-password-hashing/
- https://password.kaspersky.com/


### LICENSE

MIT