# Hyperledger Fabric Banking System

The purpose of this project is to show how the hyperledger can be feasible in a banking industry. The reason for this is to make the transaction more secure and cannot be tamper and used for fraudulent act. What I created is a working hyperledger blockchain network with chaincode that operates the transaction and a web application for the end-user to interact.

## What I've used in this project:

- `hyperledger/fabric-samples` repository. The reason this is needed is because the `/bin` directory are needed to create the hyperledger blockchain.

This came from: https://hyperledger-fabric.readthedocs.io/en/release-2.5/

- The programming languages chaincode and backend is `Golang` because they have libraries that can communicate to blockchain.

This is the library repo: https://github.com/hyperledger/fabric-sdk-go

- Lastly, the frontend I'll be using the `React` because of its Hook component like useState and useEffect.

React documentations: https://legacy.reactjs.org/docs/getting-started.html

## Prerequisites to use this:

### Mac OS

**Homebrew**: From the hyperledger-frabric website its recommended using `Homebrew` to manage the prereqs.
```bash
$ /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
$ brew --version
Homebrew 2.5.2
```

**Git**: Install the latest version of `git`
```bash
$ brew install git
$ git --version
git version 2.23.0
```

**cURL**: Install the latest version of `cURL`
```bash
$ brew install curl
$ curl --version
curl 7.64.1 (x86_64-apple-darwin19.0) libcurl/7.64.1 (SecureTransport) LibreSSL/2.8.3 zlib/1.2.11 nghttp2/1.39.2
Release-Date: 2019-03-27
```

**Docker**: Install the latest version
```bash
$ docker --version
Docker version 19.03.12, build 48a66213fe
$ docker-compose --version
docker-compose version 1.27.2, build 18f557f9
```

**Go**: Install the latest Fabric supported version of Go (required Go chaincode or SDK applications)
```bash
$ brew install go@1.23.1
$ go version
go1.23.1 darwin/amd64
```

### Windows

**WSL2**: Install the (Windows Subsystem for Linux version 2) ubuntu distro.

**Git**: Install the latest version of git if it is not already installed.
```bash
$ sudo apt-get install git
```

**cURL**
```bash
$ sudo apt-get install curl
```

**Docker**: Install the latest version of Docker if it is not already installed.

**Go**: Install the lastest version of Go.

## How to run the project:

Download the folder `hyperledger-bank` then go to that folder. Then type this command:
```bash
./setup.sh
```

Then wait for it to setup the hyperledger, backend, and frontend. Until you see this:
```bash
go-gin_1     | 2024/10/30 04:43:21 stdout: [GIN-debug] Environment variable PORT is undefined. Using port :8080 by default
go-gin_1     | 2024/10/30 04:43:21 stdout: [GIN-debug] Listening and serving HTTP on :8080
```
Or 
```bash
react-app_1  | webpack compiled successfully
react-app_1  | Compiling...
react-app_1  | Compiled successfully!
react-app_1  | webpack compiled successfully
```


Open a browser and type this URL:
```
https://localhost:3000/
```

## Walk through of the project:

Once you are in the `https://localhost:3000/` on the right top corner of the browser. You can see there's
a login button. Click that!

It will bring to the account page of **Steven** that was initialize during the `./setup.sh`.
The account has $20,000 CAD to play with it.

First we will try withdrawing some money, on the left side of the page you will see the 
`Transfer & Payment Process` section. Click the **Withdraw Money**

It will bring to `http://localhost:3000/withdraw` and there will a purple box on the center of the page.
The inputs are basically straight for choose an account and put the amount. So the example input would be:
- Chequing
- 500

Note: The Withdraw Money doesn't handle **negative integer**.

Then press **Send**, there will be an alert popping up saying
```
transaction processing!
```
Just hit **OK**. And wait for a couple of second to pop up
```
Withdraw successful!
```
Hit **OK** again.

Just hit the back arrow key of your browser. Then you will see the account balance has been updated.

Next click the `Day to Day Chequing Transaction:` on the **Bank Accounts** column. This is where you will
see all of the transaction you do with the transaction hash that came from the hyperledger blockchain.

Now let's try to deposit some money click the `Deposit Money` on the **Transfer & Payment Process** section.

The page is the same as the withdraw page, but this time it will deposit some money. Fill in the inputs so example:
- Chequing
- 1000

Note: The Deposit Money doesn't handle negative integer.

Then press **Send**, there will be an alert popping up saying
```
transaction processing!
```
Just hit **OK**. And wait for a couple of second to pop up
```
Deposit successful!
```
Hit **OK** again.

Just hit the back arrow key of your browser. Then you will see the account balance has been updated again.

You can go to the `Day to Day Chequing Transaction:` just to see the transaction hash.

Now we can now try send money to other account. Go to `Send Payment`. It will have the same design with Deposit
and Withdraw but it has another field of input.

Fill in the input:
- Account - **Chequing**
- Amount - **Up to you** as long as it does not exceed your account balance.
- To Whom: You have 2 choices.
    - NSLSC - 987654321
    - SaskEnergy - 098765432

After filling it up just hit **Send**. Then hit Ok and Ok.

Now press back arrow key of your browser. And Go to `Day to Day Chequing Transaction:` you will see
the payment was process. To confirm that the payment was received.

On your URL type:
```
http://localhost:3000/account/{the-id-you-send}
```
For NSLSC:
```
http://localhost:3000/account/987654321
```
For SaskEnergy:
```
http://localhost:3000/account/098765432
```

When you arrive to their account page you will see the amount you sent is in their account balance. And when
you go to their `Day to Day Chequing Transaction:` you will see the transaction hash is the same as yours. This goes
to show that you and the other account is part of the transaction.

Thats all for the walk through.

