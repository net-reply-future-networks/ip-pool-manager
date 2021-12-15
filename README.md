# ip-pool-manager

This ip-pool-manger Micro-Service manages IP addresses that are typically used to reserve resources for specific users or groups, or to ensure that certain ranges are used for specific purposes.

![overview](https://github.com/UErenReply/ip-pool-manager/blob/main/documentation/ipPool.jpg)

## Installation

```bash
git clone https://github.com/UErenReply/ip-pool-manager
go build 
```

## Usage

First run the server:

```bash
./server
```

or 
```bash
go run server.go --port 8080 --address 0.0.0.0 --redisPort 6378 --redisAddress 0.0.0.0
```

Use client:

```bash
go run client.go
```

## Testing

Run Unitests

```bash
go test
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
