# jwt

A simple debugger for [jwt][jwt] tokens written in Go.

# Installation

```sh
go get -u github.com/romeovs/jwt
```

# Usage

To decode a JWT, just pass it as an argument:
```sh
jwt decode "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
```

You can also pipe the token into `jwt`, like so:
```sh
echo "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c" | jwt decode
```

or pass a file as the argument:
```sh
jwt decode ./file
```

## Example output

The output looks like this:
```sh
      Type  JWT
 Algorithm  HS256
   Subject  1234567890

    Issued  2018-01-18 02:30:22 +0100 CET
   Expires  <nil>
     Valid  token is valid

{
  "iat": 1516239022,
  "name": "John Doe",
  "sub": "1234567890"
}
```

[jwt]: https://jwt.io
