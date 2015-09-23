# Configs

[![Build Status](https://drone.io/github.com/pistarlabs/configs/status.png)](https://drone.io/github.com/pistarlabs/configs/latest) [![GoDoc](https://godoc.org/github.com/pistarlabs/configs?status.svg)](https://godoc.org/github.com/pistarlabs/configs) [![Build Status](https://travis-ci.org/pistarlabs/configs.svg?branch=master)](https://travis-ci.org/pistarlabs/configs)

Package config provides convenient access methods to configuration stored as JSON. This is a simplified library from [this](https://github.com/olebedev/config)

## Usage

Example JSON configuration:
```json
{
  "development":  {
      "database": {
          "host":"localhost",
          "username":"root",
          "password":"12345",
          "port":12345,
          "name":"dev"
        }
    },
    "production": {
      "database":{
        "host":"localhost",
        "username":"root",
        "password":"12345",
        "port":12345,
        "name":"dev"
      }
    }
}
```

We can load the JSON configuration file by using:
```go
cfg, err := Load("/path/to/config.json")
if err != nil {
  panic(err)
}
```

And we can change root path by using:
```go
// Get development environment configuration
cfg, err = cfg.Get("development")
if err != nil {
  panic(err)
}
```

Get configuration value by using:
```go
// Get database host or return empty string if not exists
host := cfg.UString("database.host")

// Get database host or return default passed value if not exists
host := cfg.UString("database.host", "default")

// Get database host and error
host, err := cfg.String("database.host")
```

These are all method to get configuration value:
```go
// Load reads a JSON configuration from given filename
func Load

// Get returns a nested config according to a dotted path.
func (*Config) Get

// Return single value or default
func (*Config) UBool
func (*Config) UFloat64
func (*Config) UInt
func (*Config) UList
func (*Config) UMap
func (*Config) UString

// Return value and error
func (*Config) Bool
func (*Config) Float64
func (*Config) Int
func (*Config) List
func (*Config) Map
func (*Config) String
```
See documentation in [GoDoc](https://godoc.org/github.com/pistarlabs/configs) for more detail
