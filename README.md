# Quest

## Structure
```
.
├── README.md
├── config
│   └── config.go
├── domain
│   ├── quest.go
│   └── user.go
├── go.mod
├── go.sum
├── handler
│   └── http
│       ├── quest.handler.go
│       └── user.handler.go
├── main.go
├── repository
│   ├── quest_mysql
│   │   ├── dto.go
│   │   └── quest_mysql.repository.go
│   └── user_mysql
│       ├── dto.go
│       └── user_mysql.repository.go
├── usecase
│   ├── quest.usecase.go
│   └── user.usecase.go
└── utils
    ├── edison
    │   └── edison.go
    └── mysql
        └── mysql.go
```

- **config** -- setup ENVAR config
- **domain** -- setup entity, and business process (usecase and repository interface)
- **handler** -- it's something like controller in MVC
- **repository** -- repository implementation
- **usecase** -- usecase implementation
- **utils** -- helpers

## API DOC
https://documenter.getpostman.com/view/18749474/2s8YYPHfs7
