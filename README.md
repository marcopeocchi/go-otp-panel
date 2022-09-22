# OTP-Panel

All-in-one solution

![overview](https://i.ibb.co/2qLs169/null-1.png)

## Build
```sh
go build -o otp-panel *.go
```

## Run
```sh
GIN_MODE=release REDIS_ADDR=<address>:<port> REDIS_PASS=<password> ./otp-panel
```

## Example call
```sh
curl -XPOST -H "Content-type: application/json" -d '{
"message": "Your OTP is 390428",
"sender": "+1928731893",
"recipient": "+1092843098"
}' 'http://localhost:8080/api/publish'
```

```powershell
Invoke-RestMethod -Method POST -Uri 'http://localhost:8080/api/publish' -Verbose:$false -Headers @{
    'Content-type' = 'application/json'
} -Body '{"message": "Your OTP is 390428","sender": "+1928731893","recipient": "+1092843098"}'
```
