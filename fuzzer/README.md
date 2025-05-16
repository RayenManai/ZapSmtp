1. Download OpenSSL and add its path to `_test/unitTestConf.go`
2. Start mailhog simply with running `docker compose up -d`
3. Add mailhog config to `_test/unitTestConf.go` example:
````
Server             = "localhost"
Port        uint16 = 1025
Sender    = mail.Address{Name: "Test Sender", Address: "sender@localhost"}
Recipient = mail.Address{Name: "Test Recipient", Address: "recipient@localhost"}
````
4. `Run fuzzer/main.go`, go to `localhost:8025` for the mail client ui.